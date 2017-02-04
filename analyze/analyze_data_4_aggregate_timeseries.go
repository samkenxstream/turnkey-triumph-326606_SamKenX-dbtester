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

import "sort"

func processTimeSeries(tslice []keyNumAndMemory, unit int64, totalRequests int) []keyNumAndMemory {
	sort.Sort(keyNumAndMemorys(tslice))

	cumulKeyN := int64(0)
	maxKey := int64(0)

	rm := make(map[int64]keyNumAndMemory)

	// this data is aggregated by second
	// and we want to map number of keys to latency
	// so the range is the key
	// and the value is the cumulative throughput
	for _, ts := range tslice {
		cumulKeyN += ts.keyNum
		if cumulKeyN < unit {
			// not enough data points yet
			continue
		}

		mem := ts

		// cumulKeyN >= unit
		for cumulKeyN > maxKey {
			maxKey += unit
			rm[maxKey] = mem
		}
	}

	// fill-in empty rows
	for i := maxKey; i < int64(totalRequests); i += unit {
		if _, ok := rm[i]; !ok {
			rm[i] = keyNumAndMemory{}
		}
	}
	if _, ok := rm[int64(totalRequests)]; !ok {
		rm[int64(totalRequests)] = keyNumAndMemory{}
	}

	kss := []keyNumAndMemory{}
	delete(rm, 0)
	for k, v := range rm {
		kn := keyNumAndMemory{keyNum: k, maxMemoryMB: v.maxMemoryMB, avgMemoryMB: v.avgMemoryMB, minMemoryMB: v.minMemoryMB}
		kss = append(kss, kn)
	}
	sort.Sort(keyNumAndMemorys(kss))

	return kss
}

type keyNumAndMemory struct {
	keyNum int64

	maxMemoryMB float64
	avgMemoryMB float64
	minMemoryMB float64
}

type keyNumAndMemorys []keyNumAndMemory

func (t keyNumAndMemorys) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t keyNumAndMemorys) Len() int           { return len(t) }
func (t keyNumAndMemorys) Less(i, j int) bool { return t[i].keyNum < t[j].keyNum }
