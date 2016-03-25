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

package analyze

import (
	"fmt"

	"github.com/gyuho/dataframe"
	"github.com/spf13/cobra"
)

type (
	Flags struct {
		RawMonitorFilePaths string
	}
)

var (
	Command = &cobra.Command{
		Use:   "analyze",
		Short: "Analyzes test results specific to dbtester.",
		RunE:  CommandFunc,
	}

	globalFlags = Flags{}
)

func init() {
	Command.PersistentFlags().StringSliceVarP(&globalFlags.RawMonitorFilePaths, "monitor-data-file-paths", "m", []string{}, "Monitor file paths to aggregate.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {

	return nil
}

func aggregate(fpaths ...string) (dataframe.Frame, error) {
	if len(fpaths) == 0 {
		return nil, fmt.Errorf("no file specified")
	}

	for _, fpath := range fpaths {
		fr, err := dataframe.NewFromCSV(nil, fpath)
		if err != nil {
			return nil, err
		}
	}

}
