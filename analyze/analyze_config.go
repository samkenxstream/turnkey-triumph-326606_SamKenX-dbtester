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
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// TODO: deprecate un-used fields

// Step1 defines how to aggregate data from each machine.
type Step1 []struct {
	DataPathList      []string `yaml:"data_path_list"`
	DataBenchmarkPath string   `yaml:"data_benchmark_path"`
	OutputPath        string   `yaml:"output_path"`
}

// Step2 defines how to aggregate the data of each aggregated from Step1.
type Step2 []struct {
	DataList []struct {
		Path string `yaml:"path"`
		Name string `yaml:"name"`
	} `yaml:"data_list"`
	OutputPath string `yaml:"output_path"`
}

// Step3 defines how to plot graphs.
type Step3 []struct {
	DataPath string       `yaml:"data_path"`
	PlotList []PlotConfig `yaml:"plot_list"`
}

// Config defines analyze configuration.
type Config struct {
	Titles []string `yaml:"titles"`

	Step1        Step1        `yaml:"step1"`
	Step2        Step2        `yaml:"step2"`
	Step3        Step3        `yaml:"step3"`
	READMEConfig READMEConfig `yaml:"readme"`
}

// ReadConfig reads analyze configuration.
func ReadConfig(fpath string) (Config, error) {
	bts, err := ioutil.ReadFile(fpath)
	if err != nil {
		return Config{}, err
	}
	rs := Config{}
	if err := yaml.Unmarshal(bts, &rs); err != nil {
		return Config{}, err
	}
	return rs, nil
}
