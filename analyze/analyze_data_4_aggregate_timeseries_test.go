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
	"reflect"
	"testing"
)

func Test_processTimeSeries(t *testing.T) {
	var tslice []keyNumAndMemory
	for i := int64(0); i < 10; i++ {
		dp := keyNumAndMemory{
			keyNum:   50,
			memoryMB: float64(i + 1),
		}
		tslice = append(tslice, dp)
	}

	pss := processTimeSeries(tslice, 20, 555)
	expexcted := []keyNumAndMemory{
		{keyNum: 20, memoryMB: 1},
		{keyNum: 40, memoryMB: 1},
		{keyNum: 60, memoryMB: 1},
		{keyNum: 80, memoryMB: 2},
		{keyNum: 100, memoryMB: 2},
		{keyNum: 120, memoryMB: 3},
		{keyNum: 140, memoryMB: 3},
		{keyNum: 160, memoryMB: 3},
		{keyNum: 180, memoryMB: 4},
		{keyNum: 200, memoryMB: 4},
		{keyNum: 220, memoryMB: 5},
		{keyNum: 240, memoryMB: 5},
		{keyNum: 260, memoryMB: 5},
		{keyNum: 280, memoryMB: 6},
		{keyNum: 300, memoryMB: 6},
		{keyNum: 320, memoryMB: 7},
		{keyNum: 340, memoryMB: 7},
		{keyNum: 360, memoryMB: 7},
		{keyNum: 380, memoryMB: 8},
		{keyNum: 400, memoryMB: 8},
		{keyNum: 420, memoryMB: 9},
		{keyNum: 440, memoryMB: 9},
		{keyNum: 460, memoryMB: 9},
		{keyNum: 480, memoryMB: 10},
		{keyNum: 500, memoryMB: 10},
		{keyNum: 520, memoryMB: 0},
		{keyNum: 540, memoryMB: 0},
		{keyNum: 555, memoryMB: 0},
	}
	if len(pss) != len(expexcted) {
		t.Fatalf("expected %+v, got %+v", expexcted, pss)
	}
	for i, elem := range pss {
		if !reflect.DeepEqual(elem, expexcted[i]) {
			t.Fatalf("#%d: processed data point expected %+v, got %+v", i, expexcted[i], elem)
		}
	}
}
