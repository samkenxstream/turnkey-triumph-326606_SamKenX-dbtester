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

// sysMetricsColumnsToRead is already aggregated
// and interpolated by unix second.
var sysMetricsColumnsToRead = []string{
	"UNIX-SECOND",
	"VOLUNTARY-CTXT-SWITCHES",
	"NON-VOLUNTARY-CTXT-SWITCHES",
	"CPU-NUM",
	"LOAD-AVERAGE-1-MINUTE",
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
	filePath        string
	frontUnixSecond int64
	lastUnixSecond  int64
	frame           dataframe.Frame
}

// readSystemMetrics extracts only the columns that we need for analyze.
func readSystemMetrics(fpath string) (data testData, err error) {
	originalFrame, err := dataframe.NewFromCSV(nil, fpath)
	if err != nil {
		return testData{}, err
	}

	data.filePath = fpath
	data.frame = dataframe.New()
	var unixSecondCol dataframe.Column
	for _, name := range sysMetricsColumnsToRead {
		var column dataframe.Column
		column, err = originalFrame.Column(name)
		if err != nil {
			return testData{}, err
		}
		if err = data.frame.AddColumn(column); err != nil {
			return testData{}, err
		}
		if name == "UNIX-SECOND" {
			unixSecondCol = column
		}
	}

	// get first(minimum) unix second
	fv, ok := unixSecondCol.FrontNonNil()
	if !ok {
		return testData{}, fmt.Errorf("FrontNonNil %s has empty Unix time %v", fpath, fv)
	}
	fs, ok := fv.String()
	if !ok {
		return testData{}, fmt.Errorf("cannot String %v", fv)
	}
	data.frontUnixSecond, err = strconv.ParseInt(fs, 10, 64)
	if err != nil {
		return testData{}, err
	}

	// get last(maximum) unix second
	bv, ok := unixSecondCol.BackNonNil()
	if !ok {
		return testData{}, fmt.Errorf("BackNonNil %s has empty Unix time %v", fpath, fv)
	}
	bs, ok := bv.String()
	if !ok {
		return testData{}, fmt.Errorf("cannot String %v", bv)
	}
	data.lastUnixSecond, err = strconv.ParseInt(bs, 10, 64)
	if err != nil {
		return testData{}, err
	}

	return
}
