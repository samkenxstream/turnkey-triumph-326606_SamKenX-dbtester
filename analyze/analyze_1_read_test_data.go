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
	"strconv"

	"github.com/gyuho/dataframe"
)

var sysMetricsColumnsToRead = []string{
	"UNIX-TS", "CPU-NUM", "VMRSS-NUM",
	"READS-COMPLETED",
	"READS-COMPLETED-DIFF",
	"SECTORS-READ",
	"SECTORS-READ-DIFF",
	"WRITES-COMPLETED",
	"WRITES-COMPLETED-DIFF",
	"SECTORS-WRITTEN",
	"SECTORS-WRITTEN-DIFF",
	"RECEIVE-BYTES-NUM",
	"RECEIVE-BYTES-NUM-DIFF",
	"TRANSMIT-BYTES-NUM",
	"TRANSMIT-BYTES-NUM-DIFF",
	"EXTRA",
}

type testData struct {
	filePath    string
	frontUnixTS int64
	lastUnixTS  int64
	frame       dataframe.Frame
}

// readSystemMetrics extracts only the columns that we need for analyze.
func readSystemMetrics(fpath string) (data testData, err error) {
	originalFrame, err := dataframe.NewFromCSV(nil, fpath)
	if err != nil {
		return testData{}, err
	}

	data.filePath = fpath
	data.frame = dataframe.New()
	var unixTSColumn dataframe.Column
	for _, name := range sysMetricsColumnsToRead {
		var column dataframe.Column
		column, err = originalFrame.GetColumn(name)
		if err != nil {
			return testData{}, err
		}
		if err = data.frame.AddColumn(column); err != nil {
			return testData{}, err
		}
		if name == "UNIX-TS" {
			unixTSColumn = column
		}
	}

	// get first(minimum) unix second
	fv, ok := unixTSColumn.FrontNonNil()
	if !ok {
		return testData{}, fmt.Errorf("FrontNonNil %s has empty Unix time %v", fpath, fv)
	}
	fs, ok := fv.ToString()
	if !ok {
		return testData{}, fmt.Errorf("cannot ToString %v", fv)
	}
	data.frontUnixTS, err = strconv.ParseInt(fs, 10, 64)
	if err != nil {
		return testData{}, err
	}

	// get last(maximum) unix second
	bv, ok := unixTSColumn.BackNonNil()
	if !ok {
		return testData{}, fmt.Errorf("BackNonNil %s has empty Unix time %v", fpath, fv)
	}
	bs, ok := bv.ToString()
	if !ok {
		return testData{}, fmt.Errorf("cannot ToString %v", bv)
	}
	data.lastUnixTS, err = strconv.ParseInt(bs, 10, 64)
	if err != nil {
		return testData{}, err
	}

	return
}
