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

// Package fileinspect implements various file inspecting utils.
package fileinspect

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

// Walk walks all files in the target directory.
func Walk(targetDir string) (map[string]os.FileInfo, error) {
	rm := make(map[string]os.FileInfo)
	visit := func(path string, f os.FileInfo, err error) error {
		if f != nil {
			if !f.IsDir() {
				wd, err := os.Getwd()
				if err != nil {
					return err
				}
				rm[filepath.Join(wd, strings.Replace(path, wd, "", -1))] = f
			}
		}
		return nil
	}
	err := filepath.Walk(targetDir, visit)
	if err != nil {
		return nil, err
	}
	return rm, nil
}

// Size returns the size of target directory, in bytes.
// Same as 'du -sh $DIR'.
func Size(targetDir string) (int64, error) {
	fm, err := Walk(targetDir)
	if err != nil {
		return 0, err
	}
	var size int64
	for _, v := range fm {
		size += v.Size()
	}
	return size, nil
}

func walk(targetDir string) (map[string]os.FileInfo, error) {
	rm := make(map[string]os.FileInfo)
	visit := func(path string, f os.FileInfo, err error) error {
		if f != nil {
			if !f.IsDir() {
				if !filepath.HasPrefix(path, ".") && !strings.Contains(path, "/.") {
					wd, err := os.Getwd()
					if err != nil {
						return err
					}
					rm[filepath.Join(wd, strings.Replace(path, wd, "", -1))] = f
				}
			}
		}
		return nil
	}
	err := filepath.Walk(targetDir, visit)
	if err != nil {
		return nil, err
	}
	return rm, nil
}

type filepathSize struct {
	path    string
	size    uint64
	sizeTxt string
}

func filterByKbs(fs []filepathSize, kbLimit int) []filepathSize {
	var ns []filepathSize
	for _, v := range fs {
		if v.size > uint64(kbLimit*1024) {
			continue
		}
		ns = append(ns, v)
	}
	return ns
}

type filepathSizeSlice []filepathSize

func (f filepathSizeSlice) Len() int           { return len(f) }
func (f filepathSizeSlice) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f filepathSizeSlice) Less(i, j int) bool { return f[i].size < f[j].size }

func walkDir(targetDir string) ([]filepathSize, error) {
	rm, err := walk(targetDir)
	if err != nil {
		return nil, err
	}

	var fs []filepathSize
	for k, v := range rm {
		fv := filepathSize{
			path:    k,
			size:    uint64(v.Size()),
			sizeTxt: humanize.Bytes(uint64(v.Size())),
		}
		fs = append(fs, fv)
	}
	sort.Sort(filepathSizeSlice(fs))

	return fs, nil
}
