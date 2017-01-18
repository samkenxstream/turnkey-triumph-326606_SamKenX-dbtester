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
	"reflect"
	"testing"
)

func Test_assignRequest(t *testing.T) {
	ranges := []int{1, 10, 50, 100, 300, 500, 700, 1000, 1500, 2000, 2500}
	total := 3000000
	rs := assignRequest(ranges, total)
	expected := []int{250000, 250000, 250000, 250000, 250000, 250000, 250000, 250000, 250000, 250000, 500000}
	if !reflect.DeepEqual(rs, expected) {
		t.Fatalf("expected %+v, got %+v", expected, rs)
	}
}
