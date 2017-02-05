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

import "testing"

func Test_readConfig(t *testing.T) {
	c, err := readConfig("analyze_config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if c.Title != `Write 100K keys, 1000-client (etcd 100 TCP conns), 256-byte key, 1KB value` {
		t.Fatalf("unexpected Title %q", c.Title)
	}
	if c.WorkDir != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys" {
		t.Fatalf("unexpected WorkDir %q", c.WorkDir)
	}
	if c.AllAggregatedPath != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/aggregated.csv" {
		t.Fatalf("unexpected AllAggregatedPath %q", c.AllAggregatedPath)
	}
	if c.AllLatencyByKey != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/aggregated-data-latency-by-key-number.csv" {
		t.Fatalf("unexpected AllLatencyByKey %q", c.AllLatencyByKey)
	}
	if c.AllMemoryByKey != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/aggregated-data-memory-by-key-number.csv" {
		t.Fatalf("unexpected AllMemoryByKey %q", c.AllMemoryByKey)
	}

	if c.RawData[0].Legend != "etcd v3.1 (Go 1.7.4)" {
		t.Fatalf("unexpected c.RawData[0].Legend %q", c.RawData[0].Legend)
	}
	if makeTag(c.RawData[0].Legend) != "etcd-v3.1-go1.7.4" {
		t.Fatalf("unexpected makeTag(c.RawData[0].Legend) %q", makeTag(c.RawData[0].Legend))
	}
	if c.RawData[0].OutputPath != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/etcd-v3.1-go1.7.4-aggregated.csv" {
		t.Fatalf("unexpected c.RawData[0].OutputPath %q", c.RawData[0].OutputPath)
	}
	if c.RawData[0].DataInterpolatedSystemMetricsPaths[0] != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/etcd-v3.1-go1.7.4-1-system-metrics-interpolated.csv" {
		t.Fatalf("unexpected c.RawData[0].DataInterpolatedSystemMetricsPaths[0] %q", c.RawData[0].DataInterpolatedSystemMetricsPaths[0])
	}
	if c.RawData[0].DatasizeSummary != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/etcd-v3.1-go1.7.4-data-size-summary.csv" {
		t.Fatalf("unexpected c.RawData[0].DatasizeSummary %q", c.RawData[0].DatasizeSummary)
	}
	if c.RawData[0].DataBenchmarkLatencyPercentile != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/etcd-v3.1-go1.7.4-data-latency-distribution-percentile.csv" {
		t.Fatalf("unexpected c.RawData[0].DataBenchmarkLatencyPercentile %q", c.RawData[0].DataBenchmarkLatencyPercentile)
	}
	if c.RawData[0].DataBenchmarkLatencySummary != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/etcd-v3.1-go1.7.4-data-latency-distribution-summary.csv" {
		t.Fatalf("unexpected c.RawData[0].DataBenchmarkLatencySummary %q", c.RawData[0].DataBenchmarkLatencySummary)
	}
	if c.RawData[0].DataBenchmarkThroughput != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/etcd-v3.1-go1.7.4-data-latency-throughput-timeseries.csv" {
		t.Fatalf("unexpected c.RawData[0].DataBenchmarkThroughput %q", c.RawData[0].DataBenchmarkThroughput)
	}
	if c.RawData[0].DataBenchmarkLatencyByKey != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/etcd-v3.1-go1.7.4-data-latency-by-key-number.csv" {
		t.Fatalf("unexpected c.RawData[0].DataBenchmarkLatencyByKey %q", c.RawData[0].DataBenchmarkLatencyByKey)
	}
	if c.RawData[0].DataBenchmarkMemoryByKey != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/etcd-v3.1-go1.7.4-data-memory-by-key-number.csv" {
		t.Fatalf("unexpected c.RawData[0].DataBenchmarkMemoryByKey %q", c.RawData[0].DataBenchmarkMemoryByKey)
	}
	if c.RawData[0].ClientSystemMetricsInterpolated != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/etcd-v3.1-go1.7.4-client-system-metrics-interpolated.csv" {
		t.Fatalf("unexpected c.RawData[0].ClientSystemMetricsInterpolated %q", c.RawData[0].ClientSystemMetricsInterpolated)
	}

	if c.READMEConfig.OutputPath != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/README.md" {
		t.Fatalf("unexpected %s", c.READMEConfig.OutputPath)
	}
	if c.READMEConfig.Results[0].Images[0].ImageTitle != "2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-LATENCY-MS" {
		t.Fatalf("unexpected %s", c.READMEConfig.Results[0].Images[0].ImageTitle)
	}
	if c.READMEConfig.Results[0].Images[0].ImagePath != "https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-LATENCY-MS.svg" {
		t.Fatalf("unexpected %s", c.READMEConfig.Results[0].Images[0].ImagePath)
	}
	if c.READMEConfig.Results[0].Images[0].ImageType != "remote" {
		t.Fatalf("unexpected %s", c.READMEConfig.Results[0].Images[0].ImageType)
	}
}
