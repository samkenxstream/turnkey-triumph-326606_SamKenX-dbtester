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
	"sort"
	"strings"

	"github.com/gyuho/dataframe"
)

// aggregateAll aggregates all system metrics from 3+ nodes.
func (data *analyzeData) aggregateAll(memoryByKeyPath string, totalRequests int64) error {
	colSys, err := data.sysAgg.Column("UNIX-SECOND")
	if err != nil {
		return err
	}
	colBench, err := data.benchMetrics.frame.Column("UNIX-SECOND")
	if err != nil {
		return err
	}

	fv, ok := colBench.FrontNonNil()
	if !ok {
		return fmt.Errorf("FrontNonNil %s has empty Unix time %v", data.benchMetrics.filePath, fv)
	}
	sysStartIdx, ok := colSys.FindFirst(fv)
	if !ok {
		return fmt.Errorf("%v is not found in system metrics results", fv)
	}

	bv, ok := colBench.BackNonNil()
	if !ok {
		return fmt.Errorf("BackNonNil %s has empty Unix time %v", data.benchMetrics.filePath, fv)
	}
	sysEndIdx, ok := colSys.FindFirst(bv)
	if !ok {
		return fmt.Errorf("%v is not found in system metrics results", fv)
	}

	expectedSysRowN := sysEndIdx - sysStartIdx + 1

	var minBenchEndIdx int
	for _, col := range data.benchMetrics.frame.Columns() {
		if minBenchEndIdx == 0 {
			minBenchEndIdx = col.Count()
		}
		if minBenchEndIdx > col.Count() {
			minBenchEndIdx = col.Count()
		}
	}
	// this is index, so decrement by 1 to make it as valid index
	minBenchEndIdx--

	// sysStartIdx 3, sysEndIdx 9, expectedSysRowN 7, minBenchEndIdx 5 (5+1 < 7)
	// so benchmark has 6 rows, but system metrics has 7 rows; benchmark is short of rows
	// so we should keep system metrics [3, 9)
	// THEN sysEndIdx = 3 + 5 = 8 (keep [3, 8+1))
	//
	// sysStartIdx 3, sysEndIdx 7, expectedSysRowN 5, minBenchEndIdx 5 (5+1 > 5)
	// so benchmark has 6 rows, but system metrics has 5 rows; system metrics is short of rows
	// so we should keep benchmark [0, 5)
	// THEN minBenchEndIdx = 7 - 3 = 4 (keep [0, 4+1))
	if minBenchEndIdx+1 < expectedSysRowN {
		// benchmark is short of rows
		// adjust system metrics rows to benchmark-metrics
		// will truncate front of system metrics rows
		sysEndIdx = sysStartIdx + minBenchEndIdx
	} else {
		// system metrics is short of rows
		// adjust benchmark metrics to system-metrics
		// will truncate front of benchmark metrics rows
		minBenchEndIdx = sysEndIdx - sysStartIdx
	}

	// aggregate all system-metrics and benchmark-metrics
	data.aggregated = dataframe.New()

	// first, add bench metrics data
	// UNIX-SECOND, MIN-LATENCY-MS, AVG-LATENCY-MS, MAX-LATENCY-MS, AVG-THROUGHPUT
	for _, col := range data.benchMetrics.frame.Columns() {
		// ALWAYS KEEP FROM FIRST ROW OF BENCHMARKS
		// keeps from [a, b)
		if err = col.Keep(0, minBenchEndIdx+1); err != nil {
			return err
		}
		if err = data.aggregated.AddColumn(col); err != nil {
			return err
		}
	}
	for _, col := range data.sysAgg.Columns() {
		if col.Header() == "UNIX-SECOND" {
			continue
		}
		if err = col.Keep(sysStartIdx, sysEndIdx+1); err != nil {
			return err
		}
		if err = data.aggregated.AddColumn(col); err != nil {
			return err
		}
	}

	var (
		requestSum              int
		cumulativeThroughputCol = dataframe.NewColumn("CUMULATIVE-THROUGHPUT") // from AVG-THROUGHPUT

		sampleSize = float64(len(data.sys))

		avgClientNumCol             = dataframe.NewColumn("AVG-CLIENT-NUM")                  // from CLIENT-NUM
		avgVolCtxSwitchCol          = dataframe.NewColumn("AVG-VOLUNTARY-CTXT-SWITCHES")     // from VOLUNTARY-CTXT-SWITCHES
		avgNonVolCtxSwitchCol       = dataframe.NewColumn("AVG-NON-VOLUNTARY-CTXT-SWITCHES") // from NON-VOLUNTARY-CTXT-SWITCHES
		avgCPUCol                   = dataframe.NewColumn("AVG-CPU")                         // from CPU-NUM
		avgSystemLoadCol            = dataframe.NewColumn("AVG-SYSTEM-LOAD-1-MIN")           // from LOAD-AVERAGE-1-MINUTE
		avgVMRSSMBCol               = dataframe.NewColumn("AVG-VMRSS-MB")                    // from VMRSS-NUM
		avgReadsCompletedCol        = dataframe.NewColumn("AVG-READS-COMPLETED")             // from READS-COMPLETED
		avgReadsCompletedDeltaCol   = dataframe.NewColumn("AVG-READS-COMPLETED-DELTA")       // from READS-COMPLETED-DELTA
		avgSectorsReadCol           = dataframe.NewColumn("AVG-SECTORS-READ")                // from SECTORS-READ
		avgSectorsReadDeltaCol      = dataframe.NewColumn("AVG-SECTORS-READ-DELTA")          // from SECTORS-READ-DELTA
		avgWritesCompletedCol       = dataframe.NewColumn("AVG-WRITES-COMPLETED")            // from WRITES-COMPLETED
		avgWritesCompletedDeltaCol  = dataframe.NewColumn("AVG-WRITES-COMPLETED-DELTA")      // from WRITES-COMPLETED-DELTA
		avgSectorsWrittenCol        = dataframe.NewColumn("AVG-SECTORS-WRITTEN")             // from SECTORS-WRITTEN
		avgSectorsWrittenDeltaCol   = dataframe.NewColumn("AVG-SECTORS-WRITTEN-DELTA")       // from SECTORS-WRITTEN-DELTA
		avgReceiveBytesNumCol       = dataframe.NewColumn("AVG-RECEIVE-BYTES-NUM")           // from RECEIVE-BYTES-NUM
		avgReceiveBytesNumDeltaCol  = dataframe.NewColumn("AVG-RECEIVE-BYTES-NUM-DELTA")     // from RECEIVE-BYTES-NUM-DELTA
		avgTransmitBytesNumCol      = dataframe.NewColumn("AVG-TRANSMIT-BYTES-NUM")          // from TRANSMIT-BYTES-NUM
		avgTransmitBytesNumDeltaCol = dataframe.NewColumn("AVG-TRANSMIT-BYTES-NUM-DELTA")    // from TRANSMIT-BYTES-NUM-DELTA
	)

	sec2minVMRSSMB := make(map[int64]float64)
	sec2maxVMRSSMB := make(map[int64]float64)

	// compute average value of 3+ nodes
	// by iterating each row (horizontally) for all the columns
	for rowIdx := 0; rowIdx < minBenchEndIdx+1; rowIdx++ {
		var (
			clientNumSum             float64
			volCtxSwitchSum          float64
			nonVolCtxSwitchSum       float64
			cpuSum                   float64
			loadAvgSum               float64
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
		sc, err := data.aggregated.Column("UNIX-SECOND")
		if err != nil {
			return err
		}
		for _, col := range data.aggregated.Columns() {
			rv, err := col.Value(rowIdx)
			if err != nil {
				return err
			}
			vv, _ := rv.Float64()

			hd := col.Header()
			switch {
			// cumulative values
			case hd == "AVG-THROUGHPUT":
				requestSum += int(vv)
				cumulativeThroughputCol.PushBack(dataframe.NewStringValue(requestSum))

			// average values (need sume first!)
			case strings.HasPrefix(hd, "CLIENT-NUM-"):
				clientNumSum += vv
			case strings.HasPrefix(hd, "VOLUNTARY-CTXT-SWITCHES-"):
				volCtxSwitchSum += vv
			case strings.HasPrefix(hd, "NON-VOLUNTARY-CTXT-SWITCHES-"):
				nonVolCtxSwitchSum += vv
			case strings.HasPrefix(hd, "CPU-"): // CPU-NUM was converted to CPU-1, CPU-2, CPU-3
				cpuSum += vv
			case strings.HasPrefix(hd, "LOAD-AVERAGE-1-"): // LOAD-AVERAGE-1-MINUTE
				loadAvgSum += vv
			case strings.HasPrefix(hd, "VMRSS-MB-"): // VMRSS-NUM-NUM was converted to VMRSS-MB-1, VMRSS-MB-2, VMRSS-MB-3
				vmrssMBSum += vv

				svv, err := sc.Value(rowIdx)
				if err != nil {
					return err
				}
				ts, _ := svv.Int64()

				if v, ok := sec2minVMRSSMB[ts]; !ok {
					sec2minVMRSSMB[ts] = vv
				} else if v > vv || (v == 0.0 && vv != 0.0) {
					sec2minVMRSSMB[ts] = vv
				}
				if v, ok := sec2maxVMRSSMB[ts]; !ok {
					sec2maxVMRSSMB[ts] = vv
				} else if v < vv {
					sec2maxVMRSSMB[ts] = vv
				}

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

		avgClientNumCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", clientNumSum/sampleSize)))
		avgVolCtxSwitchCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", volCtxSwitchSum/sampleSize)))
		avgNonVolCtxSwitchCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", nonVolCtxSwitchSum/sampleSize)))
		avgCPUCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", cpuSum/sampleSize)))
		avgSystemLoadCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", loadAvgSum/sampleSize)))
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

	// add all cumulative, average columns
	if err = data.aggregated.AddColumn(cumulativeThroughputCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgClientNumCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgVolCtxSwitchCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgNonVolCtxSwitchCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgCPUCol); err != nil {
		return err
	}
	if err = data.aggregated.AddColumn(avgSystemLoadCol); err != nil {
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
	uc, err := data.aggregated.Column("UNIX-SECOND")
	if err != nil {
		return err
	}
	secondCol := dataframe.NewColumn("SECOND")
	for i := 0; i < uc.Count(); i++ {
		secondCol.PushBack(dataframe.NewStringValue(i))
	}
	if err = data.aggregated.AddColumn(secondCol); err != nil {
		return err
	}
	// move to 2th column
	if err = data.aggregated.MoveColumn("SECOND", 1); err != nil {
		return err
	}
	// move to 3th column
	if err = data.aggregated.MoveColumn("CONTROL-CLIENT-NUM", 2); err != nil {
		return err
	}
	// move to 3th column
	if err = data.aggregated.MoveColumn("AVG-CLIENT-NUM", 2); err != nil {
		return err
	}
	// move to 6th column
	if err = data.aggregated.MoveColumn("AVG-THROUGHPUT", 5); err != nil {
		return err
	}

	// currently first columns are ordered as:
	// UNIX-SECOND, SECOND, AVG-CLIENT-NUM, MIN-LATENCY-MS, AVG-LATENCY-MS, MAX-LATENCY-MS, AVG-THROUGHPUT
	//
	// re-order columns in the following order, to make it more readable
	reorder := []string{
		"CUMULATIVE-THROUGHPUT",
		"AVG-CPU",
		"AVG-SYSTEM-LOAD-1-MIN",
		"AVG-VMRSS-MB",
		"AVG-WRITES-COMPLETED",
		"AVG-WRITES-COMPLETED-DELTA",
		"AVG-SECTORS-WRITTEN",
		"AVG-SECTORS-WRITTEN-DELTA",
		"AVG-READS-COMPLETED",
		"AVG-READS-COMPLETED-DELTA",
		"AVG-SECTORS-READ",
		"AVG-SECTORS-READ-DELTA",
		"AVG-RECEIVE-BYTES-NUM",
		"AVG-RECEIVE-BYTES-NUM-DELTA",
		"AVG-TRANSMIT-BYTES-NUM",
		"AVG-TRANSMIT-BYTES-NUM-DELTA",
		"AVG-VOLUNTARY-CTXT-SWITCHES",
		"AVG-NON-VOLUNTARY-CTXT-SWITCHES",
	}
	// move to 9th
	for i := len(reorder) - 1; i >= 0; i-- {
		if err = data.aggregated.MoveColumn(reorder[i], 8); err != nil {
			return err
		}
	}

	for _, col := range data.aggregated.Columns() {
		// since we will have same headers from different databases
		col.UpdateHeader(makeHeader(col.Header(), data.databaseTag))
	}

	// aggregate memory usage by number of keys
	colUnixSecond, err := data.aggregated.Column("UNIX-SECOND")
	if err != nil {
		return err
	}
	colMemoryMB, err := data.aggregated.Column("AVG-VMRSS-MB")
	if err != nil {
		return err
	}
	colAvgThroughput, err := data.aggregated.Column("AVG-THROUGHPUT")
	if err != nil {
		return err
	}
	if colUnixSecond.Count() != colMemoryMB.Count() {
		return fmt.Errorf("SECOND column count %d, AVG-VMRSS-MB column count %d", colUnixSecond.Count(), colMemoryMB.Count())
	}
	if colAvgThroughput.Count() != colMemoryMB.Count() {
		return fmt.Errorf("AVG-THROUGHPUT column count %d, AVG-VMRSS-MB column count %d", colAvgThroughput.Count(), colMemoryMB.Count())
	}
	if colUnixSecond.Count() != colAvgThroughput.Count() {
		return fmt.Errorf("SECOND column count %d, AVG-THROUGHPUT column count %d", colUnixSecond.Count(), colAvgThroughput.Count())
	}

	var tslice []keyNumAndMemory
	for i := 0; i < colUnixSecond.Count(); i++ {
		vv0, err := colUnixSecond.Value(i)
		if err != nil {
			return err
		}
		v0, _ := vv0.Int64()

		vv1, err := colMemoryMB.Value(i)
		if err != nil {
			return err
		}
		vf1, _ := vv1.Float64()

		vv2, err := colAvgThroughput.Value(i)
		if err != nil {
			return err
		}
		vf2, _ := vv2.Float64()

		point := keyNumAndMemory{
			keyNum:      int64(vf2),
			minMemoryMB: sec2minVMRSSMB[v0],
			avgMemoryMB: vf1,
			maxMemoryMB: sec2maxVMRSSMB[v0],
		}
		tslice = append(tslice, point)
	}
	sort.Sort(keyNumAndMemorys(tslice))

	// aggregate memory by number of keys
	knms := processTimeSeries(tslice, 1000, totalRequests)
	ckk1 := dataframe.NewColumn("KEYS")
	ckk2 := dataframe.NewColumn("MIN-VMRSS-MB")
	ckk3 := dataframe.NewColumn("AVG-VMRSS-MB")
	ckk4 := dataframe.NewColumn("MAX-VMRSS-MB")
	for i := range knms {
		ckk1.PushBack(dataframe.NewStringValue(knms[i].keyNum))
		ckk2.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", knms[i].minMemoryMB)))
		ckk3.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", knms[i].avgMemoryMB)))
		ckk4.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.2f", knms[i].maxMemoryMB)))
	}
	fr := dataframe.New()
	if err := fr.AddColumn(ckk1); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(ckk2); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(ckk3); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(ckk4); err != nil {
		plog.Fatal(err)
	}
	if err := fr.CSV(memoryByKeyPath); err != nil {
		plog.Fatal(err)
	}

	return nil
}

func (data *analyzeData) save() error {
	return data.aggregated.CSV(data.allAggregatedOutputPath)
}
