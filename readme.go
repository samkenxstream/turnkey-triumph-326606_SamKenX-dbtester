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

import (
	"bytes"
	"fmt"
	"path/filepath"
)

// WriteREADME writes README.
func (cfg *Config) WriteREADME(summary string) error {
	plog.Printf("writing README at %q", cfg.ConfigAnalyzeMachineREADME.OutputPath)

	buf := new(bytes.Buffer)
	buf.WriteString("\n\n")
	buf.WriteString(fmt.Sprintf("<br><br><hr>\n##### %s", cfg.TestTitle))
	buf.WriteString("\n\n")
	buf.WriteString(cfg.TestDescription)
	buf.WriteString("\n\n```\n")
	buf.WriteString(summary)
	buf.WriteString("```\n\n\n")

	for _, img := range cfg.Images {
		switch img.Type {
		case "local":
			imgPath := "./" + filepath.Base(img.Path)
			buf.WriteString(fmt.Sprintf("![%s](%s)\n\n", img.Title, imgPath))
		case "remote":
			buf.WriteString(fmt.Sprintf(`<img src="%s" alt="%s">`, img.Path, img.Title))
			buf.WriteString("\n\n")
		default:
			return fmt.Errorf("%s is not supported", img.Type)
		}
		buf.WriteString("\n\n")
	}

	return toFile(buf.String(), cfg.ConfigAnalyzeMachineREADME.OutputPath)
}
