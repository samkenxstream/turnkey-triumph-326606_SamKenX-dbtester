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
	minUnixTS int64
	maxUnixTS int64
	sys       []testData

	// aggregated frame within [min,maxUnixTS] from sys
	sysAgg dataframe.Frame

	benchMetricsFilePath string
	benchMetrics         testData

	// aggregated from sysAgg and benchMetrics
	sysBenchAgg dataframe.Frame
}

// readSystemMetricsAll reads all system metric files
// (e.g. if cluster is 3-node, read all 3 files).
// It returns minimum and maximum common unix second and a list of frames.
func readSystemMetricsAll(fpaths ...string) (data *analyzeData, err error) {
	data = &analyzeData{}
	for i, fpath := range fpaths {
		plog.Printf("STEP #1-%d: creating dataframe from %s", i, fpath)
		sm, err := readSystemMetrics(fpath)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			data.minUnixTS = sm.frontUnixTS
			data.maxUnixTS = sm.lastUnixTS
		}
		if data.minUnixTS < sm.frontUnixTS {
			data.minUnixTS = sm.frontUnixTS
		}
		if data.maxUnixTS > sm.lastUnixTS {
			data.maxUnixTS = sm.lastUnixTS
		}
		data.sys = append(data.sys, sm)
	}
	return
}

// aggSystemMetrics aggregates all system metrics from 3+ nodes.
func (data *analyzeData) aggSystemMetrics() error {
	// monitor CSVs from multiple servers, and want them to have equal number of rows
	// Truncate all rows before data.minUnixTS and after data.maxUnixTS
	minTS := fmt.Sprintf("%d", data.minUnixTS)
	maxTS := fmt.Sprintf("%d", data.maxUnixTS)
	data.sysAgg = dataframe.New()
	for i := range data.sys {
		uc, err := data.sys[i].frame.GetColumn("UNIX-TS")
		if err != nil {
			return err
		}
		minTSIdx, ok := uc.FindValue(dataframe.NewStringValue(minTS))
		if !ok {
			return fmt.Errorf("%v does not exist in %s", minTS, data.sys[i].filePath)
		}
		maxTSIdx, ok := uc.FindValue(dataframe.NewStringValue(maxTS))
		if !ok {
			return fmt.Errorf("%v does not exist in %s", maxTS, data.sys[i].filePath)
		}

		for _, header := range data.sys[i].frame.GetHeader() {
			if i > 0 && header == "UNIX-TS" {
				// skip for other databases; we want to keep just one UNIX-TS column
				continue
			}

			var col dataframe.Column
			col, err = data.sys[i].frame.GetColumn(header)
			if err != nil {
				return err
			}
			// just keep rows from [min,maxUnixTS]
			if err = col.KeepRows(minTSIdx, maxTSIdx+1); err != nil {
				return err
			}

			if header == "UNIX-TS" {
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
				colN := col.RowNumber()
				for rowIdx := 0; rowIdx < colN; rowIdx++ {
					var rowV dataframe.Value
					rowV, err = col.GetValue(rowIdx)
					if err != nil {
						return err
					}
					fv, _ := rowV.ToNumber()
					frv := float64(fv) * 0.000001
					if err = col.SetValue(rowIdx, dataframe.NewStringValue(fmt.Sprintf("%.2f", frv))); err != nil {
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
