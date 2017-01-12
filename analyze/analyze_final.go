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

import "github.com/gyuho/dataframe"

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
			all.headerToLegend[hd] = elem.Legend
		}
	}

	plog.Println("combining data for plotting")
	for _, plotConfig := range cfg.PlotList {
		plog.Printf("plotting %q", plotConfig.Column)
		var columns []dataframe.Column
		for i, ad := range all.data {
			tag := all.databaseTags[i]
			col, err := ad.aggregated.Column(makeHeader(plotConfig.Column, tag))
			if err != nil {
				return err
			}
			columns = append(columns, col)
		}
		if err = all.draw(plotConfig, columns...); err != nil {
			return err
		}
	}

	plog.Printf("writing README at %q", cfg.READMEConfig.OutputPath)
	return writeREADME(cfg.READMEConfig)
}
