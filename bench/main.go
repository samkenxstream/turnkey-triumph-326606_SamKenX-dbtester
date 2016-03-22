// Copyright 2015 CoreOS, Inc.
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

// bench does benchmark.
package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/spf13/cobra"
)

// This represents the base command when called without any subcommands
var Command = &cobra.Command{
	Use:   "bench",
	Short: "Low-level benchmark tool for etcd, Zookeeper, etcd2, consul.",
}

var (
	database     string
	endpoints    []string
	totalConns   uint
	totalClients uint
	sample       bool
	noHistogram  bool

	csvResultPath string

	bar     *pb.ProgressBar
	results chan result
	wg      sync.WaitGroup
)

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	Command.PersistentFlags().StringVarP(&database, "database", "d", "etcd", "'etcd', 'zk'(zookeeper)")
	Command.PersistentFlags().StringSliceVar(&endpoints, "endpoints", []string{"10.240.0.9:2181", "10.240.0.10:2181", "10.240.0.14:2181"}, "gRPC endpoints")
	Command.PersistentFlags().UintVar(&totalConns, "conns", 1, "Total number of gRPC connections or Zookeeper connections")
	Command.PersistentFlags().UintVar(&totalClients, "clients", 1, "Total number of gRPC clients (only for etcd)")
	Command.PersistentFlags().BoolVar(&sample, "sample", false, "'true' to sample requests for every second.")
	Command.PersistentFlags().BoolVar(&noHistogram, "no-histogram", false, "'true' to not show results in histogram.")

	Command.PersistentFlags().StringVar(&csvResultPath, "csv-result-path", "timeseries.csv", "path to store csv results.")
}

func main() {
	log.Printf("bench started at %s\n", time.Now().String()[:19])
	if err := Command.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	log.Printf("bench ended at %s\n", time.Now().String()[:19])
}
