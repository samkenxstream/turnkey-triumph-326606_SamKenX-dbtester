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
	"strconv"
	"strings"

	"github.com/coreos/dbtester"
	humanize "github.com/dustin/go-humanize"
	"github.com/gyuho/dataframe"
	"github.com/olekukonko/tablewriter"
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

func commandFunc(cmd *cobra.Command, args []string) error {
	return do(configPath)
}

type allAggregatedData struct {
	title                       string
	data                        []*analyzeData
	headerToDatabaseID          map[string]string
	headerToDatabaseDescription map[string]string
	allDatabaseIDList           []string
}

func do(configPath string) error {
	cfg, err := dbtester.ReadConfig(configPath, true)
	if err != nil {
		return err
	}

	all := &allAggregatedData{
		title:                       cfg.TestTitle,
		data:                        make([]*analyzeData, 0, len(cfg.DatabaseIDToTestData)),
		headerToDatabaseID:          make(map[string]string),
		headerToDatabaseDescription: make(map[string]string),
		allDatabaseIDList:           cfg.AllDatabaseIDList,
	}
	for _, databaseID := range cfg.AllDatabaseIDList {
		testgroup := cfg.DatabaseIDToTestGroup[databaseID]
		testdata := cfg.DatabaseIDToTestData[databaseID]

		plog.Printf("reading system metrics data for %s", databaseID)
		ad, err := readSystemMetricsAll(testdata.ServerSystemMetricsInterpolatedPathList...)
		if err != nil {
			return err
		}
		ad.databaseTag = testgroup.DatabaseTag
		ad.legend = testgroup.DatabaseDescription
		ad.allAggregatedOutputPath = testdata.AllAggregatedOutputPath

		if err = ad.aggSystemMetrics(); err != nil {
			return err
		}
		if err = ad.importBenchMetrics(testdata.ClientLatencyThroughputTimeseriesPath); err != nil {
			return err
		}
		if err = ad.aggregateAll(testdata.ServerMemoryByKeyNumberPath, testgroup.RequestNumber); err != nil {
			return err
		}
		if err = ad.save(); err != nil {
			return err
		}

		all.data = append(all.data, ad)
		for _, hd := range ad.aggregated.Headers() {
			all.headerToDatabaseID[makeHeader(hd, testgroup.DatabaseTag)] = databaseID
			all.headerToDatabaseDescription[makeHeader(hd, testgroup.DatabaseTag)] = testgroup.DatabaseDescription
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
			databaseID := all.headerToDatabaseID[col.Header()]
			row00Header = append(row00Header, cfg.DatabaseIDToTestGroup[databaseID].DatabaseTag)
			break
		}
	}

	row17ServerReceiveBytesSum := []string{"SERVER-TOTAL-NETWORK-RX-DATA-SUM"}
	row17ServerReceiveBytesSumRaw := []string{"SERVER-TOTAL-NETWORK-RX-DATA-BYTES-SUM-RAW"}
	row18ServerTransmitBytesSum := []string{"SERVER-TOTAL-NETWORK-TX-DATA-SUM"}
	row18ServerTransmitBytesSumRaw := []string{"SERVER-TOTAL-NETWORK-TX-DATA-BYTES-SUM-RAW"}
	row21ServerMaxCPUUsage := []string{"SERVER-MAX-CPU-USAGE"}
	row22ServerMaxMemoryUsage := []string{"SERVER-MAX-MEMORY-USAGE"}
	row26ReadsCompletedDeltaSum := []string{"SERVER-AVG-READS-COMPLETED-DELTA-SUM"}
	row27SectorsReadDeltaSum := []string{"SERVER-AVG-SECTORS-READS-DELTA-SUM"}
	row28WritesCompletedDeltaSum := []string{"SERVER-AVG-WRITES-COMPLETED-DELTA-SUM"}
	row29SectorsWrittenDeltaSum := []string{"SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM"}

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

		row17ServerReceiveBytesSum = append(row17ServerReceiveBytesSum, humanize.Bytes(uint64(receiveBytesNumDeltaSum)))
		row17ServerReceiveBytesSumRaw = append(row17ServerReceiveBytesSumRaw, fmt.Sprintf("%.2f", receiveBytesNumDeltaSum))
		row18ServerTransmitBytesSum = append(row18ServerTransmitBytesSum, humanize.Bytes(uint64(transmitBytesNumDeltaSum)))
		row18ServerTransmitBytesSumRaw = append(row18ServerTransmitBytesSumRaw, fmt.Sprintf("%.2f", transmitBytesNumDeltaSum))
		row21ServerMaxCPUUsage = append(row21ServerMaxCPUUsage, fmt.Sprintf("%.2f %%", maxAvgCPU))
		row26ReadsCompletedDeltaSum = append(row26ReadsCompletedDeltaSum, humanize.Comma(int64(readsCompletedDeltaSum)))
		row27SectorsReadDeltaSum = append(row27SectorsReadDeltaSum, humanize.Comma(int64(sectorsReadDeltaSum)))
		row28WritesCompletedDeltaSum = append(row28WritesCompletedDeltaSum, humanize.Comma(int64(writesCompletedDeltaSum)))
		row29SectorsWrittenDeltaSum = append(row29SectorsWrittenDeltaSum, humanize.Comma(int64(sectorsWrittenDeltaSum)))

		// TODO: handle overflowed memory value?
		sort.Float64s(maxAvgVMRSSMBs)
		mv := maxAvgVMRSSMBs[len(maxAvgVMRSSMBs)-1]
		mb := uint64(mv * 1000000)
		row22ServerMaxMemoryUsage = append(row22ServerMaxMemoryUsage, humanize.Bytes(mb))
	}

	row01TotalSeconds := []string{"TOTAL-SECONDS"} // TOTAL-SECONDS
	row02TotalRequestNumber := []string{"TOTAL-REQUEST-NUMBER"}
	row03MaxThroughput := []string{"MAX-THROUGHPUT"}                                    // MAX AVG-THROUGHPUT
	row04AverageThroughput := []string{"AVG-THROUGHPUT"}                                // REQUESTS-PER-SECOND
	row05MinThroughput := []string{"MIN-THROUGHPUT"}                                    // MIN AVG-THROUGHPUT
	row06FastestLatency := []string{"FASTEST-LATENCY"}                                  // FASTEST-LATENCY-MS
	row07AverageLatency := []string{"AVG-LATENCY"}                                      // AVERAGE-LATENCY-MS
	row08SlowestLatency := []string{"SLOWEST-LATENCY"}                                  // SLOWEST-LATENCY-MS
	row09p10 := []string{"Latency p10"}                                                 // p10
	row10p25 := []string{"Latency p25"}                                                 // p25
	row11p50 := []string{"Latency p50"}                                                 // p50
	row12p75 := []string{"Latency p75"}                                                 // p75
	row13p90 := []string{"Latency p90"}                                                 // p90
	row14p95 := []string{"Latency p95"}                                                 // p95
	row15p99 := []string{"Latency p99"}                                                 // p99
	row16p999 := []string{"Latency p99.9"}                                              // p99.9
	row19ClientReceiveBytesSum := []string{"CLIENT-TOTAL-NETWORK-RX-SUM"}               // RECEIVE-BYTES-NUM-DELTA
	row19ClientReceiveBytesSumRaw := []string{"CLIENT-TOTAL-NETWORK-RX-BYTES-SUM-RAW"}  // RECEIVE-BYTES-NUM-DELTA
	row20ClientTransmitBytesSum := []string{"CLIENT-TOTAL-NETWORK-TX-SUM"}              // TRANSMIT-BYTES-DELTA
	row20ClientTransmitBytesSumRaw := []string{"CLIENT-TOTAL-NETWORK-TX-BYTES-SUM-RAW"} // TRANSMIT-BYTES-DELTA
	row23ClientMaxCPU := []string{"CLIENT-MAX-CPU-USAGE"}                               // CPU-NUM
	row24ClientMaxMemory := []string{"CLIENT-MAX-MEMORY-USAGE"}                         // VMRSS-NUM
	row25ClientErrorCount := []string{"CLIENT-ERROR-COUNT"}                             // ERROR:
	row30AverageDatasize := []string{"SERVER-AVG-DATA-SIZE-ON-DISK"}                    // TOTAL-DATA-SIZE

	databaseIDToErrs := make(map[string][]string)
	for i, databaseID := range cfg.AllDatabaseIDList {
		testgroup := cfg.DatabaseIDToTestGroup[databaseID]
		testdata := cfg.DatabaseIDToTestData[databaseID]

		tag := testdata.DatabaseTag
		if tag != row00Header[i+1] {
			return fmt.Errorf("analyze config has different order; expected %q, got %q", row00Header[i+1], tag)
		}
		row02TotalRequestNumber = append(row02TotalRequestNumber, humanize.Comma(testgroup.RequestNumber))

		{
			fr, err := dataframe.NewFromCSV(nil, testdata.ClientSystemMetricsInterpolatedPath)
			if err != nil {
				return err
			}

			var receiveBytesNumDeltaSum float64
			col, err := fr.Column("RECEIVE-BYTES-NUM-DELTA")
			if err != nil {
				return err
			}
			for i := 0; i < col.Count(); i++ {
				v, err := col.Value(i)
				if err != nil {
					return err
				}
				fv, _ := v.Float64()
				receiveBytesNumDeltaSum += fv
			}

			var transmitBytesNumDeltaSum float64
			col, err = fr.Column("TRANSMIT-BYTES-NUM-DELTA")
			if err != nil {
				return err
			}
			for i := 0; i < col.Count(); i++ {
				v, err := col.Value(i)
				if err != nil {
					return err
				}
				fv, _ := v.Float64()
				transmitBytesNumDeltaSum += fv
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

			row19ClientReceiveBytesSum = append(row19ClientReceiveBytesSum, humanize.Bytes(uint64(receiveBytesNumDeltaSum)))
			row19ClientReceiveBytesSumRaw = append(row19ClientReceiveBytesSumRaw, fmt.Sprintf("%.2f", receiveBytesNumDeltaSum))
			row20ClientTransmitBytesSum = append(row20ClientTransmitBytesSum, humanize.Bytes(uint64(transmitBytesNumDeltaSum)))
			row20ClientTransmitBytesSumRaw = append(row20ClientTransmitBytesSumRaw, fmt.Sprintf("%.2f", transmitBytesNumDeltaSum))
			row23ClientMaxCPU = append(row23ClientMaxCPU, fmt.Sprintf("%.2f %%", maxAvgCPU))
			row24ClientMaxMemory = append(row24ClientMaxMemory, humanize.Bytes(maxVMRSSNum))
		}
		{
			f, err := openToRead(testdata.ClientLatencyDistributionSummaryPath)
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
					row01TotalSeconds = append(row01TotalSeconds, fmt.Sprintf("%s sec", row[1]))
				case "REQUESTS-PER-SECOND":
					fv, err := strconv.ParseFloat(row[1], 64)
					if err != nil {
						return err
					}
					avg := int64(fv)
					row04AverageThroughput = append(row04AverageThroughput, fmt.Sprintf("%s req/sec", humanize.Comma(avg)))
				case "SLOWEST-LATENCY-MS":
					row08SlowestLatency = append(row08SlowestLatency, fmt.Sprintf("%s ms", row[1]))
				case "FASTEST-LATENCY-MS":
					row06FastestLatency = append(row06FastestLatency, fmt.Sprintf("%s ms", row[1]))
				case "AVERAGE-LATENCY-MS":
					row07AverageLatency = append(row07AverageLatency, fmt.Sprintf("%s ms", row[1]))
				}

				if strings.HasPrefix(row[0], "ERROR:") {
					iv, err := strconv.ParseInt(row[1], 10, 64)
					if err != nil {
						return err
					}
					totalErrCnt += iv

					c1 := strings.TrimSpace(strings.Replace(row[0], "ERROR:", "", -1))
					c2 := humanize.Comma(iv)
					es := fmt.Sprintf("%s (count %s)", c1, c2)
					if _, ok := databaseIDToErrs[databaseID]; !ok {
						databaseIDToErrs[databaseID] = []string{es}
					} else {
						databaseIDToErrs[databaseID] = append(databaseIDToErrs[databaseID], es)
					}
				}
			}
			row25ClientErrorCount = append(row25ClientErrorCount, humanize.Comma(totalErrCnt))
		}
		{
			fr, err := dataframe.NewFromCSV(nil, testdata.ClientLatencyThroughputTimeseriesPath)
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
			row03MaxThroughput = append(row03MaxThroughput, fmt.Sprintf("%s req/sec", humanize.Comma(max)))
			row05MinThroughput = append(row05MinThroughput, fmt.Sprintf("%s req/sec", humanize.Comma(min)))
		}
		{
			fr, err := dataframe.NewFromCSV(nil, testdata.ServerDatasizeOnDiskSummaryPath)
			if err != nil {
				return err
			}
			col, err := fr.Column(dbtester.DatasizeOnDiskSummaryColumns[3]) // datasize in bytes
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
			row30AverageDatasize = append(row30AverageDatasize, humanize.Bytes(avg))
		}
		{
			f, err := openToRead(testdata.ClientLatencyDistributionPercentilePath)
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
					row09p10 = append(row09p10, fmt.Sprintf("%s ms", row[1]))
				case "p25":
					row10p25 = append(row10p25, fmt.Sprintf("%s ms", row[1]))
				case "p50":
					row11p50 = append(row11p50, fmt.Sprintf("%s ms", row[1]))
				case "p75":
					row12p75 = append(row12p75, fmt.Sprintf("%s ms", row[1]))
				case "p90":
					row13p90 = append(row13p90, fmt.Sprintf("%s ms", row[1]))
				case "p95":
					row14p95 = append(row14p95, fmt.Sprintf("%s ms", row[1]))
				case "p99":
					row15p99 = append(row15p99, fmt.Sprintf("%s ms", row[1]))
				case "p99.9":
					row16p999 = append(row16p999, fmt.Sprintf("%s ms", row[1]))
				}
			}
		}
	}

	plog.Printf("saving summary data to %q", cfg.Analyze.AllAggregatedOutputPathCSV)
	aggRowsForSummaryCSV := [][]string{
		row00Header,
		row01TotalSeconds,
		row02TotalRequestNumber,
		row03MaxThroughput,
		row04AverageThroughput,
		row05MinThroughput,

		row06FastestLatency,
		row07AverageLatency,
		row08SlowestLatency,

		row09p10,
		row10p25,
		row11p50,
		row12p75,
		row13p90,
		row14p95,
		row15p99,
		row16p999,

		row17ServerReceiveBytesSum,
		row17ServerReceiveBytesSumRaw,
		row18ServerTransmitBytesSum,
		row18ServerTransmitBytesSumRaw,
		row19ClientReceiveBytesSum,
		row19ClientReceiveBytesSumRaw,
		row20ClientTransmitBytesSum,
		row20ClientTransmitBytesSumRaw,

		row21ServerMaxCPUUsage,
		row22ServerMaxMemoryUsage,
		row23ClientMaxCPU,
		row24ClientMaxMemory,

		row25ClientErrorCount,

		row26ReadsCompletedDeltaSum,
		row27SectorsReadDeltaSum,
		row28WritesCompletedDeltaSum,
		row29SectorsWrittenDeltaSum,
		row30AverageDatasize,
	}
	file, err := openToOverwrite(cfg.Analyze.AllAggregatedOutputPathCSV)
	if err != nil {
		return err
	}
	defer file.Close()
	wr := csv.NewWriter(file)
	if err := wr.WriteAll(aggRowsForSummaryCSV); err != nil {
		return err
	}
	wr.Flush()
	if err := wr.Error(); err != nil {
		return err
	}

	plog.Printf("saving summary data to %q", cfg.Analyze.AllAggregatedOutputPathTXT)
	aggRowsForSummaryTXT := [][]string{
		row00Header,
		row01TotalSeconds,
		row02TotalRequestNumber,
		row03MaxThroughput,
		row04AverageThroughput,
		row05MinThroughput,

		row06FastestLatency,
		row07AverageLatency,
		row08SlowestLatency,

		row09p10,
		row10p25,
		row11p50,
		row12p75,
		row13p90,
		row14p95,
		row15p99,
		row16p999,

		row17ServerReceiveBytesSum,
		row18ServerTransmitBytesSum,
		row19ClientReceiveBytesSum,
		row20ClientTransmitBytesSum,

		row21ServerMaxCPUUsage,
		row22ServerMaxMemoryUsage,
		row23ClientMaxCPU,
		row24ClientMaxMemory,

		row25ClientErrorCount,

		row26ReadsCompletedDeltaSum,
		row27SectorsReadDeltaSum,
		row28WritesCompletedDeltaSum,
		row29SectorsWrittenDeltaSum,
		row30AverageDatasize,
	}
	buf := new(bytes.Buffer)
	tw := tablewriter.NewWriter(buf)
	tw.SetHeader(aggRowsForSummaryTXT[0])
	for _, row := range aggRowsForSummaryTXT[1:] {
		tw.Append(row)
	}
	tw.SetAutoFormatHeaders(false)
	tw.SetAlignment(tablewriter.ALIGN_RIGHT)
	tw.Render()
	errs := ""
	for _, databaseID := range cfg.AllDatabaseIDList {
		es, ok := databaseIDToErrs[databaseID]
		if !ok {
			continue
		}
		errs = databaseID + " " + "errors:\n" + strings.Join(es, "\n") + "\n"
	}
	stxt := buf.String()
	if errs != "" {
		stxt += "\n" + "\n" + errs
	}
	if err := toFile(stxt, changeExtToTxt(cfg.Analyze.AllAggregatedOutputPathTXT)); err != nil {
		return err
	}

	// KEYS, MIN-LATENCY-MS, AVG-LATENCY-MS, MAX-LATENCY-MS
	plog.Info("combining all latency data by keys")
	allLatencyFrame := dataframe.New()
	for _, databaseID := range cfg.AllDatabaseIDList {
		testdata := cfg.DatabaseIDToTestData[databaseID]

		fr, err := dataframe.NewFromCSV(nil, testdata.ClientLatencyByKeyNumberPath)
		if err != nil {
			return err
		}
		colKeys, err := fr.Column("KEYS")
		if err != nil {
			return err
		}
		colKeys.UpdateHeader(makeHeader("KEYS", testdata.DatabaseTag))
		if err = allLatencyFrame.AddColumn(colKeys); err != nil {
			return err
		}

		colMinLatency, err := fr.Column("MIN-LATENCY-MS")
		if err != nil {
			return err
		}
		colMinLatency.UpdateHeader(makeHeader("MIN-LATENCY-MS", testdata.DatabaseTag))
		if err = allLatencyFrame.AddColumn(colMinLatency); err != nil {
			return err
		}

		colAvgLatency, err := fr.Column("AVG-LATENCY-MS")
		if err != nil {
			return err
		}
		colAvgLatency.UpdateHeader(makeHeader("AVG-LATENCY-MS", testdata.DatabaseTag))
		if err = allLatencyFrame.AddColumn(colAvgLatency); err != nil {
			return err
		}

		colMaxLatency, err := fr.Column("MAX-LATENCY-MS")
		if err != nil {
			return err
		}
		colMaxLatency.UpdateHeader(makeHeader("MAX-LATENCY-MS", testdata.DatabaseTag))
		if err = allLatencyFrame.AddColumn(colMaxLatency); err != nil {
			return err
		}
	}
	// KEYS, MIN-VMRSS-MB, AVG-VMRSS-MB, MAX-VMRSS-MB
	plog.Info("combining all server memory usage by keys")
	allMemoryFrame := dataframe.New()
	for _, databaseID := range cfg.AllDatabaseIDList {
		testdata := cfg.DatabaseIDToTestData[databaseID]

		fr, err := dataframe.NewFromCSV(nil, testdata.ServerMemoryByKeyNumberPath)
		if err != nil {
			return err
		}
		colKeys, err := fr.Column("KEYS")
		if err != nil {
			return err
		}
		colKeys.UpdateHeader(makeHeader("KEYS", testdata.DatabaseTag))
		if err = allMemoryFrame.AddColumn(colKeys); err != nil {
			return err
		}

		colMemMin, err := fr.Column("MIN-VMRSS-MB")
		if err != nil {
			return err
		}
		colMemMin.UpdateHeader(makeHeader("MIN-VMRSS-MB", testdata.DatabaseTag))
		if err = allMemoryFrame.AddColumn(colMemMin); err != nil {
			return err
		}

		colMem, err := fr.Column("AVG-VMRSS-MB")
		if err != nil {
			return err
		}
		colMem.UpdateHeader(makeHeader("AVG-VMRSS-MB", testdata.DatabaseTag))
		if err = allMemoryFrame.AddColumn(colMem); err != nil {
			return err
		}

		colMemMax, err := fr.Column("MAX-VMRSS-MB")
		if err != nil {
			return err
		}
		colMemMax.UpdateHeader(makeHeader("MAX-VMRSS-MB", testdata.DatabaseTag))
		if err = allMemoryFrame.AddColumn(colMemMax); err != nil {
			return err
		}
	}

	{
		allLatencyFrameCfg := dbtester.Plot{
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
		allLatencyFrameCfg := dbtester.Plot{
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
	{
		allMemoryFrameCfg := dbtester.Plot{
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
		allMemoryFrameCfg := dbtester.Plot{
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
			databaseID := all.allDatabaseIDList[i]
			tag := cfg.DatabaseIDToTestGroup[databaseID].DatabaseTag

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
		if err = nf1.CSV(plotConfig.OutputPathCSV); err != nil {
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
		if err = nf2.CSV(filepath.Join(filepath.Dir(plotConfig.OutputPathCSV), plotConfig.Column+"-BY-CLIENT-NUM"+".csv")); err != nil {
			return err
		}

		if len(cfg.DatabaseIDToTestGroup[cfg.AllDatabaseIDList[0]].BenchmarkOptions.ConnectionClientNumbers) > 0 {
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
			if err = nf3.CSV(filepath.Join(filepath.Dir(plotConfig.OutputPathCSV), plotConfig.Column+"-BY-CLIENT-NUM-aggregated"+".csv")); err != nil {
				return err
			}
		}
	}

	return cfg.WriteREADME(stxt)
}

func changeExtToTxt(fpath string) string {
	ext := filepath.Ext(fpath)
	return strings.Replace(fpath, ext, ".txt", -1)
}
