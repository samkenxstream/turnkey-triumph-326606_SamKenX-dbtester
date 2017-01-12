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

// combineAnalyzeData combines multiple analyzeData into one by given header.
// So, slice should consist of aggregated data of etcd, Zookeeper, Consul, etc..
// This combines them into one data frame.
func combineAnalyzeData(header string, ds ...*analyzeData) (dataframe.Frame, error) {
	minEndIndex := 0
	columns := make([]dataframe.Column, len(ds))
	for i, ad := range ds {
		col, err := ad.allDataFrame.Column(header)
		if err != nil {
			return nil, err
		}

		// since we have same headers from different databases
		col.UpdateHeader(fmt.Sprintf("%s-%s", header, ad.database))

		columns[i] = col

		if i == 0 {
			minEndIndex = col.CountRow()
		}
		if minEndIndex > col.CountRow() {
			minEndIndex = col.CountRow()
		}
	}
	// this is index, so decrement by 1 to make it as valid index
	minEndIndex--
	maxSize := minEndIndex + 1

	// make all columns have same row number
	for _, col := range columns {
		rNum := col.CountRow()
		if rNum < maxSize { // fill-in with zero values
			for i := 0; i < maxSize-rNum; i++ {
				col.PushBack(dataframe.NewStringValue(0))
			}
		}
		if rNum > maxSize {
			return nil, fmt.Errorf("something wrong with minimum end index %d (%q has %d rows)", minEndIndex, col.Header(), rNum)
		}
	}
	rNum := columns[0].CountRow()
	for _, col := range columns {
		if rNum != col.CountRow() {
			return nil, fmt.Errorf("%q has %d rows (expected %d rows as %q)", col.Header(), col.CountRow(), rNum, columns[0].Header())
		}
	}

	combined := dataframe.New()
	for _, col := range columns {
		if err := combined.AddColumn(col); err != nil {
			return nil, err
		}
	}
	return combined, nil
}
