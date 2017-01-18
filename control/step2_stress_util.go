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

func assignRequest(ranges []int, total int) (rs []int) {
	reqEach := total / (len(ranges) + 1)
	reqEach = (reqEach / 100000) * 100000
	curSum := 0
	rs = make([]int, len(ranges))
	for i := range ranges {
		if i < len(ranges)-1 {
			rs[i] = reqEach
			curSum += reqEach
		} else {
			rs[i] = total - curSum
		}
	}
	return
}
