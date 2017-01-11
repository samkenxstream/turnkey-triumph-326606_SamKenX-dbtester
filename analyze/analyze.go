// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analyze

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/gyuho/dataframe"
	"github.com/spf13/cobra"
)

// Command implements 'analyze' command.
var Command = &cobra.Command{
	Use:   "analyze",
	Short: "Analyzes test dbtester test results.",
	RunE:  commandFunc,
}

var configPath string

func init() {
	Command.PersistentFlags().StringVarP(&configPath, "config", "c", "", "YAML configuration file path.")
}

var columnsToAggregate = []string{
	"UNIX-TS", "CPU-NUM", "VMRSS-NUM",
	"READS-COMPLETED",
	"READS-COMPLETED-DIFF",
	"SECTORS-READ",
	"SECTORS-READ-DIFF",
	"WRITES-COMPLETED",
	"WRITES-COMPLETED-DIFF",
	"SECTORS-WRITTEN",
	"SECTORS-WRITTEN-DIFF",
	"RECEIVE-BYTES-NUM",
	"RECEIVE-BYTES-NUM-DIFF",
	"TRANSMIT-BYTES-NUM",
	"TRANSMIT-BYTES-NUM-DIFF",
	"EXTRA",
}

func commandFunc(cmd *cobra.Command, args []string) error {
	cfg, err := ReadConfig(configPath)
	if err != nil {
		return err
	}

	println()
	plog.Println("Step 1: aggregating each database...")
	for step1Idx, elem := range cfg.Step1 {
		var (
			frames               = []dataframe.Frame{}
			maxCommonMinUnixTime int64
			maxCommonMaxUnixTime int64
		)
		for i, monitorPath := range elem.DataPathList {
			plog.Printf("Step 1-%d-%d: creating dataframe from %s", step1Idx, i, monitorPath)
			originalFrame, err := dataframe.NewFromCSV(nil, monitorPath)
			if err != nil {
				return err
			}

			newFrame := dataframe.New()
			var tsc dataframe.Column
			for _, name := range columnsToAggregate {
				cmn, err := originalFrame.GetColumn(name)
				if err != nil {
					return err
				}
				if name == "UNIX-TS" {
					tsc = cmn
				}
				if err = newFrame.AddColumn(cmn); err != nil {
					return err
				}
			}
			frames = append(frames, newFrame)

			fv, ok := tsc.FrontNonNil()
			if !ok {
				return fmt.Errorf("FrontNonNil %s has empty Unix time %v", monitorPath, fv)
			}
			fs, ok := fv.ToString()
			if !ok {
				return fmt.Errorf("cannot ToString %v", fv)
			}
			fd, err := strconv.ParseInt(fs, 10, 64)
			if err != nil {
				return err
			}
			bv, ok := tsc.BackNonNil()
			if !ok {
				return fmt.Errorf("BackNonNil %s has empty Unix time %v", monitorPath, fv)
			}
			bs, ok := bv.ToString()
			if !ok {
				return fmt.Errorf("cannot ToString %v", bv)
			}
			bd, err := strconv.ParseInt(bs, 10, 64)
			if err != nil {
				return err
			}
			if i == 0 {
				maxCommonMinUnixTime = fd
				maxCommonMaxUnixTime = bd
			}
			if maxCommonMinUnixTime < fd {
				maxCommonMinUnixTime = fd
			}
			if maxCommonMaxUnixTime > bd {
				maxCommonMaxUnixTime = bd
			}
		}

		// monitor CSVs from multiple servers, and want them to have equal number of rows
		// Truncate all rows before maxCommonMinUnixTime and after maxCommonMinUnixTime
		minTS := fmt.Sprintf("%d", maxCommonMinUnixTime)
		maxTS := fmt.Sprintf("%d", maxCommonMaxUnixTime)
		aggregatedFrame := dataframe.New()
		for i := range frames {
			uc, err := frames[i].GetColumn("UNIX-TS")
			if err != nil {
				return err
			}
			j, ok := uc.FindValue(dataframe.NewStringValue(minTS))
			if !ok {
				return fmt.Errorf("%v does not exist in %s", minTS, elem.DataPathList[i])
			}
			k, ok := uc.FindValue(dataframe.NewStringValue(maxTS))
			if !ok {
				return fmt.Errorf("%v does not exist in %s", maxTS, elem.DataPathList[i])
			}

			for _, header := range frames[i].GetHeader() {
				if i > 0 && header == "UNIX-TS" {
					continue
				}
				var col dataframe.Column
				col, err = frames[i].GetColumn(header)
				if err != nil {
					return err
				}
				if err = col.KeepRows(j, k+1); err != nil {
					return err
				}

				// update column name with database name and its index
				// all in one aggregated CSV file
				if header != "UNIX-TS" {
					switch header {
					case "CPU-NUM":
						header = "CPU"
					case "VMRSS-NUM":
						header = "VMRSS-MB"

						// to bytes to mb
						colN := col.RowNumber()
						for rowIdx := 0; rowIdx < colN; rowIdx++ {
							var rowV dataframe.Value
							rowV, err = col.GetValue(rowIdx)
							if err != nil {
								return err
							}
							fv, _ := rowV.ToNumber()
							frv := float64(fv) * 0.000001
							if err = col.SetValue(rowIdx, dataframe.NewStringValue(fmt.Sprintf("%.2f", frv))); err != nil {
								return err
							}
						}
					case "EXTRA":
						header = "CLIENT-NUM"
					}

					col.UpdateHeader(fmt.Sprintf("%s-%d", header, i+1))
				}

				if err = aggregatedFrame.AddColumn(col); err != nil {
					return err
				}
			}
		}

		plog.Printf("Step 1-%d-%d: creating dataframe from %s", step1Idx, len(elem.DataPathList), elem.DataBenchmarkPath)
		colMonitorUnixTs, err := aggregatedFrame.GetColumn("UNIX-TS")
		if err != nil {
			return err
		}

		// need to combine aggregatedFrame to benchResultFrame by unix timestamps
		benchResultFrame, err := dataframe.NewFromCSV(nil, elem.DataBenchmarkPath)
		if err != nil {
			return err
		}
		colBenchUnixTs, err := benchResultFrame.GetColumn("UNIX-TS")
		if err != nil {
			return err
		}

		fv, ok := colBenchUnixTs.FrontNonNil()
		if !ok {
			return fmt.Errorf("FrontNonNil %s has empty Unix time %v", elem.DataBenchmarkPath, fv)
		}
		startRowMonitor, ok := colMonitorUnixTs.FindValue(fv)
		if !ok {
			return fmt.Errorf("%v is not found in monitor results %q", fv, elem.DataPathList)
		}
		bv, ok := colBenchUnixTs.BackNonNil()
		if !ok {
			return fmt.Errorf("BackNonNil %s has empty Unix time %v", elem.DataBenchmarkPath, bv)
		}
		endRowMonitor, ok := colMonitorUnixTs.FindValue(bv)
		if !ok { // monitor short of rows
			endRowMonitor = colMonitorUnixTs.RowNumber() - 1
		}

		var benchLastIdx int
		for _, col := range benchResultFrame.GetColumns() {
			if benchLastIdx == 0 {
				benchLastIdx = col.RowNumber()
			}
			if benchLastIdx > col.RowNumber() {
				benchLastIdx = col.RowNumber()
			}
		}
		benchLastIdx--

		if benchLastIdx+1 < endRowMonitor-startRowMonitor+1 { // benchmark is short of rows
			endRowMonitor = startRowMonitor + benchLastIdx
		} else { // monitor is short of rows
			benchLastIdx = endRowMonitor - startRowMonitor
		}

		for _, hd := range aggregatedFrame.GetHeader() {
			if hd == "UNIX-TS" {
				continue
			}
			var col dataframe.Column
			col, err = aggregatedFrame.GetColumn(hd)
			if err != nil {
				return err
			}
			if err = col.KeepRows(startRowMonitor, endRowMonitor); err != nil {
				return err
			}
			if err = benchResultFrame.AddColumn(col); err != nil {
				return err
			}
		}

		plog.Printf("Step 1-%d-%d: calculating average, cumulative values", step1Idx, len(elem.DataPathList)+1)
		var (
			sampleSize              = float64(len(elem.DataPathList))
			cumulativeThroughputCol = dataframe.NewColumn("CUMULATIVE-AVG-THROUGHPUT")
			totalThrougput          int
			avgCPUCol               = dataframe.NewColumn("AVG-CPU")
			avgVMRSSMBCol           = dataframe.NewColumn("AVG-VMRSS-MB")

			// TODO: average value of disk stats, network stats
		)
		for i := 0; i < benchLastIdx; i++ {
			var (
				cpuTotal    float64
				memoryTotal float64
			)
			for _, col := range benchResultFrame.GetColumns() {
				var rv dataframe.Value
				rv, err = col.GetValue(i)
				if err != nil {
					return err
				}

				fv, _ := rv.ToNumber()
				switch {
				case strings.HasPrefix(col.GetHeader(), "CPU-"):
					cpuTotal += fv

				case strings.HasPrefix(col.GetHeader(), "VMRSS-"):
					memoryTotal += fv

				case col.GetHeader() == "AVG-THROUGHPUT":
					fv, _ := rv.ToNumber()
					totalThrougput += int(fv)
					cumulativeThroughputCol.PushBack(dataframe.NewStringValue(totalThrougput))
				}
			}
			avgCPUCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", cpuTotal/sampleSize)))
			avgVMRSSMBCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", memoryTotal/sampleSize)))
		}

		plog.Printf("Step 1-%d-%d: combine %s and %q", step1Idx, len(elem.DataPathList)+2, elem.DataBenchmarkPath, elem.DataPathList)
		unixTsCol, err := benchResultFrame.GetColumn("UNIX-TS")
		if err != nil {
			return err
		}
		latencyCol, err := benchResultFrame.GetColumn("AVG-LATENCY-MS")
		if err != nil {
			return err
		}
		throughputCol, err := benchResultFrame.GetColumn("AVG-THROUGHPUT")
		if err != nil {
			return err
		}

		aggFr := dataframe.New()
		aggFr.AddColumn(unixTsCol)
		aggFr.AddColumn(latencyCol)
		aggFr.AddColumn(throughputCol)
		aggFr.AddColumn(cumulativeThroughputCol)
		for _, hd := range benchResultFrame.GetHeader() {
			col, err := benchResultFrame.GetColumn(hd)
			if err != nil {
				return err
			}
			switch {
			case strings.HasPrefix(hd, "CPU-"):
				aggFr.AddColumn(col)
			case strings.HasPrefix(hd, "VMRSS-"):
				aggFr.AddColumn(col)
			}
		}
		aggFr.AddColumn(avgCPUCol)
		aggFr.AddColumn(avgVMRSSMBCol)

		plog.Printf("Step 1-%d-%d: saving to %s", step1Idx, len(elem.DataPathList)+3, elem.OutputPath)
		if err := aggFr.ToCSV(elem.OutputPath); err != nil {
			return err
		}
		println()
	}

	println()
	plog.Println("Step 2: aggregating aggregates...")
	for step2Idx, elem := range cfg.Step2 {
		var (
			frames  = []dataframe.Frame{}
			maxSize int
		)
		for i, data := range elem.DataList {
			plog.Printf("Step 2-%d-%d: creating dataframe from %s...", step2Idx, i, data.Path)
			fr, err := dataframe.NewFromCSV(nil, data.Path)
			if err != nil {
				return err
			}
			frames = append(frames, fr)

			col, err := fr.GetColumn("UNIX-TS")
			if err != nil {
				return err
			}
			rNum := col.RowNumber()
			if maxSize < rNum {
				maxSize = rNum
			}
		}

		nf := dataframe.New()
		secondCol := dataframe.NewColumn("SECOND")
		for i := 0; i < maxSize; i++ {
			secondCol.PushBack(dataframe.NewStringValue(i))
		}
		nf.AddColumn(secondCol)

		// TODO: keep disk, network stats columns
		colsToKeep := []string{"AVG-LATENCY-MS", "AVG-THROUGHPUT", "CUMULATIVE-AVG-THROUGHPUT", "AVG-CPU", "AVG-VMRSS-MB"}
		for i, fr := range frames {
			dbID := elem.DataList[i].Name
			plog.Printf("Step 2-%d-%d: cleaning up %s...", step2Idx, i, dbID)
			for _, col := range fr.GetColumns() {
				toSkip := true
				for _, cv := range colsToKeep {
					if col.GetHeader() == cv {
						toSkip = false
						break
					}
				}
				if toSkip {
					continue
				}
				if err := col.Appends(dataframe.NewStringValueNil(), maxSize); err != nil {
					return err
				}
				col.UpdateHeader(fmt.Sprintf("%s-%s", col.GetHeader(), dbID))
				nf.AddColumn(col)
			}
		}

		plog.Printf("Step 2-%d: saving to %s...", step2Idx, elem.OutputPath)
		if err := nf.ToCSV(elem.OutputPath); err != nil {
			return err
		}
	}

	println()
	plog.Println("Step 3: plotting...")

	plot.DefaultFont = "Helvetica"
	plotter.DefaultLineStyle.Width = vg.Points(1.5)
	plotter.DefaultGlyphStyle.Radius = vg.Points(2.0)
	var (
		plotWidth  = 12 * vg.Inch
		plotHeight = 8 * vg.Inch
	)
	for step3Idx, elem := range cfg.Step3 {
		fr, err := dataframe.NewFromCSV(nil, elem.DataPath)
		if err != nil {
			return err
		}
		plog.Printf("Step 3-%d: %s with %q", step3Idx, elem.DataPath, fr.GetHeader())

		for i, pelem := range elem.PlotList {
			plog.Printf("Step 3-%d-%d: %s at %q", step3Idx, i, pelem.YAxis, pelem.OutputPathList)
			pl, err := plot.New()
			if err != nil {
				return err
			}
			pl.Title.Text = fmt.Sprintf("%s, %s", cfg.Titles[step3Idx], pelem.YAxis)
			pl.X.Label.Text = pelem.XAxis
			pl.Y.Label.Text = pelem.YAxis
			pl.Legend.Top = true

			// var args []interface{}
			// for _, line := range pelem.Lines {
			// 	col, err := fr.GetColumn(line.Column)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	pt, err := points(col)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	args = append(args, line.Legend, pt)
			// }
			// if err = plotutil.AddLines(pl, args...); err != nil {
			// 	return err
			// }

			var ps []plot.Plotter
			for j, line := range pelem.Lines {
				col, err := fr.GetColumn(line.Column)
				if err != nil {
					return err
				}
				pt, err := points(col)
				if err != nil {
					return err
				}

				l, err := plotter.NewLine(pt)
				if err != nil {
					return err
				}
				l.Color = getRGB(line.Legend, j)
				l.Dashes = plotutil.Dashes(j)
				ps = append(ps, l)

				pl.Legend.Add(line.Legend, l)
			}
			pl.Add(ps...)

			for _, outputPath := range pelem.OutputPathList {
				if err = pl.Save(plotWidth, plotHeight, outputPath); err != nil {
					return err
				}
			}
		}
	}

	println()
	plog.Println("Step 4: writing README...")
	rdBuf := new(bytes.Buffer)
	rdBuf.WriteString("\n\n")
	rdBuf.WriteString(cfg.Step4.Preface)
	rdBuf.WriteString("\n\n\n")
	for i, result := range cfg.Step4.Results {
		rdBuf.WriteString(fmt.Sprintf("<br><br>\n##### %s", cfg.Titles[i]))
		rdBuf.WriteString("\n\n")
		for _, img := range result.Images {
			imgPath := ""
			switch img.ImageType {
			case "local":
				imgPath = "./" + filepath.Base(img.ImagePath)
				rdBuf.WriteString(fmt.Sprintf("![%s](%s)\n\n", img.ImageTitle, imgPath))
			case "remote":
				rdBuf.WriteString(fmt.Sprintf(`<img src="%s" alt="%s">`, img.ImagePath, img.ImageTitle))
				rdBuf.WriteString("\n\n")
			default:
				return fmt.Errorf("%s is not supported", img.ImageType)
			}
		}
		rdBuf.WriteString("\n\n")
	}
	if err := toFile(rdBuf.String(), cfg.Step4.OutputPath); err != nil {
		return err
	}

	println()
	plog.Println("FINISHED!")
	return nil
}

func toFile(txt, fpath string) error {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
	}
	defer f.Close()
	if _, err := f.WriteString(txt); err != nil {
		return err
	}
	return nil
}

func points(col dataframe.Column) (plotter.XYs, error) {
	bv, ok := col.BackNonNil()
	if !ok {
		return nil, fmt.Errorf("BackNonNil not found")
	}
	rowN, ok := col.FindLastValue(bv)
	if !ok {
		return nil, fmt.Errorf("not found %v", bv)
	}
	pts := make(plotter.XYs, rowN)
	for i := range pts {
		v, err := col.GetValue(i)
		if err != nil {
			return nil, err
		}
		n, _ := v.ToNumber()
		pts[i].X = float64(i)
		pts[i].Y = n
	}
	return pts, nil
}
