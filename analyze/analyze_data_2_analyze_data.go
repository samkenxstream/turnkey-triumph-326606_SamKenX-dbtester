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

type analyzeData struct {
	databaseTag string
	legend      string

	minUnixSecond int64
	maxUnixSecond int64
	sys           []testData

	// aggregated frame within [min,maxUnixSecond] from sys
	sysAgg               dataframe.Frame
	benchMetricsFilePath string
	benchMetrics         testData

	// aggregated from sysAgg and benchMetrics
	aggregated dataframe.Frame

	csvOutputpath string
}

// readSystemMetricsAll reads all system metric files
// (e.g. if cluster is 3-node, read all 3 files).
// It returns minimum and maximum common unix second and a list of frames.
func readSystemMetricsAll(fpaths ...string) (data *analyzeData, err error) {
	data = &analyzeData{}
	for i, fpath := range fpaths {
		sm, err := readSystemMetrics(fpath)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			data.minUnixSecond = sm.frontUnixSecond
			data.maxUnixSecond = sm.lastUnixSecond
		}
		if data.minUnixSecond < sm.frontUnixSecond {
			data.minUnixSecond = sm.frontUnixSecond
		}
		if data.maxUnixSecond > sm.lastUnixSecond {
			data.maxUnixSecond = sm.lastUnixSecond
		}
		data.sys = append(data.sys, sm)
	}
	return
}

// aggSystemMetrics aggregates all system metrics from 3+ nodes.
func (data *analyzeData) aggSystemMetrics() error {
	// monitor CSVs from multiple servers, and want them to have equal number of rows
	// Truncate all rows before data.minUnixSecond and after data.maxUnixSecond
	minTS := fmt.Sprintf("%d", data.minUnixSecond)
	maxTS := fmt.Sprintf("%d", data.maxUnixSecond)
	data.sysAgg = dataframe.New()
	for i := range data.sys {
		uc, err := data.sys[i].frame.Column("UNIX-SECOND")
		if err != nil {
			return err
		}
		minTSIdx, ok := uc.FindFirst(dataframe.NewStringValue(minTS))
		if !ok {
			return fmt.Errorf("%v does not exist in %s", minTS, data.sys[i].filePath)
		}
		maxTSIdx, ok := uc.FindFirst(dataframe.NewStringValue(maxTS))
		if !ok {
			return fmt.Errorf("%v does not exist in %s", maxTS, data.sys[i].filePath)
		}

		for _, header := range data.sys[i].frame.Headers() {
			if i > 0 && header == "UNIX-SECOND" {
				// skip for other databases; we want to keep just one UNIX-SECOND column
				continue
			}

			var col dataframe.Column
			col, err = data.sys[i].frame.Column(header)
			if err != nil {
				return err
			}
			// just keep rows from [min,maxUnixSecond]
			if err = col.Keep(minTSIdx, maxTSIdx+1); err != nil {
				return err
			}

			if header == "UNIX-SECOND" {
				if err = data.sysAgg.AddColumn(col); err != nil {
					return err
				}
				continue
			}

			switch header {
			case "CPU-NUM":
				header = "CPU"

			case "VMRSS-NUM":
				header = "VMRSS-MB"

				// convert bytes to mb
				colN := col.Count()
				for rowIdx := 0; rowIdx < colN; rowIdx++ {
					var rowV dataframe.Value
					rowV, err = col.Value(rowIdx)
					if err != nil {
						return err
					}
					fv, _ := rowV.Float64()
					frv := float64(fv) * 0.000001
					if err = col.Set(rowIdx, dataframe.NewStringValue(fmt.Sprintf("%.2f", frv))); err != nil {
						return err
					}
				}

			case "EXTRA":
				// dbtester uses psn 'EXTRA' column as 'CLIENT-NUM'
				header = "CLIENT-NUM"
			}

			// since we are aggregating multiple system-metrics CSV files
			// suffix header with the index
			col.UpdateHeader(fmt.Sprintf("%s-%d", header, i+1))
			if err = data.sysAgg.AddColumn(col); err != nil {
				return err
			}
		}
	}

	return nil
}
