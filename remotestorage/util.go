// Copyright 2016 CoreOS, Inc.
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

package remotestorage

import (
	"os"
	"path/filepath"
	"strings"
)

func walkRecursive(dir string) (map[string]string, error) {
	fmap := make(map[string]string)
	visit := func(path string, f os.FileInfo, err error) error {
		if f != nil {
			if !f.IsDir() {
				if !filepath.HasPrefix(path, ".") && !strings.Contains(path, "/.") {
					wd, err := os.Getwd()
					if err != nil {
						return err
					}
					fmap[filepath.Join(wd, strings.Replace(path, wd, "", -1))] = path
				}
			}
		}
		return nil
	}
	if err := filepath.Walk(dir, visit); err != nil {
		return nil, err
	}
	return fmap, nil
}
