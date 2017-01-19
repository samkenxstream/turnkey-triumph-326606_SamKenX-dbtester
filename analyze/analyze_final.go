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
	"path/filepath"
	"sort"

	"github.com/gyuho/dataframe"
)

type allAggregatedData struct {
	title          string
	data           []*analyzeData
	databaseTags   []string
	headerToLegend map[string]string
}

func do(configPath string) error {
	cfg, err := readConfig(configPath)
	if err != nil {
		return err
	}

	all := &allAggregatedData{
		title:          cfg.Title,
		data:           make([]*analyzeData, 0, len(cfg.RawData)),
		headerToLegend: make(map[string]string),
	}
	for _, elem := range cfg.RawData {
		plog.Printf("reading system metrics data for %s (%q)", elem.DatabaseTag, elem.Legend)
		ad, err := readSystemMetricsAll(elem.SourceSystemMetricsPaths...)
		if err != nil {
			return err
		}
		ad.databaseTag = elem.DatabaseTag
		ad.legend = elem.Legend
		ad.csvOutputpath = elem.OutputPath

		if err = ad.aggSystemMetrics(); err != nil {
			return err
		}
		if err = ad.importBenchMetrics(elem.SourceBenchmarkMetricsPath); err != nil {
			return err
		}
		if err = ad.aggregateAll(); err != nil {
			return err
		}
		if err = ad.save(); err != nil {
			return err
		}

		all.data = append(all.data, ad)
		all.databaseTags = append(all.databaseTags, elem.DatabaseTag)
		for _, hd := range ad.aggregated.Headers() {
			all.headerToLegend[makeHeader(hd, elem.DatabaseTag)] = elem.Legend
		}
	}

	plog.Println("combining data for plotting")
	for _, plotConfig := range cfg.PlotList {
		plog.Printf("plotting %q", plotConfig.Column)
		var clientNumColumns []dataframe.Column
		var dataColumns []dataframe.Column
		for i, ad := range all.data {
			tag := all.databaseTags[i]

			avgCol, err := ad.aggregated.Column("CONTROL-CLIENT-NUM")
			if err != nil {
				return err
			}
			avgCol.UpdateHeader(makeHeader("CONTROL-CLIENT-NUM", tag))
			clientNumColumns = append(clientNumColumns, avgCol)

			col, err := ad.aggregated.Column(plotConfig.Column)
			if err != nil {
				return err
			}
			col.UpdateHeader(makeHeader(plotConfig.Column, tag))
			dataColumns = append(dataColumns, col)
		}
		if err = all.draw(plotConfig, dataColumns...); err != nil {
			return err
		}

		plog.Printf("saving data for %q of all database", plotConfig.Column)
		nf1, err := dataframe.NewFromColumns(nil, dataColumns...)
		if err != nil {
			return err
		}
		if err = nf1.CSV(filepath.Join(cfg.WorkDir, plotConfig.Column+".csv")); err != nil {
			return err
		}

		plog.Printf("saving data for %q of all database (by client number)", plotConfig.Column)
		nf2 := dataframe.New()
		for i := range clientNumColumns {
			if clientNumColumns[i].Count() != dataColumns[i].Count() {
				return fmt.Errorf("clientNumColumns[i].Count() %d != dataColumns[i].Count() %d", clientNumColumns[i].Count(), dataColumns[i].Count())
			}
			if err := nf2.AddColumn(clientNumColumns[i]); err != nil {
				return err
			}
			if err := nf2.AddColumn(dataColumns[i]); err != nil {
				return err
			}
		}
		if err = nf2.CSV(filepath.Join(cfg.WorkDir, plotConfig.Column+"-BY-CLIENT-NUM"+".csv")); err != nil {
			return err
		}

		plog.Printf("aggregating data for %q of all database (by client number)", plotConfig.Column)
		nf3 := dataframe.New()
		for i := range clientNumColumns {
			n := clientNumColumns[i].Count()
			allData := make(map[int]float64)
			for j := 0; j < n; j++ {
				v1, err := clientNumColumns[i].Value(j)
				if err != nil {
					return err
				}
				num, _ := v1.Number()

				v2, err := dataColumns[i].Value(j)
				if err != nil {
					return err
				}
				data, _ := v2.Number()

				if v, ok := allData[int(num)]; ok {
					allData[int(num)] = (v + data) / 2
				} else {
					allData[int(num)] = data
				}
			}
			var allKeys []int
			for k := range allData {
				allKeys = append(allKeys, k)
			}
			sort.Ints(allKeys)

			col1 := dataframe.NewColumn(clientNumColumns[i].Header())
			col2 := dataframe.NewColumn(dataColumns[i].Header())
			for j := range allKeys {
				col1.PushBack(dataframe.NewStringValue(allKeys[j]))
				col2.PushBack(dataframe.NewStringValue(fmt.Sprintf("%.4f", allData[allKeys[j]])))
			}
			if err := nf3.AddColumn(col1); err != nil {
				return err
			}
			if err := nf3.AddColumn(col2); err != nil {
				return err
			}
		}
		if err = nf3.CSV(filepath.Join(cfg.WorkDir, plotConfig.Column+"-BY-CLIENT-NUM-aggregated"+".csv")); err != nil {
			return err
		}
	}

	plog.Printf("writing README at %q", cfg.READMEConfig.OutputPath)
	return writeREADME(cfg.READMEConfig)
}
