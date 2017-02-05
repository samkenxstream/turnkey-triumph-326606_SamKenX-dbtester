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

	"github.com/gyuho/dataframe"
)

// importBenchMetrics adds benchmark metrics from client-side
// and aggregates this to system metrics by unix timestamps.
func (data *analyzeData) importBenchMetrics(fpath string) (err error) {
	data.benchMetricsFilePath = fpath

	var tdf dataframe.Frame
	tdf, err = dataframe.NewFromCSV(nil, fpath)
	if err != nil {
		return
	}

	var oldTSCol dataframe.Column
	oldTSCol, err = tdf.Column("UNIX-SECOND")
	if err != nil {
		return err
	}

	// get first(minimum) unix second
	fv1, ok := oldTSCol.FrontNonNil()
	if !ok {
		return fmt.Errorf("FrontNonNil %s has empty Unix time %v", fpath, fv1)
	}
	ivv1, ok := fv1.Int64()
	if !ok {
		return fmt.Errorf("cannot Int64 %v", fv1)
	}
	data.benchMetrics.frontUnixSecond = int64(ivv1)

	// get last(maximum) unix second
	fv2, ok := oldTSCol.BackNonNil()
	if !ok {
		return fmt.Errorf("BackNonNil %s has empty Unix time %v", fpath, fv2)
	}
	ivv2, ok := fv2.Int64()
	if !ok {
		return fmt.Errorf("cannot Int64 %v", fv2)
	}
	data.benchMetrics.lastUnixSecond = int64(ivv2)

	// UNIX-SECOND, CONTROL-CLIENT-NUM, MIN-LATENCY-MS, AVG-LATENCY-MS, MAX-LATENCY-MS, AVG-THROUGHPUT
	var oldControlClientNumCol dataframe.Column
	oldControlClientNumCol, err = tdf.Column("CONTROL-CLIENT-NUM")
	if err != nil {
		return err
	}
	var oldMinLatencyMSCol dataframe.Column
	oldMinLatencyMSCol, err = tdf.Column("MIN-LATENCY-MS")
	if err != nil {
		return err
	}
	var oldAvgLatencyMSCol dataframe.Column
	oldAvgLatencyMSCol, err = tdf.Column("AVG-LATENCY-MS")
	if err != nil {
		return err
	}
	var oldMaxLatencyMSCol dataframe.Column
	oldMaxLatencyMSCol, err = tdf.Column("MAX-LATENCY-MS")
	if err != nil {
		return err
	}
	var oldAvgThroughputCol dataframe.Column
	oldAvgThroughputCol, err = tdf.Column("AVG-THROUGHPUT")
	if err != nil {
		return err
	}

	sec2Data := make(map[int64]rowData)
	for i := 0; i < oldTSCol.Count(); i++ {
		tv, err := oldTSCol.Value(i)
		if err != nil {
			return err
		}
		ts, ok := tv.Int64()
		if !ok {
			return fmt.Errorf("cannot Int64 %v", tv)
		}

		cv, err := oldControlClientNumCol.Value(i)
		if err != nil {
			return err
		}
		clientN, ok := cv.Int64()
		if !ok {
			return fmt.Errorf("cannot Int64 %v", cv)
		}
		cn := int64(clientN)

		lv1, err1 := oldMinLatencyMSCol.Value(i)
		if err1 != nil {
			return err1
		}
		minLat, ok1 := lv1.Float64()
		if !ok1 {
			return fmt.Errorf("cannot Float64 %v", lv1)
		}

		lv2, err2 := oldAvgLatencyMSCol.Value(i)
		if err2 != nil {
			return err2
		}
		avgLat, ok2 := lv2.Float64()
		if !ok2 {
			return fmt.Errorf("cannot Float64 %v", lv2)
		}

		lv3, err3 := oldMaxLatencyMSCol.Value(i)
		if err3 != nil {
			return err3
		}
		maxLat, ok3 := lv3.Float64()
		if !ok3 {
			return fmt.Errorf("cannot Float64 %v", lv3)
		}

		hv, err := oldAvgThroughputCol.Value(i)
		if err != nil {
			return err
		}
		dataThr, ok := hv.Float64()
		if !ok {
			return fmt.Errorf("cannot Float64 %v", hv)
		}

		if v, ok := sec2Data[ts]; !ok {
			sec2Data[ts] = rowData{clientN: cn, minLat: minLat, avgLat: avgLat, maxLat: maxLat, throughput: dataThr}
		} else {
			oldCn := v.clientN
			if oldCn != cn {
				return fmt.Errorf("different client number with same timestamps! %d != %d", oldCn, cn)
			}
			sec2Data[ts] = rowData{
				clientN:    cn,
				minLat:     minFloat64(v.minLat, minLat),
				avgLat:     (v.avgLat + avgLat) / 2.0,
				maxLat:     maxFloat64(v.maxLat, maxLat),
				throughput: (v.throughput + dataThr) / 2.0,
			}
		}
	}

	// UNIX-SECOND, CONTROL-CLIENT-NUM, MIN-LATENCY-MS, AVG-LATENCY-MS, MAX-LATENCY-MS, AVG-THROUGHPUT
	// aggregate duplicate benchmark timestamps with average values
	// OR fill in missing timestamps with zero values
	//
	// expected row number
	expectedRowN := data.benchMetrics.lastUnixSecond - data.benchMetrics.frontUnixSecond + 1
	newSecondCol := dataframe.NewColumn("UNIX-SECOND")
	newControlClientNumCol := dataframe.NewColumn("CONTROL-CLIENT-NUM")
	newMinLatencyCol := dataframe.NewColumn("MIN-LATENCY-MS")
	newAvgLatencyCol := dataframe.NewColumn("AVG-LATENCY-MS")
	newMaxLatencyCol := dataframe.NewColumn("MAX-LATENCY-MS")
	newAvgThroughputCol := dataframe.NewColumn("AVG-THROUGHPUT")
	for i := int64(0); i < expectedRowN; i++ {
		second := data.benchMetrics.frontUnixSecond + i
		newSecondCol.PushBack(dataframe.NewStringValue(second))

		v, ok := sec2Data[second]
		if !ok {
			prev := findClosest(second, sec2Data)
			newControlClientNumCol.PushBack(dataframe.NewStringValue(prev.clientN))
			newMinLatencyCol.PushBack(dataframe.NewStringValue(0.0))
			newAvgLatencyCol.PushBack(dataframe.NewStringValue(0.0))
			newMaxLatencyCol.PushBack(dataframe.NewStringValue(0.0))
			newAvgThroughputCol.PushBack(dataframe.NewStringValue(0))
			continue
		}

		newControlClientNumCol.PushBack(dataframe.NewStringValue(v.clientN))
		newMinLatencyCol.PushBack(dataframe.NewStringValue(v.minLat))
		newAvgLatencyCol.PushBack(dataframe.NewStringValue(v.avgLat))
		newMaxLatencyCol.PushBack(dataframe.NewStringValue(v.maxLat))
		newAvgThroughputCol.PushBack(dataframe.NewStringValue(v.throughput))
	}

	df := dataframe.New()
	if err = df.AddColumn(newSecondCol); err != nil {
		return err
	}
	if err = df.AddColumn(newControlClientNumCol); err != nil {
		return err
	}
	if err = df.AddColumn(newMinLatencyCol); err != nil {
		return err
	}
	if err = df.AddColumn(newAvgLatencyCol); err != nil {
		return err
	}
	if err = df.AddColumn(newMaxLatencyCol); err != nil {
		return err
	}
	if err = df.AddColumn(newAvgThroughputCol); err != nil {
		return err
	}

	data.benchMetrics.frame = df
	return
}

type rowData struct {
	clientN    int64
	minLat     float64
	avgLat     float64
	maxLat     float64
	throughput float64
}

func findClosest(second int64, sec2Data map[int64]rowData) rowData {
	v, ok := sec2Data[second]
	if ok {
		return v
	}
	var min int64
	var max int64
	for k := range sec2Data {
		if min == 0 || min > k {
			min = k
		}
		if max == 0 || max < k {
			max = k
		}
	}
	r, ok := _findClosestLower(second, sec2Data, min, max)
	if ok {
		return r
	}
	r, ok = _findClosestUpper(second, sec2Data, min, max)
	if !ok {
		panic(fmt.Errorf("something wrong with benchmark data... too many data points are missing"))
	}
	return r
}

func _findClosestUpper(second int64, sec2Data map[int64]rowData, min, max int64) (rowData, bool) {
	if second < min || second > max {
		return rowData{}, false
	}
	v, ok := sec2Data[second]
	if ok {
		return v, true
	}
	return _findClosestUpper(second+1, sec2Data, min, max)
}

func _findClosestLower(second int64, sec2Data map[int64]rowData, min, max int64) (rowData, bool) {
	if second < min || second > max {
		return rowData{}, false
	}
	v, ok := sec2Data[second]
	if ok {
		return v, true
	}
	return _findClosestLower(second-1, sec2Data, min, max)
}
