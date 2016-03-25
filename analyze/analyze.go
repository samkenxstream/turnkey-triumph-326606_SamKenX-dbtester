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
		if err := plotAggAgg(globalFlags.AggAggFilePath, globalFlags.OutputPath, globalFlags.ImageFormat, globalFlags.ImageTitle); err != nil {
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
	fr2, err := aggMonitor(monitorPaths...)
	if err != nil {
		return nil, err
	}
	col2, err := fr2.GetColumn("unix_ts")
	if err != nil {
		return nil, err
	}

	// need to combine fr2 to fr1
	fr1, err := dataframe.NewFromCSV(nil, benchPath)
	if err != nil {
		return nil, err
	}
	col1, err := fr1.GetColumn("unix_ts")
	if err != nil {
		return nil, err
	}
	fv, ok := col1.FrontNonNil()
	if !ok {
		return nil, fmt.Errorf("FrontNonNil %s has empty Unix time %v", benchPath, fv)
	}
	startIdx, ok := col2.FindValue(fv)
	if !ok {
		return nil, fmt.Errorf("%v is not found in benchmark result %s", fv, benchPath)
	}
	bv, ok := col1.BackNonNil()
	if !ok {
		return nil, fmt.Errorf("BackNonNil %s has empty Unix time %v", benchPath, bv)
	}
	endIdx, ok := col2.FindValue(bv)
	if !ok {
		return nil, fmt.Errorf("%v is not found in benchmark result %s", bv, benchPath)
	}

	var minLen int
	for i, hd := range fr1.GetHeader() {
		var col dataframe.Column
		col, err = fr1.GetColumn(hd)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			minLen = col.RowNumber()
		}
		if minLen < col.RowNumber() {
			minLen = col.RowNumber()
		}
	}
	var delta int
	if minLen > endIdx+1-startIdx { // short of rows
		delta = minLen - (endIdx + 1 - startIdx)
	}

	for _, hd := range fr2.GetHeader() {
		if hd == "unix_ts" {
			continue
		}
		var col dataframe.Column
		col, err = fr2.GetColumn(hd)
		if err != nil {
			return nil, err
		}
		if err = col.KeepRows(startIdx, endIdx+1+delta); err != nil {
			return nil, err
		}
		if err = fr1.AddColumn(col); err != nil {
			return nil, err
		}
	}

	// get average value
	uc, err := fr1.GetColumn("unix_ts")
	if err != nil {
		return nil, err
	}
	var (
		rowNum                  = uc.RowNumber()
		sampleSize              = float64(len(monitorPaths))
		cumulativeThroughputCol = dataframe.NewColumn("cumulative_throughput")
		totalThrougput          int
		avgCpuCol               = dataframe.NewColumn("avg_cpu")
		avgMemCol               = dataframe.NewColumn("avg_memory_mb")
	)
	for i := 0; i < rowNum; i++ {
		var (
			cpuTotal    float64
			memoryTotal float64
		)
		for _, hd := range fr1.GetHeader() {
			var col dataframe.Column
			col, err = fr1.GetColumn(hd)
			if err != nil {
				return nil, err
			}
			var rv dataframe.Value
			rv, err = col.GetValue(i)
			if err != nil {
				return nil, err
			}
			fv, _ := rv.ToNumber()
			switch {
			case strings.HasPrefix(hd, "cpu_"):
				cpuTotal += fv
			case strings.HasPrefix(hd, "memory_"):
				memoryTotal += fv
			case hd == "throughput":
				fv, _ := rv.ToNumber()
				totalThrougput += int(fv)
				cumulativeThroughputCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%d", totalThrougput)))
			}
		}
		avgCpuCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", cpuTotal/sampleSize)))
		avgMemCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", memoryTotal/sampleSize)))
	}

	unixTsCol, err := fr1.GetColumn("unix_ts")
	if err != nil {
		return nil, err
	}
	latencyCol, err := fr1.GetColumn("avg_latency_ms")
	if err != nil {
		return nil, err
	}
	throughputCol, err := fr1.GetColumn("throughput")
	if err != nil {
		return nil, err
	}

	nf := dataframe.New()
	nf.AddColumn(unixTsCol)
	nf.AddColumn(latencyCol)
	nf.AddColumn(throughputCol)
	nf.AddColumn(cumulativeThroughputCol)

	for _, hd := range fr1.GetHeader() {
		col, err := fr1.GetColumn(hd)
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
		secondCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%d", i)))
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
			case strings.Contains(fpaths[i], "etcd"):
				dbID = "etcd3"
			case strings.Contains(fpaths[i], "zk"):
				dbID = "zk"
			}

			col.UpdateHeader(fmt.Sprintf("%s_%s", col.GetHeader(), dbID))
			nf.AddColumn(col)
		}
	}

	return nf, nil
}

func plotAggAgg(fpath, outputPath, imageFormat, imageTitle string) error {
	fr, err := dataframe.NewFromCSV(nil, fpath)
	if err != nil {
		return err
	}

	plot.DefaultFont = "Helvetica"
	plotter.DefaultLineStyle.Width /= 2
	plotter.DefaultGlyphStyle.Radius = vg.Points(2.0)
	var (
		defaultSize    = 5 * vg.Inch
		avgLatencyPath = outputPath + fmt.Sprintf("-avg-latency-ms.%s", imageFormat)
		throughputPath = outputPath + fmt.Sprintf("-throughput.%s", imageFormat)
		avgCpuPath     = outputPath + fmt.Sprintf("-avg-cpu.%s", imageFormat)
		avgMemPath     = outputPath + fmt.Sprintf("-avg-mem.%s", imageFormat)
	)

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
	plotAvgLatency.Y.Label.Text = "latency(ms)"
	if err := plotutil.AddLinePoints(
		plotAvgLatency,
		"consul", plotAvgLatencyConsulPoints,
		"etcd3", plotAvgLatencyEtcd3Points,
		"etcd2", plotAvgLatencyEtcd2Points,
		"zk", plotAvgLatencyZkPoints,
	); err != nil {
		return err
	}
	if err := plotAvgLatency.Save(defaultSize, defaultSize, avgLatencyPath); err != nil {
		return err
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
	if err := plotutil.AddLinePoints(
		plotThroughput,
		"consul", plotThroughputConsulPoints,
		"etcd3", plotThroughputEtcd3Points,
		"etcd2", plotThroughputEtcd2Points,
		"zk", plotThroughputZkPoints,
	); err != nil {
		return err
	}
	if err := plotThroughput.Save(defaultSize, defaultSize, throughputPath); err != nil {
		return err
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
	if err := plotutil.AddLinePoints(
		plotAvgCpu,
		"consul", plotAvgCpuConsulPoints,
		"etcd3", plotAvgCpuEtcd3Points,
		"etcd2", plotAvgCpuEtcd2Points,
		"zk", plotAvgCpuZkPoints,
	); err != nil {
		return err
	}
	if err := plotAvgCpu.Save(defaultSize, defaultSize, avgCpuPath); err != nil {
		return err
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
	if err := plotutil.AddLinePoints(
		plotAvgMem,
		"consul", plotAvgMemConsulPoints,
		"etcd3", plotAvgMemEtcd3Points,
		"etcd2", plotAvgMemEtcd2Points,
		"zk", plotAvgMemZkPoints,
	); err != nil {
		return err
	}
	if err := plotAvgMem.Save(defaultSize, defaultSize, avgMemPath); err != nil {
		return err
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
