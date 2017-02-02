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
	"strings"
)

func makeHeader(column string, tag string) string {
	return fmt.Sprintf("%s-%s", column, tag)
}

func makeTag(legend string) string {
	legend = strings.ToLower(legend)
	legend = strings.Replace(legend, "go ", "go", -1)
	legend = strings.Replace(legend, "java ", "java", -1)
	legend = strings.Replace(legend, "(", "", -1)
	legend = strings.Replace(legend, ")", "", -1)
	return strings.Replace(legend, " ", "-", -1)
}
