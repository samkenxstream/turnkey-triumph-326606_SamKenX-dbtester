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

	"github.com/gonum/plot/plotutil"
)

// IsValidDatabaseID returns false if the database id is not supported.
func IsValidDatabaseID(id string) bool {
	_, ok := DatabaseID_value[id]
	return ok
}

func GetRGBI(databaseID string, i int) color.Color {
	switch databaseID {
	case "etcd__v2_3":
		return color.RGBA{218, 97, 229, 255} // purple
	case "etcd__v3_1":
		return color.RGBA{24, 90, 169, 255} // blue
	case "etcd__v3_2":
		return color.RGBA{63, 81, 181, 255} // indigo
	case "etcd__tip":
		return color.RGBA{0, 229, 255, 255} // cyan
	case "zookeeper__r3_5_2_alpha":
		return color.RGBA{38, 169, 24, 255} // green
	case "consul__v0_7_5":
		return color.RGBA{198, 53, 53, 255} // red
	case "zetcd__beta":
		return color.RGBA{251, 206, 0, 255} // yellow
	case "cetcd__beta":
		return color.RGBA{205, 220, 57, 255} // lime
	}
	return plotutil.Color(i)
}

func GetRGBII(databaseID string, i int) color.Color {
	switch databaseID {
	case "etcd__v2_3":
		return color.RGBA{229, 212, 231, 255} // light-purple
	case "etcd__v3_1":
		return color.RGBA{129, 212, 247, 255} // light-blue
	case "etcd__v3_2":
		return color.RGBA{159, 168, 218, 255} // light-indigo
	case "etcd__tip":
		return color.RGBA{132, 255, 255, 255} // light-cyan
	case "zookeeper__r3_5_2_alpha":
		return color.RGBA{129, 247, 152, 255} // light-green
	case "consul__v0_7_5":
		return color.RGBA{247, 156, 156, 255} // light-red
	case "zetcd__beta":
		return color.RGBA{245, 247, 166, 255} // light-yellow
	case "cetcd__beta":
		return color.RGBA{238, 255, 65, 255} // light-lime
	}
	return plotutil.Color(i)
}

func GetRGBIII(databaseID string, i int) color.Color {
	switch databaseID {
	case "etcd__v2_3":
		return color.RGBA{165, 8, 180, 255} // deep-purple
	case "etcd__v3_1":
		return color.RGBA{37, 29, 191, 255} // deep-blue
	case "etcd__v3_2":
		return color.RGBA{26, 35, 126, 255} // deep-indigo
	case "etcd__tip":
		return color.RGBA{0, 96, 100, 255} // deep-cyan
	case "zookeeper__r3_5_2_alpha":
		return color.RGBA{7, 64, 35, 255} // deep-green
	case "consul__v0_7_5":
		return color.RGBA{212, 8, 46, 255} // deep-red
	case "zetcd__beta":
		return color.RGBA{229, 255, 0, 255} // deep-yellow
	case "cetcd__beta":
		return color.RGBA{205, 220, 57, 255} // deep-lime
	}
	return plotutil.Color(i)
}
