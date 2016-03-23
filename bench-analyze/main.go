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

// bench-analyze analyzes bench results.
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gyuho/psn/ps"
)

var (
	dbtesterBenchColumns = map[string]int{
		"unix_ts":        0,
		"avg_latency_ms": 1,
		"throughput":     2,
	}
)

func main() {
	var (
		prefixes = []string{
			"testdata/bench-01-consul-",
			"testdata/bench-01-etcd-",
			"testdata/bench-01-etcd2-",
			"testdata/bench-01-zk-",
		}
		comparedPath = "testdata/bench-01-compared.csv"
	)

	tbs := []ps.Table{}
	tableToSuffix := make(map[int]string)

	tableToLatencyIdx := make(map[int]int)
	tableToThroughputIdx := make(map[int]int)
	tableToCpuIdx := make(map[int]int)
	tableToMemoryIdx := make(map[int]int)
	maxSize := 0

	var compareColumns = map[string]int{
		"second": 0,
	}
	initSize := len(compareColumns)
	for i, prefix := range prefixes {
		tb, err := aggregate(prefix)
		if err != nil {
			log.Fatal(err)
		}
		tbs = append(tbs, tb)
		if strings.Contains(prefix, "-etcd-") {
			tableToSuffix[i] = "etcd"
		} else if strings.Contains(prefix, "-zk-") {
			tableToSuffix[i] = "zk"
		} else if strings.Contains(prefix, "-etcd2-") {
			tableToSuffix[i] = "etcd2"
		} else if strings.Contains(prefix, "-consul-") {
			tableToSuffix[i] = "consul"
		}

		tableToLatencyIdx[i] = tb.Columns["avg_latency_ms"]
		tableToThroughputIdx[i] = tb.Columns["throughput"]
		tableToCpuIdx[i] = tb.Columns["avg_cpu"]
		tableToMemoryIdx[i] = tb.Columns["avg_memory_mb"]

		if maxSize < len(tb.Rows) {
			maxSize = len(tb.Rows)
		}

		compareColumns["avg_latency_ms_"+tableToSuffix[i]] = 4*i + initSize
		compareColumns["throughput_"+tableToSuffix[i]] = 4*i + initSize + 1
		compareColumns["avg_cpu_"+tableToSuffix[i]] = 4*i + initSize + 2
		compareColumns["avg_memory_mb_"+tableToSuffix[i]] = 4*i + initSize + 3
	}

	columnSlice := make([]string, len(compareColumns))
	for k, v := range compareColumns {
		columnSlice[v] = k
	}

	cTable := ps.Table{}
	cTable.Columns = compareColumns
	cTable.ColumnSlice = columnSlice
	crows := make([][]string, maxSize)
	for i, tb := range tbs {
		latencyIdx := tableToLatencyIdx[i]
		throughputIdx := tableToThroughputIdx[i]
		cpuIdx := tableToCpuIdx[i]
		memoryIdx := tableToMemoryIdx[i]
		for j, row := range tb.Rows {
			if len(crows[j]) == 0 {
				crows[j] = []string{fmt.Sprintf("%d", j)}
			}
			crows[j] = append(crows[j], row[latencyIdx], row[throughputIdx], row[cpuIdx], row[memoryIdx])
		}
		for k := len(tb.Rows); k < maxSize; k++ {
			crows[k] = append(crows[k], "", "", "", "")
		}
	}
	cTable.Rows = crows

	if err := cTable.ToCSV(comparedPath); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully saved compared.csv")

	// TODO:
	// add plotting
}

func aggregate(prefix string) (ps.Table, error) {
	var (
		dbtesterBenchColumns = map[string]int{
			"unix_ts":        0,
			"avg_latency_ms": 1,
			"throughput":     2,
		}
		benchmarkResultPath = fmt.Sprintf("%stimeseries.csv", prefix)
		testPaths           = []string{
			fmt.Sprintf("%s1-monitor.csv", prefix),
			fmt.Sprintf("%s2-monitor.csv", prefix),
			fmt.Sprintf("%s3-monitor.csv", prefix),
		}
		finalPath = fmt.Sprintf("%sfinal.csv", prefix)
	)
	log.Printf("Aggregating %q\n", testPaths)

	tbResultCombined, err := ps.ReadCSVs(ps.ColumnsPS, testPaths...)
	if err != nil {
		return ps.Table{}, err
	}

	tbResultBench, err := ps.ReadCSV(dbtesterBenchColumns, benchmarkResultPath)
	if err != nil {
		return ps.Table{}, err
	}

	tIdx := 0
	for i := range tbResultCombined.Rows {
		ts, err := strconv.ParseInt(tbResultCombined.Rows[i][0], 10, 64)
		if err != nil {
			return ps.Table{}, err
		}
		if ts == tbResultBench.MinTS {
			tbResultCombined.MinTS = tbResultBench.MinTS
			tIdx = i
		}
	}
	tbResultCombined.Rows = tbResultCombined.Rows[tIdx:]

	// now combine tbResultBench with tbResultCombined
	tbFinal := ps.Table{}
	tbFinal.MinTS = tbResultBench.MinTS
	tbFinal.MaxTS = tbResultBench.MaxTS
	tbFinal.Columns = tbResultBench.Columns
	cSize := len(tbResultBench.Columns)
	for k, v := range tbResultCombined.Columns {
		if v == 0 {
			continue // skip unix_ts
		}
		tbFinal.Columns[k] = v + cSize - 1
	}
	fSize := len(tbFinal.Columns)
	tbFinal.Columns["avg_cpu"] = fSize
	tbFinal.Columns["avg_memory_mb"] = fSize + 1

	columnSlice := make([]string, len(tbFinal.Columns))
	for k, v := range tbFinal.Columns {
		columnSlice[v] = k
	}
	tbFinal.ColumnSlice = columnSlice

	cpuIdxs, memoryIdxs := []int{}, []int{}
	for i := range tbFinal.ColumnSlice {
		if strings.HasPrefix(tbFinal.ColumnSlice[i], "cpu_") {
			cpuIdxs = append(cpuIdxs, i-cSize+1)
		}
		if strings.HasPrefix(tbFinal.ColumnSlice[i], "memory_") {
			memoryIdxs = append(memoryIdxs, i-cSize+1)
		}
	}

	nrows := make([][]string, len(tbResultBench.Rows))
	for i, row := range tbResultBench.Rows {
		resultRow := tbResultCombined.Rows[i][1:]

		var totalCpu float64
		for _, idx := range cpuIdxs {
			f, err := strconv.ParseFloat(tbResultCombined.Rows[i][idx], 64)
			if err != nil {
				return ps.Table{}, err
			}
			totalCpu += f
		}
		avgCpu := totalCpu / float64(len(testPaths))
		var totalMemory float64
		for _, idx := range memoryIdxs {
			f, err := strconv.ParseFloat(tbResultCombined.Rows[i][idx], 64)
			if err != nil {
				return ps.Table{}, err
			}
			totalMemory += f
		}
		avgMemory := totalMemory / float64(len(testPaths))
		resultRow = append(resultRow, fmt.Sprintf("%.2f", avgCpu), fmt.Sprintf("%.2f", avgMemory))

		nrows[i] = append(row, resultRow...)
	}
	tbFinal.Rows = nrows

	if err := tbFinal.ToCSV(finalPath); err != nil {
		return ps.Table{}, err
	}

	log.Printf("Successfully saved %s\n", finalPath)
	return tbFinal, nil
}
