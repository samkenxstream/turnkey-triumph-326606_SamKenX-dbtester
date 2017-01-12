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

import "github.com/spf13/cobra"

// Command implements 'analyze' command.
var Command = &cobra.Command{
	Use:   "analyze",
	Short: "Analyzes test dbtester test results.",
	RunE:  commandFunc,
}

var configPath string

func init() {
	Command.PersistentFlags().StringVarP(&configPath, "config", "c", "", "YAML configuration file path.")
}

func commandFunc(cmd *cobra.Command, args []string) error {
	return do(configPath)
}
