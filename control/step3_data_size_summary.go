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

package control

import (
	"github.com/coreos/dbtester/agent/agentpb"
	humanize "github.com/dustin/go-humanize"
	"github.com/gyuho/dataframe"
)

var DataSummaryColumns = []string{
	"INDEX",
	"DATABASE-ENDPOINT",
	"TOTAL-DATA-SIZE",
	"TOTAL-DATA-SIZE-BYTES-NUM",
}

func saveDatasizeSummary(cfg Config, idxToResponse map[int]agentpb.Response) {
	c1 := dataframe.NewColumn(DataSummaryColumns[0])
	c2 := dataframe.NewColumn(DataSummaryColumns[1])
	c3 := dataframe.NewColumn(DataSummaryColumns[2])
	c4 := dataframe.NewColumn(DataSummaryColumns[3])
	for i := range cfg.DatabaseEndpoints {
		c1.PushBack(dataframe.NewStringValue(i))
		c2.PushBack(dataframe.NewStringValue(cfg.DatabaseEndpoints[i]))
		c3.PushBack(dataframe.NewStringValue(humanize.Bytes(uint64(idxToResponse[i].Datasize))))
		c4.PushBack(dataframe.NewStringValue(idxToResponse[i].Datasize))
	}
	fr := dataframe.New()
	if err := fr.AddColumn(c1); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(c2); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(c3); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(c4); err != nil {
		plog.Fatal(err)
	}
	if err := fr.CSV(cfg.DatasizeSummary); err != nil {
		plog.Fatal(err)
	}
}
