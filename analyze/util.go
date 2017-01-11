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
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gonum/plot/plotutil"
)

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("copy: mkdirall: %v", err)
	}

	r, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("copy: open(%q): %v", src, err)
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("copy: create(%q): %v", dst, err)
	}
	defer w.Close()

	// func Copy(dst Writer, src Reader) (written int64, err error)
	if _, err = io.Copy(w, r); err != nil {
		return err
	}
	if err := w.Sync(); err != nil {
		return err
	}
	if _, err := w.Seek(0, 0); err != nil {
		return err
	}
	return nil
}

func getRGB(legend string, i int) color.Color {
	legend = strings.ToLower(strings.TrimSpace(legend))
	if strings.HasPrefix(legend, "etcd") {
		return color.RGBA{24, 90, 169, 255} // blue
	}
	if strings.HasPrefix(legend, "zookeeper") {
		return color.RGBA{38, 169, 24, 255} // green
	}
	if strings.HasPrefix(legend, "consul") {
		return color.RGBA{198, 53, 53, 255} // red
	}
	if strings.HasPrefix(legend, "zetcd") {
		return color.RGBA{251, 206, 0, 255} // yellow
	}
	if strings.HasPrefix(legend, "cetcd") {
		return color.RGBA{116, 24, 169, 255} // purple
	}
	return plotutil.Color(i)
}
