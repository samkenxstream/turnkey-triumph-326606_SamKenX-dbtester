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
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/gyuho/dataframe"
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
		plog.Printf("reading system metrics data for %s (%q)", elem.DatabaseTag, elem.Legend)
		ad, err := readSystemMetricsAll(elem.DataSystemMetricsPaths...)
		if err != nil {
			return err
		}
		ad.databaseTag = elem.DatabaseTag
		ad.legend = elem.Legend
		ad.csvOutputpath = elem.OutputPath

		if err = ad.aggSystemMetrics(); err != nil {
			return err
		}
		if err = ad.importBenchMetrics(elem.DataBenchmarkThroughput); err != nil {
			return err
		}
		if err = ad.aggregateAll(); err != nil {
			return err
		}
		if err = ad.save(); err != nil {
			return err
		}

		all.data = append(all.data, ad)
		all.databaseTags = append(all.databaseTags, elem.DatabaseTag)
		for _, hd := range ad.aggregated.Headers() {
			all.headerToLegend[makeHeader(hd, elem.DatabaseTag)] = elem.Legend
		}
	}

	// aggregated everything
	// 1. sum of all network usage per database
	// 2. throughput, latency percentiles distribution
	//
	// FIRST ROW FOR HEADER: etcd, Zookeeper, Consul, ...
	// FIRST COLUMN FOR LABELS: READS-COMPLETED-DELTA, ...
	// SECOND COLUMNS ~ FOR DATA
	row1Header := []string{""} // first is empty
	for _, ad := range all.data {
		// per database
		for _, col := range ad.aggregated.Columns() {
			legend := all.headerToLegend[col.Header()]
			row1Header = append(row1Header, makeTag(legend))
			break
		}
	}
	row2ReadsCompletedDeltaSum := []string{"READS-COMPLETED-DELTA"}
	row3SectorsReadDeltaSum := []string{"SECTORS-READS-DELTA-SUM"}
	row4WritesCompletedDeltaSum := []string{"WRITES-COMPLETED-DELTA-SUM"}
	row5SectorsWrittenDeltaSum := []string{"SECTORS-WRITTEN-DELTA-SUM"}
	row6ReceiveBytesSum := []string{"RECEIVE-BYTES-SUM"}
	row7ReceiveBytesNumDeltaSum := []string{"RECEIVE-BYTES-NUM-DELTA-SUM"}
	row8TransmitBytesSum := []string{"TRANSMIT-BYTES-SUM"}
	row9TransmitBytesNumDeltaSum := []string{"TRANSMIT-BYTES-NUM-DELTA-SUM"}
	row10AverageThroughput := []string{"AVG-THROUGHPUT"}
	row11AverageLatency := []string{"AVG-LATENCY"}

	// iterate each database's all data
	for _, ad := range all.data {
		// ad.benchMetrics.frame.Co
		var (
			readsCompletedDeltaSum   float64
			sectorsReadDeltaSum      float64
			writesCompletedDeltaSum  float64
			sectorsWrittenDeltaSum   float64
			receiveBytesNumDeltaSum  float64
			transmitBytesNumDeltaSum float64
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
					fv, _ := vv.Number()
					readsCompletedDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "SECTORS-READS-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Number()
					sectorsReadDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "WRITES-COMPLETED-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Number()
					writesCompletedDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "SECTORS-WRITTEN-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Number()
					sectorsWrittenDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "RECEIVE-BYTES-NUM-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Number()
					receiveBytesNumDeltaSum += fv
				}
			case strings.HasPrefix(hdr, "TRANSMIT-BYTES-NUM-DELTA-"):
				cnt := col.Count()
				for j := 0; j < cnt; j++ {
					vv, err := col.Value(j)
					if err != nil {
						return err
					}
					fv, _ := vv.Number()
					transmitBytesNumDeltaSum += fv
				}
			}
		}

		row2ReadsCompletedDeltaSum = append(row2ReadsCompletedDeltaSum, fmt.Sprintf("%d", uint64(readsCompletedDeltaSum)))
		row3SectorsReadDeltaSum = append(row3SectorsReadDeltaSum, fmt.Sprintf("%d", uint64(sectorsReadDeltaSum)))
		row4WritesCompletedDeltaSum = append(row4WritesCompletedDeltaSum, fmt.Sprintf("%d", uint64(writesCompletedDeltaSum)))
		row5SectorsWrittenDeltaSum = append(row5SectorsWrittenDeltaSum, fmt.Sprintf("%d", uint64(sectorsWrittenDeltaSum)))
		row6ReceiveBytesSum = append(row6ReceiveBytesSum, humanize.Bytes(uint64(receiveBytesNumDeltaSum)))
		row7ReceiveBytesNumDeltaSum = append(row7ReceiveBytesNumDeltaSum, fmt.Sprintf("%d", uint64(receiveBytesNumDeltaSum)))
		row8TransmitBytesSum = append(row8TransmitBytesSum, humanize.Bytes(uint64(transmitBytesNumDeltaSum)))
		row9TransmitBytesNumDeltaSum = append(row9TransmitBytesNumDeltaSum, fmt.Sprintf("%d", uint64(transmitBytesNumDeltaSum)))
	}

	aggDf := dataframe.New()

	plog.Printf("saving data to %q", cfg.AllAggregatedPath)
	if err := aggDf.CSVHorizontal(cfg.AllAggregatedPath); err != nil {
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
				return fmt.Errorf("clientNumColumns[i].Count() %d != dataColumns[i].Count() %d", clientNumColumns[i].Count(), dataColumns[i].Count())
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
				num, _ := v1.Number()

				v2, err := dataColumns[i].Value(j)
				if err != nil {
					return err
				}
				data, _ := v2.Number()

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
				// col1 := dataframe.NewColumn(clientNumColumns[i].Header())
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

	plog.Printf("writing README at %q", cfg.READMEConfig.OutputPath)
	return writeREADME(cfg.READMEConfig)
}
