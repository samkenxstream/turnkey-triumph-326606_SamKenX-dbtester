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

package dbtester

import "time"

func toMillisecond(d time.Duration) float64 {
	return d.Seconds() * 1000
}

func max(n1, n2 int64) int64 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func assignRequest(ranges []int64, total int64) (rs []int64) {
	reqEach := int(float64(total) / float64(len(ranges)))
	// truncate 10000th digits
	if reqEach > 10000 {
		reqEach = (reqEach / 10000) * 10000
	}
	// truncate 1000th digits
	if reqEach > 1000 {
		reqEach = (reqEach / 1000) * 1000
	}

	curSum := int64(0)
	rs = make([]int64, len(ranges))
	for i := range ranges {
		if i < len(ranges)-1 {
			rs[i] = int64(reqEach)
			curSum += int64(reqEach)
		} else {
			rs[i] = int64(total) - curSum
		}
	}
	return
}
