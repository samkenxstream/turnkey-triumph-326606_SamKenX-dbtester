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
	row01ReadsCompletedDeltaSum := []string{"READS-COMPLETED-DELTA-SUM"}
	row02SectorsReadDeltaSum := []string{"SECTORS-READS-DELTA-SUM"}
	row03WritesCompletedDeltaSum := []string{"WRITES-COMPLETED-DELTA-SUM"}
	row04SectorsWrittenDeltaSum := []string{"SECTORS-WRITTEN-DELTA-SUM"}
	row06ReceiveBytesSum := []string{"NETWORK-RECEIVE-DATA-SUM"}
	row07TransmitBytesSum := []string{"NETWORK-TRANSMIT-DATA-SUM"}
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

		row01ReadsCompletedDeltaSum = append(row01ReadsCompletedDeltaSum, humanize.Comma(int64(readsCompletedDeltaSum)))
		row02SectorsReadDeltaSum = append(row02SectorsReadDeltaSum, humanize.Comma(int64(sectorsReadDeltaSum)))
		row03WritesCompletedDeltaSum = append(row03WritesCompletedDeltaSum, humanize.Comma(int64(writesCompletedDeltaSum)))
		row04SectorsWrittenDeltaSum = append(row04SectorsWrittenDeltaSum, humanize.Comma(int64(sectorsWrittenDeltaSum)))
		row06ReceiveBytesSum = append(row06ReceiveBytesSum, humanize.Bytes(uint64(receiveBytesNumDeltaSum)))
		row07TransmitBytesSum = append(row07TransmitBytesSum, humanize.Bytes(uint64(transmitBytesNumDeltaSum)))
		row08MaxCPUUsage = append(row08MaxCPUUsage, fmt.Sprintf("%.2f %%", maxAvgCPU))

		// TODO: handle overflowed memory value?
		sort.Float64s(maxAvgVMRSSMBs)
		mv := maxAvgVMRSSMBs[len(maxAvgVMRSSMBs)-1]
		mb := uint64(mv * 1000000)
		row09MaxMemoryUsage = append(row09MaxMemoryUsage, humanize.Bytes(mb))
	}

	row05AverageDatasize := []string{"AVG-DATA-SIZE-ON-DISK"}              // TOTAL-DATA-SIZE
	row10TotalSeconds := []string{"TOTAL-SECONDS"}                         // TOTAL-SECONDS
	row11MaxThroughput := []string{"MAX-THROUGHPUT"}                       // MAX AVG-THROUGHPUT
	row12AverageThroughput := []string{"AVG-THROUGHPUT"}                   // REQUESTS-PER-SECOND
	row13MinThroughput := []string{"MIN-THROUGHPUT"}                       // MIN AVG-THROUGHPUT
	row14FastestLatency := []string{"FASTEST-LATENCY"}                     // FASTEST-LATENCY-MS
	row15AverageLatency := []string{"AVG-LATENCY"}                         // AVERAGE-LATENCY-MS
	row16SlowestLatency := []string{"SLOWEST-LATENCY"}                     // SLOWEST-LATENCY-MS
	row17p10 := []string{"Latency p10"}                                    // p10
	row18p25 := []string{"Latency p25"}                                    // p25
	row19p50 := []string{"Latency p50"}                                    // p50
	row20p75 := []string{"Latency p75"}                                    // p75
	row21p90 := []string{"Latency p90"}                                    // p90
	row22p95 := []string{"Latency p95"}                                    // p95
	row23p99 := []string{"Latency p99"}                                    // p99
	row24p999 := []string{"Latency p99.9"}                                 // p99.9
	row25ClientReceiveBytesSum := []string{"CLIENT-NETWORK-RECEIVE-SUM"}   // RECEIVE-BYTES-NUM-DELTA
	row26ClientTransmitBytesSum := []string{"CLIENT-NETWORK-TRANSMIT-SUM"} // TRANSMIT-BYTES-DELTA
	row27ClientMaxCPU := []string{"CLIENT-MAX-CPU-USAGE"}                  // CPU-NUM
	row28ClientMaxMemory := []string{"CLIENT-MAX-MEMORY-USAGE"}            // VMRSS-NUM
	row29ErrorCount := []string{"CLIENT-ERROR-COUNT"}                      // ERROR:

	for i, rcfg := range cfg.RawData {
		tag := makeTag(rcfg.Legend)
		if tag != row00Header[i+1] {
			return fmt.Errorf("analyze config has different order; expected %q, got %q", row00Header[i+1], tag)
		}

		{
			fr, err := dataframe.NewFromCSV(nil, rcfg.ClientSystemMetricsInterpolated)
			if err != nil {
				return err
			}

			var receiveBytesNumDeltaSum uint64
			col, err := fr.Column("RECEIVE-BYTES-NUM-DELTA")
			if err != nil {
				return err
			}
			for i := 0; i < col.Count(); i++ {
				v, err := col.Value(i)
				if err != nil {
					return err
				}
				iv, _ := v.Uint64()
				receiveBytesNumDeltaSum += iv
			}

			var transmitBytesNumDeltaSum uint64
			col, err = fr.Column("TRANSMIT-BYTES-NUM-DELTA")
			if err != nil {
				return err
			}
			for i := 0; i < col.Count(); i++ {
				v, err := col.Value(i)
				if err != nil {
					return err
				}
				iv, _ := v.Uint64()
				transmitBytesNumDeltaSum += iv
			}

			var maxAvgCPU float64
			col, err = fr.Column("CPU-NUM")
			if err != nil {
				return err
			}
			for i := 0; i < col.Count(); i++ {
				v, err := col.Value(i)
				if err != nil {
					return err
				}
				fv, _ := v.Float64()
				if maxAvgCPU == 0 || fv > maxAvgCPU {
					maxAvgCPU = fv
				}
			}

			var maxVMRSSNum uint64
			col, err = fr.Column("VMRSS-NUM")
			if err != nil {
				return err
			}
			for i := 0; i < col.Count(); i++ {
				v, err := col.Value(i)
				if err != nil {
					return err
				}
				iv, _ := v.Uint64()
				if maxVMRSSNum == 0 || iv > maxVMRSSNum {
					maxVMRSSNum = iv
				}
			}

			row25ClientReceiveBytesSum = append(row25ClientReceiveBytesSum, humanize.Bytes(receiveBytesNumDeltaSum))
			row26ClientTransmitBytesSum = append(row26ClientTransmitBytesSum, humanize.Bytes(transmitBytesNumDeltaSum))
			row27ClientMaxCPU = append(row27ClientMaxCPU, fmt.Sprintf("%.2f %%", maxAvgCPU))
			row28ClientMaxMemory = append(row28ClientMaxMemory, humanize.Bytes(maxVMRSSNum))
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

			var totalErrCnt int64
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
					row16SlowestLatency = append(row16SlowestLatency, fmt.Sprintf("%s ms", row[1]))
				case "FASTEST-LATENCY-MS":
					row14FastestLatency = append(row14FastestLatency, fmt.Sprintf("%s ms", row[1]))
				case "AVERAGE-LATENCY-MS":
					row15AverageLatency = append(row15AverageLatency, fmt.Sprintf("%s ms", row[1]))
				}

				if strings.HasPrefix(row[0], "ERROR:") {
					iv, err := strconv.ParseInt(row[1], 10, 64)
					if err != nil {
						return err
					}
					totalErrCnt += iv
				}
			}
			row29ErrorCount = append(row29ErrorCount, humanize.Comma(totalErrCnt))
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
			row11MaxThroughput = append(row11MaxThroughput, fmt.Sprintf("%s req/sec", humanize.Comma(max)))
			row13MinThroughput = append(row13MinThroughput, fmt.Sprintf("%s req/sec", humanize.Comma(min)))
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
		row11MaxThroughput,
		row12AverageThroughput,
		row13MinThroughput,
		row14FastestLatency,
		row15AverageLatency,
		row16SlowestLatency,
		row17p10,
		row18p25,
		row19p50,
		row20p75,
		row21p90,
		row22p95,
		row23p99,
		row24p999,
		row25ClientReceiveBytesSum,
		row26ClientTransmitBytesSum,
		row27ClientMaxCPU,
		row28ClientMaxMemory,
		row29ErrorCount,
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

	// KEYS, MIN-LATENCY-MS, AVG-LATENCY-MS, MAX-LATENCY-MS
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

		colMinLatency, err := fr.Column("MIN-LATENCY-MS")
		if err != nil {
			return err
		}
		colMinLatency.UpdateHeader(makeHeader("MIN-LATENCY-MS", makeTag(elem.Legend)))
		if err = allLatencyFrame.AddColumn(colMinLatency); err != nil {
			return err
		}

		colAvgLatency, err := fr.Column("AVG-LATENCY-MS")
		if err != nil {
			return err
		}
		colAvgLatency.UpdateHeader(makeHeader("AVG-LATENCY-MS", makeTag(elem.Legend)))
		if err = allLatencyFrame.AddColumn(colAvgLatency); err != nil {
			return err
		}

		colMaxLatency, err := fr.Column("MAX-LATENCY-MS")
		if err != nil {
			return err
		}
		colMaxLatency.UpdateHeader(makeHeader("MAX-LATENCY-MS", makeTag(elem.Legend)))
		if err = allLatencyFrame.AddColumn(colMaxLatency); err != nil {
			return err
		}
	}
	if err := allLatencyFrame.CSV(cfg.AllLatencyByKey); err != nil {
		return err
	}
	{
		allLatencyFrameCfg := PlotConfig{
			Column:         "AVG-LATENCY-MS",
			XAxis:          "Cumulative Number of Keys",
			YAxis:          "Latency(millisecond) by Keys",
			OutputPathList: make([]string, len(cfg.PlotList[0].OutputPathList)),
		}
		allLatencyFrameCfg.OutputPathList[0] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-LATENCY-MS-BY-KEY.svg")
		allLatencyFrameCfg.OutputPathList[1] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-LATENCY-MS-BY-KEY.png")
		plog.Printf("plotting %v", allLatencyFrameCfg.OutputPathList)
		var pairs []pair
		allCols := allLatencyFrame.Columns()
		for i := 0; i < len(allCols)-3; i += 4 {
			pairs = append(pairs, pair{
				x: allCols[i],   // x
				y: allCols[i+2], // avg
			})
		}
		if err = all.drawXY(allLatencyFrameCfg, pairs...); err != nil {
			return err
		}
		newCSV := dataframe.New()
		for _, p := range pairs {
			if err = newCSV.AddColumn(p.x); err != nil {
				return err
			}
			if err = newCSV.AddColumn(p.y); err != nil {
				return err
			}
		}
		csvPath := filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-LATENCY-MS-BY-KEY.csv")
		if err := newCSV.CSV(csvPath); err != nil {
			return err
		}
	}
	{
		// with error points
		allLatencyFrameCfg := PlotConfig{
			Column:         "AVG-LATENCY-MS",
			XAxis:          "Cumulative Number of Keys",
			YAxis:          "Latency(millisecond) by Keys",
			OutputPathList: make([]string, len(cfg.PlotList[0].OutputPathList)),
		}
		allLatencyFrameCfg.OutputPathList[0] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg")
		allLatencyFrameCfg.OutputPathList[1] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.png")
		plog.Printf("plotting %v", allLatencyFrameCfg.OutputPathList)
		var triplets []triplet
		allCols := allLatencyFrame.Columns()
		for i := 0; i < len(allCols)-3; i += 4 {
			triplets = append(triplets, triplet{
				x:      allCols[i],
				minCol: allCols[i+1],
				avgCol: allCols[i+2],
				maxCol: allCols[i+3],
			})
		}
		if err = all.drawXYWithErrorPoints(allLatencyFrameCfg, triplets...); err != nil {
			return err
		}
		newCSV := dataframe.New()
		for _, tri := range triplets {
			if err = newCSV.AddColumn(tri.x); err != nil {
				return err
			}
			if err = newCSV.AddColumn(tri.minCol); err != nil {
				return err
			}
			if err = newCSV.AddColumn(tri.avgCol); err != nil {
				return err
			}
			if err = newCSV.AddColumn(tri.maxCol); err != nil {
				return err
			}
		}
		csvPath := filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.csv")
		if err := newCSV.CSV(csvPath); err != nil {
			return err
		}
	}

	// KEYS, MIN-VMRSS-MB, AVG-VMRSS-MB, MAX-VMRSS-MB
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

		colMemMin, err := fr.Column("MIN-VMRSS-MB")
		if err != nil {
			return err
		}
		colMemMin.UpdateHeader(makeHeader("MIN-VMRSS-MB", makeTag(elem.Legend)))
		if err = allMemoryFrame.AddColumn(colMemMin); err != nil {
			return err
		}

		colMem, err := fr.Column("AVG-VMRSS-MB")
		if err != nil {
			return err
		}
		colMem.UpdateHeader(makeHeader("AVG-VMRSS-MB", makeTag(elem.Legend)))
		if err = allMemoryFrame.AddColumn(colMem); err != nil {
			return err
		}

		colMemMax, err := fr.Column("MAX-VMRSS-MB")
		if err != nil {
			return err
		}
		colMemMax.UpdateHeader(makeHeader("MAX-VMRSS-MB", makeTag(elem.Legend)))
		if err = allMemoryFrame.AddColumn(colMemMax); err != nil {
			return err
		}
	}
	if err := allMemoryFrame.CSV(cfg.AllMemoryByKey); err != nil {
		return err
	}
	{
		allMemoryFrameCfg := PlotConfig{
			Column:         "AVG-VMRSS-MB",
			XAxis:          "Cumulative Number of Keys",
			YAxis:          "Memory(MB) by Keys",
			OutputPathList: make([]string, len(cfg.PlotList[0].OutputPathList)),
		}
		allMemoryFrameCfg.OutputPathList[0] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-VMRSS-MB-BY-KEY.svg")
		allMemoryFrameCfg.OutputPathList[1] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-VMRSS-MB-BY-KEY.png")
		plog.Printf("plotting %v", allMemoryFrameCfg.OutputPathList)
		var pairs []pair
		allCols := allMemoryFrame.Columns()
		for i := 0; i < len(allCols)-3; i += 4 {
			pairs = append(pairs, pair{
				x: allCols[i],   // x
				y: allCols[i+2], // avg
			})
		}
		if err = all.drawXY(allMemoryFrameCfg, pairs...); err != nil {
			return err
		}
		newCSV := dataframe.New()
		for _, p := range pairs {
			if err = newCSV.AddColumn(p.x); err != nil {
				return err
			}
			if err = newCSV.AddColumn(p.y); err != nil {
				return err
			}
		}
		csvPath := filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-VMRSS-MB-BY-KEY.csv")
		if err := newCSV.CSV(csvPath); err != nil {
			return err
		}
	}
	{
		// with error points
		allMemoryFrameCfg := PlotConfig{
			Column:         "AVG-VMRSS-MB",
			XAxis:          "Cumulative Number of Keys",
			YAxis:          "Memory(MB) by Keys",
			OutputPathList: make([]string, len(cfg.PlotList[0].OutputPathList)),
		}
		allMemoryFrameCfg.OutputPathList[0] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg")
		allMemoryFrameCfg.OutputPathList[1] = filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.png")
		plog.Printf("plotting %v", allMemoryFrameCfg.OutputPathList)
		var triplets []triplet
		allCols := allMemoryFrame.Columns()
		for i := 0; i < len(allCols)-3; i += 4 {
			triplets = append(triplets, triplet{
				x:      allCols[i],
				minCol: allCols[i+1],
				avgCol: allCols[i+2],
				maxCol: allCols[i+3],
			})
		}
		if err = all.drawXYWithErrorPoints(allMemoryFrameCfg, triplets...); err != nil {
			return err
		}
		newCSV := dataframe.New()
		for _, tri := range triplets {
			if err = newCSV.AddColumn(tri.x); err != nil {
				return err
			}
			if err = newCSV.AddColumn(tri.minCol); err != nil {
				return err
			}
			if err = newCSV.AddColumn(tri.avgCol); err != nil {
				return err
			}
			if err = newCSV.AddColumn(tri.maxCol); err != nil {
				return err
			}
		}
		csvPath := filepath.Join(filepath.Dir(cfg.PlotList[0].OutputPathList[0]), "AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.csv")
		if err := newCSV.CSV(csvPath); err != nil {
			return err
		}
	}

	plog.Println("combining data for plotting")
	for _, plotConfig := range cfg.PlotList {
		plog.Printf("plotting %q", plotConfig.Column)
		var clientNumColumns []dataframe.Column
		var pairs []pair
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
			pairs = append(pairs, pair{y: col})
			dataColumns = append(dataColumns, col)
		}
		if err = all.draw(plotConfig, pairs...); err != nil {
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
