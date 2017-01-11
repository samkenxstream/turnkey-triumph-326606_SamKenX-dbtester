package analyze

import (
	"fmt"

	"strings"

	"github.com/gyuho/dataframe"
)

// aggSystemBenchMetrics aggregates all system metrics from 3+ nodes.
func (data *analyzeData) aggSystemBenchMetrics() error {
	plog.Println("STEP #3: aggregating system metrics and benchmark metrics")
	colSys, err := data.sysAgg.GetColumn("UNIX-TS")
	if err != nil {
		return err
	}

	colBench, err := data.benchMetrics.frame.GetColumn("UNIX-TS")
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

	sysStartIdx, ok := colSys.FindValue(fv)
	if !ok {
		return fmt.Errorf("%v is not found in system metrics results", fv)
	}
	sysEndIdx, ok := colSys.FindValue(bv)
	if !ok {
		return fmt.Errorf("%v is not found in system metrics results", fv)
	}
	sysRowN := sysEndIdx - sysStartIdx + 1

	var minBenchEndIdx int
	for _, col := range data.benchMetrics.frame.GetColumns() {
		if minBenchEndIdx == 0 {
			minBenchEndIdx = col.RowNumber()
		}
		if minBenchEndIdx > col.RowNumber() {
			minBenchEndIdx = col.RowNumber()
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
	for _, col := range data.benchMetrics.frame.GetColumns() {
		// ALWAYS KEEP FROM FIRST ROW OF BENCHMARKS
		if err = col.KeepRows(0, minBenchEndIdx); err != nil {
			return err
		}
		if err = data.sysBenchAgg.AddColumn(col); err != nil {
			return err
		}
	}

	for _, col := range data.sysAgg.GetColumns() {
		if col.GetHeader() == "UNIX-TS" {
			continue
		}
		if err = col.KeepRows(sysStartIdx, sysEndIdx); err != nil {
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
		avgCPUCol         = dataframe.NewColumn("AVG-CPU")
		avgVMRSSMBCol     = dataframe.NewColumn("AVG-VMRSS-MB")
	)
	// iterate horizontally across all the columns
	for rowIdx := 0; rowIdx < minBenchEndIdx; rowIdx++ {
		var (
			cpuSum     float64
			vmrssMBSum float64
		)
		for _, col := range data.sysBenchAgg.GetColumns() {
			rv, err := col.GetValue(rowIdx)
			if err != nil {
				return err
			}
			vv, _ := rv.ToNumber()

			switch {
			case col.GetHeader() == "AVG-THROUGHPUT":
				requestSum += int(vv)
				cumulativeThroughputCol.PushBack(dataframe.NewStringValue(requestSum))

			case strings.HasPrefix(col.GetHeader(), "CPU-"):
				cpuSum += vv

			case strings.HasPrefix(col.GetHeader(), "VMRSS-MB-"):
				vmrssMBSum += vv
			}
		}
	}

	return nil
}
