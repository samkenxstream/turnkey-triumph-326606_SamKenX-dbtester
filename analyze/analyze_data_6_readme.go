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
	"bytes"
	"fmt"
	"path/filepath"
)

// READMEConfig defines how to write README.
type READMEConfig struct {
	Preface    string `yaml:"preface"`
	OutputPath string `yaml:"output_path"`
	Results    []struct {
		Title  string
		Images []struct {
			ImageTitle string `yaml:"image_title"`
			ImagePath  string `yaml:"image_path"`
			ImageType  string `yaml:"image_type"`
		} `yaml:"images"`
	} `yaml:"results"`
}

func writeREADME(summary string, cfg READMEConfig) error {
	buf := new(bytes.Buffer)

	buf.WriteString("\n\n")
	buf.WriteString(cfg.Preface)
	buf.WriteString("\n\n\n")

	for _, result := range cfg.Results {
		buf.WriteString(fmt.Sprintf("<br><br>\n##### %s", result.Title))
		buf.WriteString("\n\n```\n")
		buf.WriteString(summary)
		buf.WriteString("```\n\n\n")
		for _, img := range result.Images {
			imgPath := ""
			switch img.ImageType {
			case "local":
				imgPath = "./" + filepath.Base(img.ImagePath)
				buf.WriteString(fmt.Sprintf("![%s](%s)\n\n", img.ImageTitle, imgPath))
			case "remote":
				buf.WriteString(fmt.Sprintf(`<img src="%s" alt="%s">`, img.ImagePath, img.ImageTitle))
				buf.WriteString("\n\n")
			default:
				return fmt.Errorf("%s is not supported", img.ImageType)
			}
		}
		buf.WriteString("\n\n")
	}

	return toFile(buf.String(), cfg.OutputPath)
}
