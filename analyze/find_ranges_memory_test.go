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

func Test_findRangesMemory(t *testing.T) {
	var tslice []keyNumAndMemory
	for i := int64(0); i < 10; i++ {
		dp := keyNumAndMemory{
			keyNum:      50,
			avgMemoryMB: float64(i + 1),
		}
		tslice = append(tslice, dp)
	}

	pss := findRangesMemory(tslice, 20, 555)
	expexcted := []keyNumAndMemory{
		{keyNum: 20, avgMemoryMB: 1},
		{keyNum: 40, avgMemoryMB: 1},
		{keyNum: 60, avgMemoryMB: 1},
		{keyNum: 80, avgMemoryMB: 2},
		{keyNum: 100, avgMemoryMB: 2},
		{keyNum: 120, avgMemoryMB: 3},
		{keyNum: 140, avgMemoryMB: 3},
		{keyNum: 160, avgMemoryMB: 3},
		{keyNum: 180, avgMemoryMB: 4},
		{keyNum: 200, avgMemoryMB: 4},
		{keyNum: 220, avgMemoryMB: 5},
		{keyNum: 240, avgMemoryMB: 5},
		{keyNum: 260, avgMemoryMB: 5},
		{keyNum: 280, avgMemoryMB: 6},
		{keyNum: 300, avgMemoryMB: 6},
		{keyNum: 320, avgMemoryMB: 7},
		{keyNum: 340, avgMemoryMB: 7},
		{keyNum: 360, avgMemoryMB: 7},
		{keyNum: 380, avgMemoryMB: 8},
		{keyNum: 400, avgMemoryMB: 8},
		{keyNum: 420, avgMemoryMB: 9},
		{keyNum: 440, avgMemoryMB: 9},
		{keyNum: 460, avgMemoryMB: 9},
		{keyNum: 480, avgMemoryMB: 10},
		{keyNum: 500, avgMemoryMB: 10},
		{keyNum: 520, avgMemoryMB: 0},
		{keyNum: 540, avgMemoryMB: 0},
		{keyNum: 555, avgMemoryMB: 0},
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
