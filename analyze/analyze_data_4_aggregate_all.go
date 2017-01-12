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
	"strings"

	"github.com/gyuho/dataframe"
)

// aggregateAll aggregates all system metrics from 3+ nodes.
func (data *analyzeData) aggregateAll() error {
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
	// this is index, so decrement by 1 to make it as valid index
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
	data.aggregated = dataframe.New()

	// first, add bench metrics data
	// UNIX-TS, AVG-LATENCY-MS, AVG-THROUGHPUT
	for _, col := range data.benchMetrics.frame.Columns() {
		// ALWAYS KEEP FROM FIRST ROW OF BENCHMARKS
		if err = col.Keep(0, minBenchEndIdx); err != nil {
			return err
		}
		if err = data.aggregated.AddColumn(col); err != nil {
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
		if err = data.aggregated.AddColumn(col); err != nil {
			return err
		}
	}

	var (
		requestSum              int
		cumulativeThroughputCol = dataframe.NewColumn("CUMULATIVE-THROUGHPUT")

		sampleSize = float64(len(data.sys))

		avgCPUCol                   = dataframe.NewColumn("AVG-CPU")                      // from CPU-NUM
		avgVMRSSMBCol               = dataframe.NewColumn("AVG-VMRSS-MB")                 // from VMRSS-NUM
		avgReadsCompletedCol        = dataframe.NewColumn("AVG-READS-COMPLETED")          // from READS-COMPLETED
		avgReadsCompletedDeltaCol   = dataframe.NewColumn("AVG-READS-COMPLETED-DELTA")    // from READS-COMPLETED-DELTA
		avgSectorsReadCol           = dataframe.NewColumn("AVG-SECTORS-READ")             // from SECTORS-READ
		avgSectorsReadDeltaCol      = dataframe.NewColumn("AVG-SECTORS-READ-DELTA")       // from SECTORS-READ-DELTA
		avgWritesCompletedCol       = dataframe.NewColumn("AVG-WRITES-COMPLETED")         // from WRITES-COMPLETED
		avgWritesCompletedDeltaCol  = dataframe.NewColumn("AVG-WRITES-COMPLETED-DELTA")   // from WRITES-COMPLETED-DELTA
		avgSectorsWrittenCol        = dataframe.NewColumn("AVG-SECTORS-WRITTEN")          // from SECTORS-WRITTEN
		avgSectorsWrittenDeltaCol   = dataframe.NewColumn("AVG-SECTORS-WRITTEN-DELTA")    // from SECTORS-WRITTEN-DELTA
		avgReceiveBytesNumCol       = dataframe.NewColumn("AVG-RECEIVE-BYTES-NUM")        // from RECEIVE-BYTES-NUM
		avgReceiveBytesNumDeltaCol  = dataframe.NewColumn("AVG-RECEIVE-BYTES-NUM-DELTA")  // from RECEIVE-BYTES-NUM-DELTA
		avgTransmitBytesNumCol      = dataframe.NewColumn("AVG-TRANSMIT-BYTES-NUM")       // from TRANSMIT-BYTES-NUM
		avgTransmitBytesNumDeltaCol = dataframe.NewColumn("AVG-TRANSMIT-BYTES-NUM-DELTA") // from TRANSMIT-BYTES-NUM-DELTA
	)

	// compute average value of 3+ nodes
	// by iterating each row (horizontally) for all the columns
	for rowIdx := 0; rowIdx < minBenchEndIdx; rowIdx++ {
		var (
			cpuSum                   float64
			vmrssMBSum               float64
			readsCompletedSum        float64
			readsCompletedDeltaSum   float64
			sectorsReadSum           float64
			sectorsReadDeltaSum      float64
			writesCompletedSum       float64
			writesCompletedDeltaSum  float64
			sectorsWrittenSum        float64
			sectorsWrittenDeltaSum   float64
			receiveBytesNumSum       float64
			receiveBytesNumDeltaSum  float64
			transmitBytesNumSum      float64
			transmitBytesNumDeltaSum float64
		)
		for _, col := range data.aggregated.Columns() {
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
			case strings.HasPrefix(hd, "READS-COMPLETED-DELTA-"): // match this first!
				readsCompletedDeltaSum += vv
			case strings.HasPrefix(hd, "READS-COMPLETED-"):
				readsCompletedSum += vv
			case strings.HasPrefix(hd, "SECTORS-READ-DELTA-"):
				sectorsReadDeltaSum += vv
			case strings.HasPrefix(hd, "SECTORS-READ-"):
				sectorsReadSum += vv
			case strings.HasPrefix(hd, "WRITES-COMPLETED-DELTA-"):
				writesCompletedDeltaSum += vv
			case strings.HasPrefix(hd, "WRITES-COMPLETED-"):
				writesCompletedSum += vv
			case strings.HasPrefix(hd, "SECTORS-WRITTEN-DELTA-"):
				sectorsWrittenDeltaSum += vv
			case strings.HasPrefix(hd, "SECTORS-WRITTEN-"):
				sectorsWrittenSum += vv
			case strings.HasPrefix(hd, "RECEIVE-BYTES-NUM-DELTA-"):
				receiveBytesNumDeltaSum += vv
			case strings.HasPrefix(hd, "RECEIVE-BYTES-NUM-"):
				receiveBytesNumSum += vv
			case strings.HasPrefix(hd, "TRANSMIT-BYTES-NUM-DELTA-"):
				transmitBytesNumDeltaSum += vv
			case strings.HasPrefix(hd, "TRANSMIT-BYTES-NUM-"):
				transmitBytesNumSum += vv
			}
		}
		avgCPUCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", cpuSum/sampleSize)))
		avgVMRSSMBCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", vmrssMBSum/sampleSize)))
		avgReadsCompletedCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", readsCompletedSum/sampleSize)))
		avgReadsCompletedDeltaCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", readsCompletedDeltaSum/sampleSize)))
		avgSectorsReadCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", sectorsReadSum/sampleSize)))
		avgSectorsReadDeltaCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", sectorsReadDeltaSum/sampleSize)))
		avgWritesCompletedCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", writesCompletedSum/sampleSize)))
		avgWritesCompletedDeltaCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", writesCompletedDeltaSum/sampleSize)))
		avgSectorsWrittenCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", sectorsWrittenSum/sampleSize)))
		avgSectorsWrittenDeltaCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", sectorsWrittenDeltaSum/sampleSize)))
		avgReceiveBytesNumCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", receiveBytesNumSum/sampleSize)))
		avgReceiveBytesNumDeltaCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", receiveBytesNumDeltaSum/sampleSize)))
		avgTransmitBytesNumCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", transmitBytesNumSum/sampleSize)))
		avgTransmitBytesNumDeltaCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", transmitBytesNumDeltaSum/sampleSize)))
	}

	// move CLIENT-NUM to second column
	if err = data.aggregated.MoveColumn("CLIENT-NUM", 1); err != nil {
		return err
	}

	// add all cumulative, average columns
	if err = data.aggregated.AddColumn(cumulativeThroughputCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgCPUCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgVMRSSMBCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgReadsCompletedCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgReadsCompletedDeltaCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgSectorsReadCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgSectorsReadDeltaCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgWritesCompletedCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgWritesCompletedDeltaCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgSectorsWrittenCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgSectorsWrittenDeltaCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgReceiveBytesNumCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgReceiveBytesNumDeltaCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgTransmitBytesNumCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgTransmitBytesNumDeltaCol); err != nil {
		return err
	}

	// add SECOND column
	uc, err := data.aggregated.Column("UNIX-TS")
	if err != nil {
		return err
	}
	secondCol := dataframe.NewColumn("SECOND")
	for i := 0; i < uc.CountRow(); i++ {
		secondCol.PushBack(dataframe.NewStringValue(i))
	}
	if err = data.aggregated.AddColumn(secondCol); err != nil {
		return err
	}
	// move to 2nd column
	if err = data.aggregated.MoveColumn("SECOND", 1); err != nil {
		return err
	}

	// currently first columns are ordered as:
	// UNIX-TS, SECOND, CLIENT-NUM, AVG-LATENCY-MS, AVG-THROUGHPUT
	//
	// re-order columns in the following order, to make it more readable
	reorder := []string{
		"CUMULATIVE-THROUGHPUT",
		"AVG-CPU",
		"AVG-VMRSS-MB",
		"AVG-READS-COMPLETED-DELTA",
		"AVG-SECTORS-READ-DELTA",
		"AVG-WRITES-COMPLETED-DELTA",
		"AVG-SECTORS-WRITTEN-DELTA",
		"AVG-RECEIVE-BYTES-NUM-DELTA",
		"AVG-TRANSMIT-BYTES-NUM-DELTA",
		"AVG-READS-COMPLETED",
		"AVG-SECTORS-READ",
		"AVG-WRITES-COMPLETED",
		"AVG-SECTORS-WRITTEN",
		"AVG-RECEIVE-BYTES-NUM",
		"AVG-TRANSMIT-BYTES-NUM",
	}
	for i := len(reorder) - 1; i >= 0; i-- {
		if err = data.aggregated.MoveColumn(reorder[i], 5); err != nil {
			return err
		}
	}

	for _, col := range data.aggregated.Columns() {
		// since we will have same headers from different databases
		col.UpdateHeader(makeHeader(col.Header(), data.databaseTag))
	}
	return nil
}

func (data *analyzeData) save() error {
	return data.aggregated.CSV(data.csvOutputpath)
}

func makeHeader(column string, tag string) string {
	return fmt.Sprintf("%s-%s", column, tag)
}
