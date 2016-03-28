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
	"fmt"
	"strconv"
	"strings"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/gyuho/dataframe"
	"github.com/spf13/cobra"
)

type (
	Flags struct {
		OutputPath          string
		BenchmarkFilePath   string
		RawMonitorFilePaths []string
		AggregatedFilePaths []string
		AggAggFilePath      string
		ImageFormat         string
		ImageTitle          string
		MultiTagTitle       string
		SameDatabase        bool
	}
)

var (
	Command = &cobra.Command{
		Use:   "analyze",
		Short: "Analyzes test results specific to dbtester.",
		RunE:  CommandFunc,
	}

	globalFlags = Flags{}
)

func init() {
	Command.PersistentFlags().StringVarP(&globalFlags.OutputPath, "output", "o", "", "Output file path.")

	Command.PersistentFlags().StringVarP(&globalFlags.BenchmarkFilePath, "bench-file-path", "b", "", "Benchmark result CSV file path.")
	Command.PersistentFlags().StringSliceVarP(&globalFlags.RawMonitorFilePaths, "monitor-data-file-paths", "m", []string{}, "Monitor CSV file paths.")

	Command.PersistentFlags().StringSliceVarP(&globalFlags.AggregatedFilePaths, "aggregated-file-paths", "a", []string{}, "Already aggregated file paths per database.")

	Command.PersistentFlags().StringVarP(&globalFlags.AggAggFilePath, "file-to-plot", "p", "", "Aggregated CSV file path to plot.")
	Command.PersistentFlags().StringVarP(&globalFlags.ImageFormat, "image-format", "f", "png", "Image format (png, svg).")
	Command.PersistentFlags().StringVarP(&globalFlags.ImageTitle, "image-title", "t", "", "Image title.")
	Command.PersistentFlags().StringVarP(&globalFlags.MultiTagTitle, "multi-tag-title", "g", "", "Special title for *multi test.")
	Command.PersistentFlags().BoolVarP(&globalFlags.SameDatabase, "same-database", "s", false, "'true' when testing same database.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	switch {
	case len(globalFlags.RawMonitorFilePaths) > 0:
		fr, err := aggBenchAndMonitor(globalFlags.BenchmarkFilePath, globalFlags.RawMonitorFilePaths...)
		if err != nil {
			return err
		}
		if err := fr.ToCSV(globalFlags.OutputPath); err != nil {
			return err
		}

	case len(globalFlags.AggregatedFilePaths) > 0:
		fr, err := aggAgg(globalFlags.AggregatedFilePaths...)
		if err != nil {
			return err
		}
		if err := fr.ToCSV(globalFlags.OutputPath); err != nil {
			return err
		}

	case len(globalFlags.AggAggFilePath) > 0:
		if err := plotAggAgg(globalFlags.AggAggFilePath, globalFlags.OutputPath, globalFlags.ImageFormat, globalFlags.ImageTitle, globalFlags.MultiTagTitle); err != nil {
			return err
		}
	}

	return nil
}

// aggMonitor aggregates monitor CSV files.
func aggMonitor(monitorPaths ...string) (dataframe.Frame, error) {
	if len(monitorPaths) == 0 {
		return nil, fmt.Errorf("no file specified")
	}

	var (
		frames               = []dataframe.Frame{}
		maxCommonMinUnixTime int64
		maxCommonMaxUnixTime int64
	)
	for i, fpath := range monitorPaths {
		fr, err := dataframe.NewFromCSV(nil, fpath)
		if err != nil {
			return nil, err
		}
		nf := dataframe.New()
		c1, err := fr.GetColumn("unix_ts")
		if err != nil {
			return nil, err
		}
		if err = nf.AddColumn(c1); err != nil {
			return nil, err
		}
		c2, err := fr.GetColumn("CpuUsageFloat64")
		if err != nil {
			return nil, err
		}
		if err = nf.AddColumn(c2); err != nil {
			return nil, err
		}
		c3, err := fr.GetColumn("VmRSSBytes")
		if err != nil {
			return nil, err
		}
		if err = nf.AddColumn(c3); err != nil {
			return nil, err
		}
		frames = append(frames, nf)

		fv, ok := c1.FrontNonNil()
		if !ok {
			return nil, fmt.Errorf("FrontNonNil %s has empty Unix time %v", fpath, fv)
		}
		fs, ok := fv.ToString()
		if !ok {
			return nil, fmt.Errorf("cannot ToString %v", fv)
		}
		fd, err := strconv.ParseInt(fs, 10, 64)
		if err != nil {
			return nil, err
		}
		bv, ok := c1.BackNonNil()
		if !ok {
			return nil, fmt.Errorf("BackNonNil %s has empty Unix time %v", fpath, fv)
		}
		bs, ok := bv.ToString()
		if !ok {
			return nil, fmt.Errorf("cannot ToString %v", bv)
		}
		bd, err := strconv.ParseInt(bs, 10, 64)
		if err != nil {
			return nil, err
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

	// make all columns have equal row number, based on the column unix_ts
	// truncate all rows before maxCommonMinUnixTime and after maxCommonMinUnixTime
	minTS := fmt.Sprintf("%d", maxCommonMinUnixTime)
	maxTS := fmt.Sprintf("%d", maxCommonMaxUnixTime)
	nf := dataframe.New()
	for i := range frames {
		uc, err := frames[i].GetColumn("unix_ts")
		if err != nil {
			return nil, err
		}
		j, ok := uc.FindValue(dataframe.NewStringValue(minTS))
		if !ok {
			return nil, fmt.Errorf("%v does not exist in %s", minTS, monitorPaths[i])
		}
		k, ok := uc.FindValue(dataframe.NewStringValue(maxTS))
		if !ok {
			return nil, fmt.Errorf("%v does not exist in %s", maxTS, monitorPaths[i])
		}

		for _, hd := range frames[i].GetHeader() {
			if i > 0 && hd == "unix_ts" {
				continue
			}
			var col dataframe.Column
			col, err = frames[i].GetColumn(hd)
			if err != nil {
				return nil, err
			}
			if err = col.KeepRows(j, k+1); err != nil {
				return nil, err
			}
			if hd != "unix_ts" {
				switch hd {
				case "CpuUsageFloat64":
					hd = "cpu"
				case "VmRSSBytes":
					hd = "memory_mb"

					// to bytes to mb
					colN := col.RowNumber()
					for rowIdx := 0; rowIdx < colN; rowIdx++ {
						var rowV dataframe.Value
						rowV, err = col.GetValue(rowIdx)
						if err != nil {
							return nil, err
						}
						fv, _ := rowV.ToNumber()
						nfv := float64(fv) * 0.000001
						if err = col.SetValue(rowIdx, dataframe.NewStringValue(fmt.Sprintf("%.2f", nfv))); err != nil {
							return nil, err
						}
					}
				}
				col.UpdateHeader(fmt.Sprintf("%s_%d", hd, i+1))
			}
			if err = nf.AddColumn(col); err != nil {
				return nil, err
			}
		}
	}
	return nf, nil
}

// aggBenchAndMonitor combines benchmark latency results and monitor CSV files.
func aggBenchAndMonitor(benchPath string, monitorPaths ...string) (dataframe.Frame, error) {
	frMonitor, err := aggMonitor(monitorPaths...)
	if err != nil {
		return nil, err
	}
	colMonitorUnixTs, err := frMonitor.GetColumn("unix_ts")
	if err != nil {
		return nil, err
	}

	// need to combine frMonitor to frBench
	frBench, err := dataframe.NewFromCSV(nil, benchPath)
	if err != nil {
		return nil, err
	}
	colBenchUnixTs, err := frBench.GetColumn("unix_ts")
	if err != nil {
		return nil, err
	}

	fv, ok := colBenchUnixTs.FrontNonNil()
	if !ok {
		return nil, fmt.Errorf("FrontNonNil %s has empty Unix time %v", benchPath, fv)
	}
	startRowMonitor, ok := colMonitorUnixTs.FindValue(fv)
	if !ok {
		return nil, fmt.Errorf("%v is not found in monitor results %q", fv, monitorPaths)
	}
	bv, ok := colBenchUnixTs.BackNonNil()
	if !ok {
		return nil, fmt.Errorf("BackNonNil %s has empty Unix time %v", benchPath, bv)
	}
	endRowMonitor, ok := colMonitorUnixTs.FindValue(bv)
	if !ok { // monitor short of rows
		endRowMonitor = colMonitorUnixTs.RowNumber() - 1
	}

	var benchLastIdx int
	for _, col := range frBench.GetColumns() {
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

	for _, hd := range frMonitor.GetHeader() {
		if hd == "unix_ts" {
			continue
		}
		var col dataframe.Column
		col, err = frMonitor.GetColumn(hd)
		if err != nil {
			return nil, err
		}
		if err = col.KeepRows(startRowMonitor, endRowMonitor); err != nil {
			return nil, err
		}
		if err = frBench.AddColumn(col); err != nil {
			return nil, err
		}
	}

	var (
		sampleSize              = float64(len(monitorPaths))
		cumulativeThroughputCol = dataframe.NewColumn("cumulative_throughput")
		totalThrougput          int
		avgCpuCol               = dataframe.NewColumn("avg_cpu")
		avgMemCol               = dataframe.NewColumn("avg_memory_mb")
	)
	for i := 0; i < benchLastIdx; i++ {
		var (
			cpuTotal    float64
			memoryTotal float64
		)
		for _, col := range frBench.GetColumns() {
			var rv dataframe.Value
			rv, err = col.GetValue(i)
			if err != nil {
				return nil, err
			}
			fv, _ := rv.ToNumber()
			switch {
			case strings.HasPrefix(col.GetHeader(), "cpu_"):
				cpuTotal += fv
			case strings.HasPrefix(col.GetHeader(), "memory_"):
				memoryTotal += fv
			case col.GetHeader() == "throughput":
				fv, _ := rv.ToNumber()
				totalThrougput += int(fv)
				cumulativeThroughputCol.PushBack(dataframe.NewStringValue(totalThrougput))
			}
		}
		avgCpuCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", cpuTotal/sampleSize)))
		avgMemCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", memoryTotal/sampleSize)))
	}

	unixTsCol, err := frBench.GetColumn("unix_ts")
	if err != nil {
		return nil, err
	}
	latencyCol, err := frBench.GetColumn("avg_latency_ms")
	if err != nil {
		return nil, err
	}
	throughputCol, err := frBench.GetColumn("throughput")
	if err != nil {
		return nil, err
	}

	nf := dataframe.New()
	nf.AddColumn(unixTsCol)
	nf.AddColumn(latencyCol)
	nf.AddColumn(throughputCol)
	nf.AddColumn(cumulativeThroughputCol)

	for _, hd := range frBench.GetHeader() {
		col, err := frBench.GetColumn(hd)
		if err != nil {
			return nil, err
		}
		switch {
		case strings.HasPrefix(hd, "cpu_"):
			nf.AddColumn(col)
		case strings.HasPrefix(hd, "memory_"):
			nf.AddColumn(col)
		}
	}

	nf.AddColumn(avgCpuCol)
	nf.AddColumn(avgMemCol)

	return nf, nil
}

func aggAgg(fpaths ...string) (dataframe.Frame, error) {
	if len(fpaths) == 0 {
		return nil, fmt.Errorf("no file specified")
	}
	var (
		frames  = []dataframe.Frame{}
		maxSize int
	)
	for _, fpath := range fpaths {
		fr, err := dataframe.NewFromCSV(nil, fpath)
		if err != nil {
			return nil, err
		}
		frames = append(frames, fr)

		col, err := fr.GetColumn("unix_ts")
		if err != nil {
			return nil, err
		}
		rNum := col.RowNumber()
		if maxSize < rNum {
			maxSize = rNum
		}
	}

	nf := dataframe.New()
	secondCol := dataframe.NewColumn("second")
	for i := 0; i < maxSize; i++ {
		secondCol.PushBack(dataframe.NewStringValue(i))
	}
	nf.AddColumn(secondCol)

	colsToKeep := []string{"avg_latency_ms", "throughput", "cumulative_throughput", "avg_cpu", "avg_memory_mb"}
	for i, fr := range frames {
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
				return nil, err
			}

			var dbID string
			switch {
			case strings.Contains(fpaths[i], "consul"):
				dbID = "consul"
			case strings.Contains(fpaths[i], "etcd2"):
				dbID = "etcd2"
			case strings.Contains(fpaths[i], "etcdmulti"):
				dbID = "etcd3multi"
			case strings.Contains(fpaths[i], "etcd"):
				dbID = "etcd3"
			case strings.Contains(fpaths[i], "zk"):
				dbID = "zk"
			}

			if !globalFlags.SameDatabase {
				col.UpdateHeader(fmt.Sprintf("%s_%s", col.GetHeader(), dbID))
			} else {
				col.UpdateHeader(fmt.Sprintf("%s_%s_%d", col.GetHeader(), dbID, i))
			}
			nf.AddColumn(col)
		}
	}

	return nf, nil
}

func plotAggAgg(fpath, outputPath, imageFormat, imageTitle, multiTagTitle string) error {
	fr, err := dataframe.NewFromCSV(nil, fpath)
	if err != nil {
		return err
	}

	plot.DefaultFont = "Helvetica"
	plotter.DefaultLineStyle.Width = vg.Points(1.5)
	plotter.DefaultGlyphStyle.Radius = vg.Points(2.0)
	var (
		defaultPlotWidth  = 12 * vg.Inch
		defaultPlotHeight = 8 * vg.Inch
		avgLatencyPath    = outputPath + fmt.Sprintf("-avg-latency-ms.%s", imageFormat)
		throughputPath    = outputPath + fmt.Sprintf("-throughput.%s", imageFormat)
		avgCpuPath        = outputPath + fmt.Sprintf("-avg-cpu.%s", imageFormat)
		avgMemPath        = outputPath + fmt.Sprintf("-avg-mem.%s", imageFormat)
	)

	if globalFlags.SameDatabase {
		plotAvgLatencyEtcd3_0, err := fr.GetColumn("avg_latency_ms_etcd3_0")
		if err != nil {
			return err
		}
		plotAvgLatencyEtcd3Points_0, err := points(plotAvgLatencyEtcd3_0)
		if err != nil {
			return err
		}
		plotAvgLatencyEtcd3_1, err := fr.GetColumn("avg_latency_ms_etcd3_1")
		if err != nil {
			return err
		}
		plotAvgLatencyEtcd3Points_1, err := points(plotAvgLatencyEtcd3_1)
		if err != nil {
			return err
		}
		plotAvgLatency, err := plot.New()
		if err != nil {
			return err
		}
		plotAvgLatency.Title.Text = fmt.Sprintf("%s, Latency", imageTitle)
		plotAvgLatency.X.Label.Text = "second"
		plotAvgLatency.Y.Label.Text = "Latency(ms)"
		plotAvgLatency.Legend.Top = true
		if err = plotutil.AddLines(
			plotAvgLatency,
			"etcd3_0", plotAvgLatencyEtcd3Points_0,
			"etcd3_1", plotAvgLatencyEtcd3Points_1,
		); err != nil {
			return err
		}
		if err = plotAvgLatency.Save(defaultPlotWidth, defaultPlotHeight, avgLatencyPath); err != nil {
			return err
		}

		plotThroughputEtcd3_0, err := fr.GetColumn("throughput_etcd3_0")
		if err != nil {
			return err
		}
		plotThroughputEtcd3Points_0, err := points(plotThroughputEtcd3_0)
		if err != nil {
			return err
		}
		plotThroughputEtcd3_1, err := fr.GetColumn("throughput_etcd3_1")
		if err != nil {
			return err
		}
		plotThroughputEtcd3Points_1, err := points(plotThroughputEtcd3_1)
		if err != nil {
			return err
		}
		plotThroughput, err := plot.New()
		if err != nil {
			return err
		}
		plotThroughput.Title.Text = fmt.Sprintf("%s, Throughput", imageTitle)
		plotThroughput.X.Label.Text = "second"
		plotThroughput.Y.Label.Text = "Throughput"
		plotThroughput.Legend.Top = true
		if err = plotutil.AddLines(
			plotThroughput,
			"etcd3_0", plotThroughputEtcd3Points_0,
			"etcd3_1", plotThroughputEtcd3Points_1,
		); err != nil {
			return err
		}
		if err = plotThroughput.Save(defaultPlotWidth, defaultPlotHeight, throughputPath); err != nil {
			return err
		}

		plotAvgCpuEtcd3_0, err := fr.GetColumn("avg_cpu_etcd3_0")
		if err != nil {
			return err
		}
		plotAvgCpuEtcd3Points_0, err := points(plotAvgCpuEtcd3_0)
		if err != nil {
			return err
		}
		plotAvgCpuEtcd3_1, err := fr.GetColumn("avg_cpu_etcd3_1")
		if err != nil {
			return err
		}
		plotAvgCpuEtcd3Points_1, err := points(plotAvgCpuEtcd3_1)
		if err != nil {
			return err
		}
		plotAvgCpu, err := plot.New()
		if err != nil {
			return err
		}
		plotAvgCpu.Title.Text = fmt.Sprintf("%s, CPU", imageTitle)
		plotAvgCpu.X.Label.Text = "second"
		plotAvgCpu.Y.Label.Text = "CPU"
		plotAvgCpu.Legend.Top = true
		if err = plotutil.AddLines(
			plotAvgCpu,
			"etcd3_0", plotAvgCpuEtcd3Points_0,
			"etcd3_1", plotAvgCpuEtcd3Points_1,
		); err != nil {
			return err
		}
		if err = plotAvgCpu.Save(defaultPlotWidth, defaultPlotHeight, avgCpuPath); err != nil {
			return err
		}

		plotAvgMemoryEtcd3_0, err := fr.GetColumn("avg_memory_mb_etcd3_0")
		if err != nil {
			return err
		}
		plotAvgMemoryEtcd3Points_0, err := points(plotAvgMemoryEtcd3_0)
		if err != nil {
			return err
		}
		plotAvgMemoryEtcd3_1, err := fr.GetColumn("avg_memory_mb_etcd3_1")
		if err != nil {
			return err
		}
		plotAvgMemoryEtcd3Points_1, err := points(plotAvgMemoryEtcd3_1)
		if err != nil {
			return err
		}
		plotAvgMemory, err := plot.New()
		if err != nil {
			return err
		}
		plotAvgMemory.Title.Text = fmt.Sprintf("%s, Memory", imageTitle)
		plotAvgMemory.X.Label.Text = "second"
		plotAvgMemory.Y.Label.Text = "Memory(MB)"
		plotAvgMemory.Legend.Top = true
		if err = plotutil.AddLines(
			plotAvgMemory,
			"etcd3_0", plotAvgMemoryEtcd3Points_0,
			"etcd3_1", plotAvgMemoryEtcd3Points_1,
		); err != nil {
			return err
		}
		if err = plotAvgMemory.Save(defaultPlotWidth, defaultPlotHeight, avgMemPath); err != nil {
			return err
		}

		return nil
	}

	plotAvgLatencyConsul, err := fr.GetColumn("avg_latency_ms_consul")
	if err != nil {
		return err
	}
	plotAvgLatencyConsulPoints, err := points(plotAvgLatencyConsul)
	if err != nil {
		return err
	}
	plotAvgLatencyEtcd3, err := fr.GetColumn("avg_latency_ms_etcd3")
	if err != nil {
		return err
	}
	plotAvgLatencyEtcd3Points, err := points(plotAvgLatencyEtcd3)
	if err != nil {
		return err
	}
	plotAvgLatencyEtcd2, err := fr.GetColumn("avg_latency_ms_etcd2")
	if err != nil {
		return err
	}
	plotAvgLatencyEtcd2Points, err := points(plotAvgLatencyEtcd2)
	if err != nil {
		return err
	}
	plotAvgLatency, err := plot.New()
	if err != nil {
		return err
	}
	plotAvgLatencyZk, err := fr.GetColumn("avg_latency_ms_zk")
	if err != nil {
		return err
	}
	plotAvgLatencyZkPoints, err := points(plotAvgLatencyZk)
	if err != nil {
		return err
	}
	plotAvgLatency.Title.Text = fmt.Sprintf("%s, Latency", imageTitle)
	plotAvgLatency.X.Label.Text = "second"
	plotAvgLatency.Y.Label.Text = "Latency(ms)"
	plotAvgLatency.Legend.Top = true
	plotAvgLatencyEtcd3Multi, err := fr.GetColumn("avg_latency_ms_etcd3multi")
	if err == nil {
		var plotAvgLatencyEtcd3MultiPoints plotter.XYs
		plotAvgLatencyEtcd3MultiPoints, err = points(plotAvgLatencyEtcd3Multi)
		if err != nil {
			return err
		}
		if err = plotutil.AddLines(
			plotAvgLatency,
			"consul", plotAvgLatencyConsulPoints,
			"etcd3", plotAvgLatencyEtcd3Points,
			strings.Replace("etcd3multi", "multi", "-"+multiTagTitle, -1), plotAvgLatencyEtcd3MultiPoints,
			"etcd2", plotAvgLatencyEtcd2Points,
			"zk", plotAvgLatencyZkPoints,
		); err != nil {
			return err
		}
		if err = plotAvgLatency.Save(defaultPlotWidth, defaultPlotHeight, avgLatencyPath); err != nil {
			return err
		}
	} else {
		if err = plotutil.AddLines(
			plotAvgLatency,
			"consul", plotAvgLatencyConsulPoints,
			"etcd3", plotAvgLatencyEtcd3Points,
			"etcd2", plotAvgLatencyEtcd2Points,
			"zk", plotAvgLatencyZkPoints,
		); err != nil {
			return err
		}
		if err = plotAvgLatency.Save(defaultPlotWidth, defaultPlotHeight, avgLatencyPath); err != nil {
			return err
		}
	}

	plotThroughputConsul, err := fr.GetColumn("throughput_consul")
	if err != nil {
		return err
	}
	plotThroughputConsulPoints, err := points(plotThroughputConsul)
	if err != nil {
		return err
	}
	plotThroughputEtcd3, err := fr.GetColumn("throughput_etcd3")
	if err != nil {
		return err
	}
	plotThroughputEtcd3Points, err := points(plotThroughputEtcd3)
	if err != nil {
		return err
	}
	plotThroughputEtcd2, err := fr.GetColumn("throughput_etcd2")
	if err != nil {
		return err
	}
	plotThroughputEtcd2Points, err := points(plotThroughputEtcd2)
	if err != nil {
		return err
	}
	plotThroughputZk, err := fr.GetColumn("throughput_zk")
	if err != nil {
		return err
	}
	plotThroughputZkPoints, err := points(plotThroughputZk)
	if err != nil {
		return err
	}
	plotThroughput, err := plot.New()
	if err != nil {
		return err
	}
	plotThroughput.Title.Text = fmt.Sprintf("%s, Throughput", imageTitle)
	plotThroughput.X.Label.Text = "second"
	plotThroughput.Y.Label.Text = "Throughput"
	plotThroughput.Legend.Top = true
	plotThroughputEtcd3Multi, err := fr.GetColumn("throughput_etcd3multi")
	if err == nil {
		var plotThroughputEtcd3MultiPoints plotter.XYs
		plotThroughputEtcd3MultiPoints, err = points(plotThroughputEtcd3Multi)
		if err != nil {
			return err
		}
		if err = plotutil.AddLines(
			plotThroughput,
			"consul", plotThroughputConsulPoints,
			"etcd3", plotThroughputEtcd3Points,
			strings.Replace("etcd3multi", "multi", "-"+multiTagTitle, -1), plotThroughputEtcd3MultiPoints,
			"etcd2", plotThroughputEtcd2Points,
			"zk", plotThroughputZkPoints,
		); err != nil {
			return err
		}
		if err = plotThroughput.Save(defaultPlotWidth, defaultPlotHeight, throughputPath); err != nil {
			return err
		}
	} else {
		if err = plotutil.AddLines(
			plotThroughput,
			"consul", plotThroughputConsulPoints,
			"etcd3", plotThroughputEtcd3Points,
			"etcd2", plotThroughputEtcd2Points,
			"zk", plotThroughputZkPoints,
		); err != nil {
			return err
		}
		if err = plotThroughput.Save(defaultPlotWidth, defaultPlotHeight, throughputPath); err != nil {
			return err
		}
	}

	plotAvgCpuConsul, err := fr.GetColumn("avg_cpu_consul")
	if err != nil {
		return err
	}
	plotAvgCpuConsulPoints, err := points(plotAvgCpuConsul)
	if err != nil {
		return err
	}
	plotAvgCpuEtcd3, err := fr.GetColumn("avg_cpu_etcd3")
	if err != nil {
		return err
	}
	plotAvgCpuEtcd3Points, err := points(plotAvgCpuEtcd3)
	if err != nil {
		return err
	}
	plotAvgCpuEtcd2, err := fr.GetColumn("avg_cpu_etcd2")
	if err != nil {
		return err
	}
	plotAvgCpuEtcd2Points, err := points(plotAvgCpuEtcd2)
	if err != nil {
		return err
	}
	plotAvgCpuZk, err := fr.GetColumn("avg_cpu_zk")
	if err != nil {
		return err
	}
	plotAvgCpuZkPoints, err := points(plotAvgCpuZk)
	if err != nil {
		return err
	}
	plotAvgCpu, err := plot.New()
	if err != nil {
		return err
	}
	plotAvgCpu.Title.Text = fmt.Sprintf("%s, CPU", imageTitle)
	plotAvgCpu.X.Label.Text = "second"
	plotAvgCpu.Y.Label.Text = "CPU"
	plotAvgCpu.Legend.Top = true
	plotAvgCpuEtcd3Multi, err := fr.GetColumn("avg_cpu_etcd3multi")
	if err == nil {
		var plotAvgCpuEtcd3MultiPoints plotter.XYs
		plotAvgCpuEtcd3MultiPoints, err = points(plotAvgCpuEtcd3Multi)
		if err != nil {
			return err
		}
		if err = plotutil.AddLines(
			plotAvgCpu,
			"consul", plotAvgCpuConsulPoints,
			"etcd3", plotAvgCpuEtcd3Points,
			strings.Replace("etcd3multi", "multi", "-"+multiTagTitle, -1), plotAvgCpuEtcd3MultiPoints,
			"etcd2", plotAvgCpuEtcd2Points,
			"zk", plotAvgCpuZkPoints,
		); err != nil {
			return err
		}
		if err = plotAvgCpu.Save(defaultPlotWidth, defaultPlotHeight, avgCpuPath); err != nil {
			return err
		}
	} else {
		if err = plotutil.AddLines(
			plotAvgCpu,
			"consul", plotAvgCpuConsulPoints,
			"etcd3", plotAvgCpuEtcd3Points,
			"etcd2", plotAvgCpuEtcd2Points,
			"zk", plotAvgCpuZkPoints,
		); err != nil {
			return err
		}
		if err = plotAvgCpu.Save(defaultPlotWidth, defaultPlotHeight, avgCpuPath); err != nil {
			return err
		}
	}

	plotAvgMemConsul, err := fr.GetColumn("avg_memory_mb_consul")
	if err != nil {
		return err
	}
	plotAvgMemConsulPoints, err := points(plotAvgMemConsul)
	if err != nil {
		return err
	}
	plotAvgMemEtcd3, err := fr.GetColumn("avg_memory_mb_etcd3")
	if err != nil {
		return err
	}
	plotAvgMemEtcd3Points, err := points(plotAvgMemEtcd3)
	if err != nil {
		return err
	}
	plotAvgMemEtcd2, err := fr.GetColumn("avg_memory_mb_etcd2")
	if err != nil {
		return err
	}
	plotAvgMemEtcd2Points, err := points(plotAvgMemEtcd2)
	if err != nil {
		return err
	}
	plotAvgMemZk, err := fr.GetColumn("avg_memory_mb_zk")
	if err != nil {
		return err
	}
	plotAvgMemZkPoints, err := points(plotAvgMemZk)
	if err != nil {
		return err
	}
	plotAvgMem, err := plot.New()
	if err != nil {
		return err
	}
	plotAvgMem.Title.Text = fmt.Sprintf("%s, Memory", imageTitle)
	plotAvgMem.X.Label.Text = "second"
	plotAvgMem.Y.Label.Text = "Memory(MB)"
	plotAvgMem.Legend.Top = true
	plotAvgMemEtcd3Multi, err := fr.GetColumn("avg_memory_mb_etcd3multi")
	if err == nil {
		var plotAvgMemEtcd3MultiPoints plotter.XYs
		plotAvgMemEtcd3MultiPoints, err = points(plotAvgMemEtcd3Multi)
		if err != nil {
			return err
		}
		if err = plotutil.AddLines(
			plotAvgMem,
			"consul", plotAvgMemConsulPoints,
			"etcd3", plotAvgMemEtcd3Points,
			strings.Replace("etcd3multi", "multi", "-"+multiTagTitle, -1), plotAvgMemEtcd3MultiPoints,
			"etcd2", plotAvgMemEtcd2Points,
			"zk", plotAvgMemZkPoints,
		); err != nil {
			return err
		}
		if err = plotAvgMem.Save(defaultPlotWidth, defaultPlotHeight, avgMemPath); err != nil {
			return err
		}
	} else {
		if err = plotutil.AddLines(
			plotAvgMem,
			"consul", plotAvgMemConsulPoints,
			"etcd3", plotAvgMemEtcd3Points,
			"etcd2", plotAvgMemEtcd2Points,
			"zk", plotAvgMemZkPoints,
		); err != nil {
			return err
		}
		if err = plotAvgMem.Save(defaultPlotWidth, defaultPlotHeight, avgMemPath); err != nil {
			return err
		}
	}

	return nil
}

func points(col dataframe.Column) (plotter.XYs, error) {
	bv, ok := col.BackNonNil()
	if !ok {
		return nil, fmt.Errorf("BackNonNil not found")
	}
	rowN, ok := col.FindValue(bv)
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
