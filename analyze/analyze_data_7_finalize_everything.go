// Copyright 2017 CoreOS, Inc.
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
	"encoding/csv"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"strconv"

	"github.com/coreos/dbtester/control"
	humanize "github.com/dustin/go-humanize"
	"github.com/gyuho/dataframe"
	"github.com/olekukonko/tablewriter"
)

type allAggregatedData struct {
	title          string
	data           []*analyzeData
	databaseTags   []string
	headerToLegend map[string]string
}

func do(configPath string) error {
	cfg, err := readConfig(configPath)
	if err != nil {
		return err
	}

	all := &allAggregatedData{
		title:          cfg.Title,
		data:           make([]*analyzeData, 0, len(cfg.RawData)),
		headerToLegend: make(map[string]string),
	}
	for _, elem := range cfg.RawData {
		plog.Printf("reading system metrics data for %s (%q)", makeTag(elem.Legend), elem.Legend)
		ad, err := readSystemMetricsAll(elem.DataInterpolatedSystemMetricsPaths...)
		if err != nil {
			return err
		}
		ad.databaseTag = makeTag(elem.Legend)
		ad.legend = elem.Legend
		ad.csvOutputpath = elem.OutputPath

		if err = ad.aggSystemMetrics(); err != nil {
			return err
		}
		if err = ad.importBenchMetrics(elem.DataBenchmarkThroughput); err != nil {
			return err
		}
		if err = ad.aggregateAll(elem.DataBenchmarkMemoryByKey, elem.TotalRequests); err != nil {
			return err
		}
		if err = ad.save(); err != nil {
			return err
		}

		all.data = append(all.data, ad)
		all.databaseTags = append(all.databaseTags, makeTag(elem.Legend))
		for _, hd := range ad.aggregated.Headers() {
			all.headerToLegend[makeHeader(hd, makeTag(elem.Legend))] = elem.Legend
		}
	}

	// aggregated everything
	// 1. sum of all network usage per database
	// 2. throughput, latency percentiles distribution
	//
	// FIRST ROW FOR HEADER: etcd, Zookeeper, Consul, ...
	// FIRST COLUMN FOR LABELS: READS-COMPLETED-DELTA, ...
	// SECOND COLUMNS ~ FOR DATA
	row00Header := []string{""} // first is empty
	for _, ad := range all.data {
		// per database
		for _, col := range ad.aggregated.Columns() {
			legend := all.headerToLegend[col.Header()]
			row00Header = append(row00Header, makeTag(legend))
			break
		}
	}
	row01ReadsCompletedDeltaSum := []string{"READS-COMPLETED-DELTA"}
	row02SectorsReadDeltaSum := []string{"SECTORS-READS-DELTA-SUM"}
	row03WritesCompletedDeltaSum := []string{"WRITES-COMPLETED-DELTA-SUM"}
	row04SectorsWrittenDeltaSum := []string{"SECTORS-WRITTEN-DELTA-SUM"}
	row06ReceiveBytesSum := []string{"RECEIVE-BYTES-SUM"}
	row07TransmitBytesSum := []string{"TRANSMIT-BYTES-SUM"}
	row08MaxCPUUsage := []string{"MAX-CPU-USAGE"}
	row09MaxMemoryUsage := []string{"MAX-MEMORY-USAGE"}

	// iterate each database's all data
	for _, ad := range all.data {
		// per database
		var (
			readsCompletedDeltaSum   float64
			sectorsReadDeltaSum      float64
			writesCompletedDeltaSum  float64
			sectorsWrittenDeltaSum   float64
			receiveBytesNumDeltaSum  float64
			transmitBytesNumDeltaSum float64
			maxAvgCPU                float64
			maxAvgVMRSSMBs           []float64
		)
		for _, col := range ad.aggregated.Columns() {
			hdr := col.Header()
			switch {
			case strings.HasPrefix(hdr, "READS-COMPLETED-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Float64()
					readsCompletedDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "SECTORS-READS-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Float64()
					sectorsReadDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "WRITES-COMPLETED-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Float64()
					writesCompletedDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "SECTORS-WRITTEN-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Float64()
					sectorsWrittenDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "RECEIVE-BYTES-NUM-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Float64()
					receiveBytesNumDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "TRANSMIT-BYTES-NUM-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Float64()
					transmitBytesNumDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "AVG-CPU-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Float64()
					maxAvgCPU = maxFloat64(maxAvgCPU, fv)
				}
			case strings.HasPrefix(hdr, "AVG-VMRSS-MB-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Float64()
					maxAvgVMRSSMBs = append(maxAvgVMRSSMBs, fv)
				}
			}
		}

		row01ReadsCompletedDeltaSum = append(row01ReadsCompletedDeltaSum, fmt.Sprintf("%d", uint64(readsCompletedDeltaSum)))
		row02SectorsReadDeltaSum = append(row02SectorsReadDeltaSum, fmt.Sprintf("%d", uint64(sectorsReadDeltaSum)))
		row03WritesCompletedDeltaSum = append(row03WritesCompletedDeltaSum, fmt.Sprintf("%d", uint64(writesCompletedDeltaSum)))
		row04SectorsWrittenDeltaSum = append(row04SectorsWrittenDeltaSum, fmt.Sprintf("%d", uint64(sectorsWrittenDeltaSum)))
		row06ReceiveBytesSum = append(row06ReceiveBytesSum, humanize.Bytes(uint64(receiveBytesNumDeltaSum)))
		row07TransmitBytesSum = append(row07TransmitBytesSum, humanize.Bytes(uint64(transmitBytesNumDeltaSum)))
		row08MaxCPUUsage = append(row08MaxCPUUsage, fmt.Sprintf("%.2f %%", maxAvgCPU))

		// TODO: linux sometimes returns overflowed value...
		sort.Float64s(maxAvgVMRSSMBs)
		row09MaxMemoryUsage = append(row09MaxMemoryUsage, fmt.Sprintf("%.2f MB", maxAvgVMRSSMBs[len(maxAvgVMRSSMBs)-2]))
	}

	row05AverageDatasize := []string{"AVG-DATA-SIZE"}    // TOTAL-DATA-SIZE
	row10TotalSeconds := []string{"TOTAL-SECONDS"}       // TOTAL-SECONDS
	row11MinThroughput := []string{"MIN-THROUGHPUT"}     // MIN AVG-THROUGHPUT
	row12AverageThroughput := []string{"AVG-THROUGHPUT"} // REQUESTS-PER-SECOND
	row13MaxThroughput := []string{"MAX-THROUGHPUT"}     // MAX AVG-THROUGHPUT
	row14SlowestLatency := []string{"SLOWEST-LATENCY"}   // SLOWEST-LATENCY-MS
	row15AverageLatency := []string{"AVG-LATENCY"}       // AVERAGE-LATENCY-MS
	row16FastestLatency := []string{"FASTEST-LATENCY"}   // FASTEST-LATENCY-MS
	row17p10 := []string{"Latency p10"}                  // p10
	row18p25 := []string{"Latency p25"}                  // p25
	row19p50 := []string{"Latency p50"}                  // p50
	row20p75 := []string{"Latency p75"}                  // p75
	row21p90 := []string{"Latency p90"}                  // p90
	row22p95 := []string{"Latency p95"}                  // p95
	row23p99 := []string{"Latency p99"}                  // p99
	row24p999 := []string{"Latency p99.9"}               // p99.9

	for i, rcfg := range cfg.RawData {
		tag := makeTag(rcfg.Legend)
		if tag != row00Header[i+1] {
			return fmt.Errorf("analyze config has different order; expected %q, got %q", row00Header[i+1], tag)
		}

		{
			f, err := openToRead(rcfg.DataBenchmarkLatencySummary)
			if err != nil {
				return err
			}
			defer f.Close()

			rd := csv.NewReader(f)

			// FieldsPerRecord is the number of expected fields per record.
			// If FieldsPerRecord is positive, Read requires each record to
			// have the given number of fields. If FieldsPerRecord is 0, Read sets it to
			// the number of fields in the first record, so that future records must
			// have the same field count. If FieldsPerRecord is negative, no check is
			// made and records may have a variable number of fields.
			rd.FieldsPerRecord = -1

			rows, err := rd.ReadAll()
			if err != nil {
				return err
			}

			for _, row := range rows {
				switch row[0] {
				case "TOTAL-SECONDS":
					row10TotalSeconds = append(row10TotalSeconds, fmt.Sprintf("%s sec", row[1]))
				case "REQUESTS-PER-SECOND":
					fv, err := strconv.ParseFloat(row[1], 64)
					if err != nil {
						return err
					}
					avg := int64(fv)
					row12AverageThroughput = append(row12AverageThroughput, fmt.Sprintf("%s req/sec", humanize.Comma(avg)))
				case "SLOWEST-LATENCY-MS":
					row14SlowestLatency = append(row14SlowestLatency, fmt.Sprintf("%s ms", row[1]))
				case "FASTEST-LATENCY-MS":
					row16FastestLatency = append(row16FastestLatency, fmt.Sprintf("%s ms", row[1]))
				case "AVERAGE-LATENCY-MS":
					row15AverageLatency = append(row15AverageLatency, fmt.Sprintf("%s ms", row[1]))
				}
			}
		}

		{
			fr, err := dataframe.NewFromCSV(nil, rcfg.DataBenchmarkThroughput)
			if err != nil {
				return err
			}
			col, err := fr.Column("AVG-THROUGHPUT")
			if err != nil {
				return err
			}
			var min int64
			var max int64
			for i := 0; i < col.Count(); i++ {
				val, err := col.Value(i)
				if err != nil {
					return err
				}
				fv, _ := val.Float64()

				if i == 0 {
					min = int64(fv)
				}
				if max < int64(fv) {
					max = int64(fv)
				}
				if min > int64(fv) {
					min = int64(fv)
				}
			}
			row11MinThroughput = append(row11MinThroughput, fmt.Sprintf("%s req/sec", humanize.Comma(min)))
			row13MaxThroughput = append(row13MaxThroughput, fmt.Sprintf("%s req/sec", humanize.Comma(max)))
		}

		{
			fr, err := dataframe.NewFromCSV(nil, rcfg.DatasizeSummary)
			if err != nil {
				return err
			}
			col, err := fr.Column(control.DataSummaryColumns[3]) // datasize in bytes
			if err != nil {
				return err
			}
			var sum float64
			for i := 0; i < col.Count(); i++ {
				val, err := col.Value(i)
				if err != nil {
					return err
				}
				fv, _ := val.Float64()
				sum += fv
			}
			avg := uint64(sum / float64(col.Count()))
			row05AverageDatasize = append(row05AverageDatasize, humanize.Bytes(avg))
		}

		{
			f, err := openToRead(rcfg.DataBenchmarkLatencyPercentile)
			if err != nil {
				return err
			}
			defer f.Close()

			rd := csv.NewReader(f)

			// FieldsPerRecord is the number of expected fields per record.
			// If FieldsPerRecord is positive, Read requires each record to
			// have the given number of fields. If FieldsPerRecord is 0, Read sets it to
			// the number of fields in the first record, so that future records must
			// have the same field count. If FieldsPerRecord is negative, no check is
			// made and records may have a variable number of fields.
			rd.FieldsPerRecord = -1

			rows, err := rd.ReadAll()
			if err != nil {
				return err
			}

			for ri, row := range rows {
				if ri == 0 {
					continue // skip header
				}
				switch row[0] {
				case "p10":
					row17p10 = append(row17p10, fmt.Sprintf("%s ms", row[1]))
				case "p25":
					row18p25 = append(row18p25, fmt.Sprintf("%s ms", row[1]))
				case "p50":
					row19p50 = append(row19p50, fmt.Sprintf("%s ms", row[1]))
				case "p75":
					row20p75 = append(row20p75, fmt.Sprintf("%s ms", row[1]))
				case "p90":
					row21p90 = append(row21p90, fmt.Sprintf("%s ms", row[1]))
				case "p95":
					row22p95 = append(row22p95, fmt.Sprintf("%s ms", row[1]))
				case "p99":
					row23p99 = append(row23p99, fmt.Sprintf("%s ms", row[1]))
				case "p99.9":
					row24p999 = append(row24p999, fmt.Sprintf("%s ms", row[1]))
				}
			}
		}
	}

	aggRows := [][]string{
		row00Header,
		row01ReadsCompletedDeltaSum,
		row02SectorsReadDeltaSum,
		row03WritesCompletedDeltaSum,
		row04SectorsWrittenDeltaSum,
		row05AverageDatasize,
		row06ReceiveBytesSum,
		row07TransmitBytesSum,
		row08MaxCPUUsage,
		row09MaxMemoryUsage,
		row10TotalSeconds,
		row11MinThroughput,
		row12AverageThroughput,
		row13MaxThroughput,
		row14SlowestLatency,
		row15AverageLatency,
		row16FastestLatency,
		row17p10,
		row18p25,
		row19p50,
		row20p75,
		row21p90,
		row22p95,
		row23p99,
		row24p999,
	}
	plog.Printf("saving data to %q", cfg.AllAggregatedPath)
	file, err := openToOverwrite(cfg.AllAggregatedPath)
	if err != nil {
		return err
	}
	defer file.Close()
	wr := csv.NewWriter(file)
	if err := wr.WriteAll(aggRows); err != nil {
		return err
	}
	wr.Flush()
	if err := wr.Error(); err != nil {
		return err
	}

	// KEYS, AVG-LATENCY-MS
	plog.Printf("combining data to %q", cfg.AllLatencyByKey)
	allLatencyFrame := dataframe.New()
	for _, elem := range cfg.RawData {
		fr, err := dataframe.NewFromCSV(nil, elem.DataBenchmarkLatencyByKey)
		if err != nil {
			return err
		}
		colKeys, err := fr.Column("KEYS")
		if err != nil {
			return err
		}
		colKeys.UpdateHeader(makeHeader("KEYS", makeTag(elem.Legend)))
		if err = allLatencyFrame.AddColumn(colKeys); err != nil {
			return err
		}

		colLatency, err := fr.Column("AVG-LATENCY-MS")
		if err != nil {
			return err
		}
		colLatency.UpdateHeader(makeHeader("AVG-LATENCY-MS", makeTag(elem.Legend)))
		if err = allLatencyFrame.AddColumn(colLatency); err != nil {
			return err
		}
	}
	if err := allLatencyFrame.CSV(cfg.AllLatencyByKey); err != nil {
		return err
	}

	allLatencyFrameCfg := PlotConfig{
		Column:         "AVG-LATENCY-MS",
		XAxis:          "Keys",
		YAxis:          "Latency(millisecond)",
		OutputPathList: make([]string, len(cfg.PlotList[0].OutputPathList)),
	}
	allLatencyFrameCfg.OutputPathList[0] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-LATENCY-MS-BY-KEY.svg")
	allLatencyFrameCfg.OutputPathList[1] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-LATENCY-MS-BY-KEY.png")
	plog.Printf("plotting %v", allLatencyFrameCfg.OutputPathList)
	if err = all.draw(allLatencyFrameCfg, allLatencyFrame.Columns()...); err != nil {
		return err
	}

	// KEYS, AVG-VMRSS-MB
	plog.Printf("combining data to %q", cfg.AllMemoryByKey)
	allMemoryFrame := dataframe.New()
	for _, elem := range cfg.RawData {
		fr, err := dataframe.NewFromCSV(nil, elem.DataBenchmarkMemoryByKey)
		if err != nil {
			return err
		}
		colKeys, err := fr.Column("KEYS")
		if err != nil {
			return err
		}
		colKeys.UpdateHeader(makeHeader("KEYS", makeTag(elem.Legend)))
		if err = allMemoryFrame.AddColumn(colKeys); err != nil {
			return err
		}

		colMem, err := fr.Column("AVG-VMRSS-MB")
		if err != nil {
			return err
		}
		colMem.UpdateHeader(makeHeader("AVG-LATENCY-MS", makeTag(elem.Legend)))
		if err = allMemoryFrame.AddColumn(colMem); err != nil {
			return err
		}
	}
	if err := allMemoryFrame.CSV(cfg.AllMemoryByKey); err != nil {
		return err
	}

	allMemoryFrameCfg := PlotConfig{
		Column:         "AVG-VMRSS-MB",
		XAxis:          "Keys",
		YAxis:          "Memory(MB)",
		OutputPathList: make([]string, len(cfg.PlotList[0].OutputPathList)),
	}
	allMemoryFrameCfg.OutputPathList[0] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-VMRSS-MB-BY-KEY.svg")
	allMemoryFrameCfg.OutputPathList[1] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-VMRSS-MB-BY-KEY.png")
	plog.Printf("plotting %v", allMemoryFrameCfg.OutputPathList)
	if err = all.draw(allMemoryFrameCfg, allMemoryFrame.Columns()...); err != nil {
		return err
	}

	plog.Println("combining data for plotting")
	for _, plotConfig := range cfg.PlotList {
		plog.Printf("plotting %q", plotConfig.Column)
		var clientNumColumns []dataframe.Column
		var dataColumns []dataframe.Column
		for i, ad := range all.data {
			tag := all.databaseTags[i]

			avgCol, err := ad.aggregated.Column("CONTROL-CLIENT-NUM")
			if err != nil {
				return err
			}
			avgCol.UpdateHeader(makeHeader("CONTROL-CLIENT-NUM", tag))
			clientNumColumns = append(clientNumColumns, avgCol)

			col, err := ad.aggregated.Column(plotConfig.Column)
			if err != nil {
				return err
			}
			col.UpdateHeader(makeHeader(plotConfig.Column, tag))
			dataColumns = append(dataColumns, col)
		}
		if err = all.draw(plotConfig, dataColumns...); err != nil {
			return err
		}

		plog.Printf("saving data for %q of all database", plotConfig.Column)
		nf1, err := dataframe.NewFromColumns(nil, dataColumns...)
		if err != nil {
			return err
		}
		if err = nf1.CSV(filepath.Join(cfg.WorkDir, plotConfig.Column+".csv")); err != nil {
			return err
		}

		plog.Printf("saving data for %q of all database (by client number)", plotConfig.Column)
		nf2 := dataframe.New()
		for i := range clientNumColumns {
			if clientNumColumns[i].Count() != dataColumns[i].Count() {
				return fmt.Errorf("%q row count %d != %q row count %d",
					clientNumColumns[i].Header(),
					clientNumColumns[i].Count(),
					dataColumns[i].Header(),
					dataColumns[i].Count(),
				)
			}
			if err := nf2.AddColumn(clientNumColumns[i]); err != nil {
				return err
			}
			if err := nf2.AddColumn(dataColumns[i]); err != nil {
				return err
			}
		}
		if err = nf2.CSV(filepath.Join(cfg.WorkDir, plotConfig.Column+"-BY-CLIENT-NUM"+".csv")); err != nil {
			return err
		}

		plog.Printf("aggregating data for %q of all database (by client number)", plotConfig.Column)
		nf3 := dataframe.New()
		var firstKeys []int
		for i := range clientNumColumns {
			n := clientNumColumns[i].Count()
			allData := make(map[int]float64)
			for j := 0; j < n; j++ {
				v1, err := clientNumColumns[i].Value(j)
				if err != nil {
					return err
				}
				num, _ := v1.Int64()

				v2, err := dataColumns[i].Value(j)
				if err != nil {
					return err
				}
				data, _ := v2.Float64()

				if v, ok := allData[int(num)]; ok {
					allData[int(num)] = (v + data) / 2
				} else {
					allData[int(num)] = data
				}
			}
			var allKeys []int
			for k := range allData {
				allKeys = append(allKeys, k)
			}
			sort.Ints(allKeys)

			if i == 0 {
				firstKeys = allKeys
			}
			if !reflect.DeepEqual(firstKeys, allKeys) {
				return fmt.Errorf("all keys must be %+v, got %+v", firstKeys, allKeys)
			}

			if i == 0 {
				col1 := dataframe.NewColumn("CONTROL-CLIENT-NUM")
				for j := range allKeys {
					col1.PushBack(dataframe.NewStringValue(allKeys[j]))
				}
				if err := nf3.AddColumn(col1); err != nil {
					return err
				}
			}
			col2 := dataframe.NewColumn(dataColumns[i].Header())
			for j := range allKeys {
				col2.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.4f", allData[allKeys[j]])))
			}
			if err := nf3.AddColumn(col2); err != nil {
				return err
			}
		}
		if err = nf3.CSV(filepath.Join(cfg.WorkDir, plotConfig.Column+"-BY-CLIENT-NUM-aggregated"+".csv")); err != nil {
			return err
		}
	}
	buf := new(bytes.Buffer)
	tw := tablewriter.NewWriter(buf)
	tw.SetHeader(aggRows[0])
	for _, row := range aggRows[1:] {
		tw.Append(row)
	}
	tw.SetAutoFormatHeaders(false)
	tw.SetAlignment(tablewriter.ALIGN_RIGHT)
	tw.Render()
	if err := toFile(buf.String(), changeExtToTxt(cfg.AllAggregatedPath)); err != nil {
		return err
	}

	plog.Printf("writing README at %q", cfg.READMEConfig.OutputPath)
	return writeREADME(buf.String(), cfg.READMEConfig)
}

func changeExtToTxt(fpath string) string {
	ext := filepath.Ext(fpath)
	return strings.Replace(fpath, ext, ".txt", -1)
}
