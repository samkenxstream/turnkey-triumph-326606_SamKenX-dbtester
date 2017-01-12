package analyze

import (
	"fmt"

	"strings"

	"github.com/gyuho/dataframe"
)

// aggSystemBenchMetrics aggregates all system metrics from 3+ nodes.
func (data *analyzeData) aggSystemBenchMetrics() error {
	plog.Println("STEP #3: aggregating system metrics and benchmark metrics")
	colSys, err := data.sysAgg.Column("UNIX-TS")
	if err != nil {
		return err
	}

	colBench, err := data.benchMetrics.frame.Column("UNIX-TS")
	if err != nil {
		return err
	}
	fv, ok := colBench.FrontNonNil()
	if !ok {
		return fmt.Errorf("FrontNonNil %s has empty Unix time %v", data.benchMetrics.filePath, fv)
	}
	bv, ok := colBench.BackNonNil()
	if !ok {
		return fmt.Errorf("BackNonNil %s has empty Unix time %v", data.benchMetrics.filePath, fv)
	}

	sysStartIdx, ok := colSys.FindFirst(fv)
	if !ok {
		return fmt.Errorf("%v is not found in system metrics results", fv)
	}
	sysEndIdx, ok := colSys.FindFirst(bv)
	if !ok {
		return fmt.Errorf("%v is not found in system metrics results", fv)
	}
	sysRowN := sysEndIdx - sysStartIdx + 1

	var minBenchEndIdx int
	for _, col := range data.benchMetrics.frame.Columns() {
		if minBenchEndIdx == 0 {
			minBenchEndIdx = col.CountRow()
		}
		if minBenchEndIdx > col.CountRow() {
			minBenchEndIdx = col.CountRow()
		}
	}
	minBenchEndIdx--

	// sysStartIdx 3, sysEndIdx 9, sysRowN 7, minBenchEndIdx 5 (5+1 < 7)
	// THEN sysEndIdx = 3 + 5 = 8
	//
	// sysStartIdx 3, sysEndIdx 7, sysRowN 5, minBenchEndIdx 5 (5+1 > 5)
	// THEN minBenchEndIdx = 7 - 3 = 4
	if minBenchEndIdx+1 < sysRowN {
		// benchmark is short of rows
		// adjust system-metrics rows to benchmark-metrics
		sysEndIdx = sysStartIdx + minBenchEndIdx
	} else {
		// system-metrics is short of rows
		// adjust benchmark-metrics to system-metrics
		minBenchEndIdx = sysEndIdx - sysStartIdx
	}

	// aggregate all system-metrics and benchmark-metrics
	data.sysBenchAgg = dataframe.New()

	// first, add bench metrics data
	// UNIX-TS, AVG-LATENCY-MS, AVG-THROUGHPUT
	for _, col := range data.benchMetrics.frame.Columns() {
		// ALWAYS KEEP FROM FIRST ROW OF BENCHMARKS
		if err = col.Keep(0, minBenchEndIdx); err != nil {
			return err
		}
		if err = data.sysBenchAgg.AddColumn(col); err != nil {
			return err
		}
	}

	for _, col := range data.sysAgg.Columns() {
		if col.Header() == "UNIX-TS" {
			continue
		}
		if err = col.Keep(sysStartIdx, sysEndIdx); err != nil {
			return err
		}
		if err = data.sysBenchAgg.AddColumn(col); err != nil {
			return err
		}
	}

	plog.Println("STEP #4: computing average,cumulative values in system metrics and benchmark")
	var (
		requestSum              int
		cumulativeThroughputCol = dataframe.NewColumn("CUMULATIVE-THROUGHPUT")

		systemMetricsSize = float64(len(data.sys))

		avgCPUCol     = dataframe.NewColumn("AVG-CPU")
		avgVMRSSMBCol = dataframe.NewColumn("AVG-VMRSS-MB")
	)

	// compute average value of 3+ nodes
	// by iterating each row (horizontally) for all the columns
	for rowIdx := 0; rowIdx < minBenchEndIdx; rowIdx++ {
		var (
			cpuSum     float64
			vmrssMBSum float64
		)
		for _, col := range data.sysBenchAgg.Columns() {
			rv, err := col.Value(rowIdx)
			if err != nil {
				return err
			}
			vv, _ := rv.Number()

			hd := col.Header()
			switch {
			// cumulative values
			case hd == "AVG-THROUGHPUT":
				requestSum += int(vv)
				cumulativeThroughputCol.PushBack(dataframe.NewStringValue(requestSum))

			// average values (need sume first!)
			case strings.HasPrefix(hd, "CPU-"):
				// CPU-NUM was converted to CPU-1, CPU-2, CPU-3
				cpuSum += vv
			case strings.HasPrefix(hd, "VMRSS-MB-"):
				// VMRSS-NUM-NUM was converted to VMRSS-MB-1, VMRSS-MB-2, VMRSS-MB-3
				vmrssMBSum += vv
			case strings.HasPrefix(hd, "CPU-"):
			case strings.HasPrefix(hd, "CPU-"):
			case strings.HasPrefix(hd, "CPU-"):
			case strings.HasPrefix(hd, "CPU-"):
			case strings.HasPrefix(hd, "CPU-"):
			case strings.HasPrefix(hd, "CPU-"):
			case strings.HasPrefix(hd, "CPU-"):
			case strings.HasPrefix(hd, "CPU-"):
			}
		}
	}

	return nil
}
