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
	"UNIX-TS",
	"VOLUNTARY-CTXT-SWITCHES",
	"NON-VOLUNTARY-CTXT-SWITCHES",
	"CPU-NUM",
	"VMRSS-NUM",
	"READS-COMPLETED",
	"READS-COMPLETED-DELTA",
	"SECTORS-READ",
	"SECTORS-READ-DELTA",
	"WRITES-COMPLETED",
	"WRITES-COMPLETED-DELTA",
	"SECTORS-WRITTEN",
	"SECTORS-WRITTEN-DELTA",
	"RECEIVE-BYTES-NUM",
	"RECEIVE-BYTES-NUM-DELTA",
	"TRANSMIT-BYTES-NUM",
	"TRANSMIT-BYTES-NUM-DELTA",
	"EXTRA", // will be converted to 'CLIENT-NUM'
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
		column, err = originalFrame.Column(name)
		if err != nil {
			return testData{}, err
		}
		if err = data.frame.AddColumn(column); err != nil {
			return testData{}, err
		}
		if name == "UNIX-TS" {
			// TODO: UNIX-TS from pkg/report data is time.Time.Unix
			// UNIX-TS from psn.CSV data is time.Time.UnixNano
			// we need some kind of way to combine those with matching timestamps
			//
			// this unixTSColumn is unix nanoseconds
			unixTSColumn = column
		}
	}

	// get first(minimum) unix second
	fv, ok := unixTSColumn.FrontNonNil()
	if !ok {
		return testData{}, fmt.Errorf("FrontNonNil %s has empty Unix time %v", fpath, fv)
	}
	fs, ok := fv.String()
	if !ok {
		return testData{}, fmt.Errorf("cannot String %v", fv)
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
	bs, ok := bv.String()
	if !ok {
		return testData{}, fmt.Errorf("cannot String %v", bv)
	}
	data.lastUnixTS, err = strconv.ParseInt(bs, 10, 64)
	if err != nil {
		return testData{}, err
	}

	return
}
