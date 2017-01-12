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

// importBenchMetrics adds benchmark metrics from client-side
// and aggregates this to system metrics by unix timestamps.
func (data *analyzeData) importBenchMetrics(fpath string) (err error) {
	data.benchMetricsFilePath = fpath
	data.benchMetrics.frame, err = dataframe.NewFromCSV(nil, fpath)
	if err != nil {
		return
	}

	var unixTSColumn dataframe.Column
	unixTSColumn, err = data.benchMetrics.frame.Column("UNIX-TS")
	if err != nil {
		return err
	}

	// get first(minimum) unix second
	fv, ok := unixTSColumn.FrontNonNil()
	if !ok {
		return fmt.Errorf("FrontNonNil %s has empty Unix time %v", fpath, fv)
	}
	fs, ok := fv.String()
	if !ok {
		return fmt.Errorf("cannot String %v", fv)
	}
	data.benchMetrics.frontUnixTS, err = strconv.ParseInt(fs, 10, 64)
	if err != nil {
		return err
	}

	// get last(maximum) unix second
	bv, ok := unixTSColumn.BackNonNil()
	if !ok {
		return fmt.Errorf("BackNonNil %s has empty Unix time %v", fpath, fv)
	}
	bs, ok := bv.String()
	if !ok {
		return fmt.Errorf("cannot String %v", bv)
	}
	data.benchMetrics.lastUnixTS, err = strconv.ParseInt(bs, 10, 64)
	if err != nil {
		return err
	}

	return
}
