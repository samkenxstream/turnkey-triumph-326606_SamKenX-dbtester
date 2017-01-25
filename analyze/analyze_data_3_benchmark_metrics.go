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
	oldTSCol, err = tdf.Column("UNIX-TS")
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
	data.benchMetrics.frontUnixTS = int64(ivv1)

	// get last(maximum) unix second
	fv2, ok := oldTSCol.BackNonNil()
	if !ok {
		return fmt.Errorf("BackNonNil %s has empty Unix time %v", fpath, fv2)
	}
	ivv2, ok := fv2.Int64()
	if !ok {
		return fmt.Errorf("cannot Int64 %v", fv2)
	}
	data.benchMetrics.lastUnixTS = int64(ivv2)

	// UNIX-TS, CONTROL-CLIENT-NUM, AVG-LATENCY-MS, AVG-THROUGHPUT
	var oldControlClientNumCol dataframe.Column
	oldControlClientNumCol, err = tdf.Column("CONTROL-CLIENT-NUM")
	if err != nil {
		return err
	}
	var oldAvgLatencyMSCol dataframe.Column
	oldAvgLatencyMSCol, err = tdf.Column("AVG-LATENCY-MS")
	if err != nil {
		return err
	}
	var oldAvgThroughputCol dataframe.Column
	oldAvgThroughputCol, err = tdf.Column("AVG-THROUGHPUT")
	if err != nil {
		return err
	}

	type rowData struct {
		clientN    int64
		latency    float64
		throughput float64
	}
	tsToData := make(map[int64]rowData)
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

		lv, err := oldAvgLatencyMSCol.Value(i)
		if err != nil {
			return err
		}
		dataLat, ok := lv.Float64()
		if !ok {
			return fmt.Errorf("cannot Float64 %v", lv)
		}

		hv, err := oldAvgThroughputCol.Value(i)
		if err != nil {
			return err
		}
		dataThr, ok := hv.Float64()
		if !ok {
			return fmt.Errorf("cannot Float64 %v", hv)
		}

		if v, ok := tsToData[ts]; !ok {
			tsToData[ts] = rowData{clientN: cn, latency: dataLat, throughput: dataThr}
		} else {
			oldCn := v.clientN
			if oldCn != cn {
				return fmt.Errorf("different client number with same timestamps! %d != %d", oldCn, cn)
			}
			tsToData[ts] = rowData{clientN: cn, latency: (v.latency + dataLat) / 2.0, throughput: (v.throughput + dataThr) / 2.0}
		}
	}

	// UNIX-TS, CONTROL-CLIENT-NUM, AVG-LATENCY-MS, AVG-THROUGHPUT
	// aggregate duplicate benchmark timestamps with average values
	// OR fill in missing timestamps with zero values
	//
	// expected row number
	rowN := data.benchMetrics.lastUnixTS - data.benchMetrics.frontUnixTS + 1
	newTSCol := dataframe.NewColumn("UNIX-TS")
	newControlClientNumCol := dataframe.NewColumn("CONTROL-CLIENT-NUM")
	newAvgLatencyCol := dataframe.NewColumn("AVG-LATENCY-MS")
	newAvgThroughputCol := dataframe.NewColumn("AVG-THROUGHPUT")
	for i := int64(0); i < rowN; i++ {
		ts := data.benchMetrics.frontUnixTS + i
		newTSCol.PushBack(dataframe.NewStringValue(fmt.Sprintf("%d", ts)))

		v, ok := tsToData[ts]
		if !ok {
			prev, pok := tsToData[ts-1]
			if !pok {
				prev, pok = tsToData[ts+1]
				if !pok {
					return fmt.Errorf("benchmark missing a lot of rows around %d", ts)
				}
			}
			newControlClientNumCol.PushBack(dataframe.NewStringValue(prev.clientN))

			// just add empty values
			newAvgLatencyCol.PushBack(dataframe.NewStringValue("0.0"))
			newAvgThroughputCol.PushBack(dataframe.NewStringValue(0))
		} else {
			newControlClientNumCol.PushBack(dataframe.NewStringValue(v.clientN))
			newAvgLatencyCol.PushBack(dataframe.NewStringValue(v.latency))
			newAvgThroughputCol.PushBack(dataframe.NewStringValue(v.throughput))
		}
	}

	df := dataframe.New()
	if err = df.AddColumn(newTSCol); err != nil {
		return err
	}
	if err = df.AddColumn(newControlClientNumCol); err != nil {
		return err
	}
	if err = df.AddColumn(newAvgLatencyCol); err != nil {
		return err
	}
	if err = df.AddColumn(newAvgThroughputCol); err != nil {
		return err
	}

	data.benchMetrics.frame = df
	return
}
