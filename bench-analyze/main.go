package main

import (
	"log"
	"strconv"

	"github.com/gyuho/psn/ps"
)

var (
	dbtesterBenchColumns = map[string]int{
		"unix_ts":        0,
		"avg_latency_ms": 1,
		"throughput":     2,
	}
	benchmarkResultPath = "testdata/test-01-etcd-timeseries.csv"
	testPaths           = []string{
		"testdata/test-01-etcd-server-1.csv",
		"testdata/test-01-etcd-server-2.csv",
		"testdata/test-01-etcd-server-3.csv",
	}
	combinedPath = "testdata/test-01-etcd-combined.csv"
)

func main() {
	tTest, err := ps.ReadCSVs(ps.ColumnsPS, testPaths...)
	if err != nil {
		log.Fatal(err)
	}

	tBench, err := ps.ReadCSV(dbtesterBenchColumns, benchmarkResultPath)
	if err != nil {
		log.Fatal(err)
	}

	tIdx := 0
	for i := range tTest.Rows {
		ts, err := strconv.ParseInt(tTest.Rows[i][0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		if ts == tBench.MinTS {
			log.Println("Truncating tTest from tBench's minimum ts", tBench.MinTS, "at row index", i)
			tTest.MinTS = tBench.MinTS
			tIdx = i
		}
	}
	tTest.Rows = tTest.Rows[tIdx:]

	// now combine tBench with tTest
	tCombined := ps.Table{}

	tCombined.Columns = tBench.Columns
	cSize := len(tCombined.Columns)
	for k, v := range tTest.Columns {
		if v == 0 {
			continue // skip unix_ts
		}
		tCombined.Columns[k] = v + cSize - 1
	}

	columnSlice := make([]string, len(tCombined.Columns))
	for k, v := range tCombined.Columns {
		columnSlice[v] = k
	}
	tCombined.ColumnSlice = columnSlice

	nrows := make([][]string, len(tBench.Rows))
	for i, row := range tBench.Rows {
		nrows[i] = append(row, tTest.Rows[i][1:]...)
	}
	tCombined.Rows = nrows

	if err := tCombined.ToCSV(combinedPath); err != nil {
		log.Fatal(err)
	}
}
