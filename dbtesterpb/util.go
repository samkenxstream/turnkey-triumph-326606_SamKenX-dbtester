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

package dbtesterpb

import (
	"image/color"
	"sort"

	"gonum.org/v1/plot/plotutil"
)

// IsValidDatabaseID returns false if the database id is not supported.
func IsValidDatabaseID(id string) bool {
	_, ok := DatabaseID_value[id]
	return ok
}

// GetAllDatabaseIDs returns all database ids.
func GetAllDatabaseIDs() []string {
	var ids []string
	for k := range DatabaseID_value {
		ids = append(ids, k)
	}
	sort.Strings(ids)
	return ids
}

func GetRGBI(databaseID string, i int) color.Color {
	switch databaseID {
	case "etcd__other":
		return color.RGBA{94, 191, 30, 255} // green
	case "etcd__tip":
		return color.RGBA{24, 90, 169, 255} // blue
	case "etcd__v3_2":
		return color.RGBA{0, 229, 255, 255} // cyan
	case "etcd__v3_3":
		return color.RGBA{63, 81, 181, 255} // indigo
	case "zookeeper__r3_5_3_beta":
		return color.RGBA{94, 191, 30, 255} // green
	case "consul__v1_0_2":
		return color.RGBA{254, 25, 102, 255} // red
	case "zetcd__beta":
		return color.RGBA{251, 206, 0, 255} // yellow
	case "cetcd__beta":
		return color.RGBA{205, 220, 57, 255} // lime
	}
	return plotutil.Color(i)
}

func GetRGBII(databaseID string, i int) color.Color {
	switch databaseID {
	case "etcd__other":
		return color.RGBA{155, 176, 29, 255} // light-green
	case "etcd__tip":
		return color.RGBA{129, 212, 247, 255} // light-blue
	case "etcd__v3_2":
		return color.RGBA{132, 255, 255, 255} // light-cyan
	case "etcd__v3_3":
		return color.RGBA{159, 168, 218, 255} // light-indigo
	case "zookeeper__r3_5_3_beta":
		return color.RGBA{155, 176, 29, 255} // light-green
	case "consul__v1_0_2":
		return color.RGBA{255, 202, 178, 255} // light-red
	case "zetcd__beta":
		return color.RGBA{245, 247, 166, 255} // light-yellow
	case "cetcd__beta":
		return color.RGBA{238, 255, 65, 255} // light-lime
	}
	return plotutil.Color(i)
}

func GetRGBIII(databaseID string, i int) color.Color {
	switch databaseID {
	case "etcd__other":
		return color.RGBA{129, 210, 178, 255} // mid-cyan
	case "etcd__tip":
		return color.RGBA{37, 29, 191, 255} // deep-blue
	case "etcd__v3_2":
		return color.RGBA{0, 96, 100, 255} // deep-cyan
	case "etcd__v3_3":
		return color.RGBA{26, 35, 126, 255} // deep-indigo
	case "zookeeper__r3_5_3_beta":
		return color.RGBA{129, 210, 178, 255} // mid-cyan
	case "consul__v1_0_2":
		return color.RGBA{245, 144, 84, 255} // deep-red
	case "zetcd__beta":
		return color.RGBA{229, 255, 0, 255} // deep-yellow
	case "cetcd__beta":
		return color.RGBA{205, 220, 57, 255} // deep-lime
	}
	return plotutil.Color(i)
}
