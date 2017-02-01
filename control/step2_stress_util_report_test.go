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
	"testing"

	"time"

	"reflect"

	"github.com/coreos/etcd/pkg/report"
)

func Test_processDataPoints(t *testing.T) {
	var tslice report.TimeSeries
	for i := int64(0); i < 10; i++ {
		dp := report.DataPoint{
			Timestamp:  i,
			AvgLatency: time.Duration(i + 1),
			ThroughPut: 50,
		}
		tslice = append(tslice, dp)
	}

	pss := processDataPoints(tslice, 20)
	expexcted := []keyNumToAvgLatency{
		{keyNum: 20, avgLat: 1},
		{keyNum: 40, avgLat: 1},
		{keyNum: 60, avgLat: 1},
		{keyNum: 80, avgLat: 2},
		{keyNum: 100, avgLat: 2},
		{keyNum: 120, avgLat: 3},
		{keyNum: 140, avgLat: 3},
		{keyNum: 160, avgLat: 3},
		{keyNum: 180, avgLat: 4},
		{keyNum: 200, avgLat: 4},
		{keyNum: 220, avgLat: 5},
		{keyNum: 240, avgLat: 5},
		{keyNum: 260, avgLat: 5},
		{keyNum: 280, avgLat: 6},
		{keyNum: 300, avgLat: 6},
		{keyNum: 320, avgLat: 7},
		{keyNum: 340, avgLat: 7},
		{keyNum: 360, avgLat: 7},
		{keyNum: 380, avgLat: 8},
		{keyNum: 400, avgLat: 8},
		{keyNum: 420, avgLat: 9},
		{keyNum: 440, avgLat: 9},
		{keyNum: 460, avgLat: 9},
		{keyNum: 480, avgLat: 10},
		{keyNum: 500, avgLat: 10},
	}
	for i, elem := range pss {
		if !reflect.DeepEqual(elem, expexcted[i]) {
			t.Fatalf("#%d: processed data point expected %+v, got %+v", i, expexcted[i], elem)
		}
	}
}
