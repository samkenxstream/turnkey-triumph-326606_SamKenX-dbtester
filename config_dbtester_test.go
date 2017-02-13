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
	"reflect"
	"testing"

	"github.com/coreos/dbtester/dbtesterpb"
)

func TestConfig(t *testing.T) {
	cfg, err := ReadConfig("config_dbtester_test.yaml", false)
	if err != nil {
		t.Fatal(err)
	}
	expected := &Config{
		TestTitle: "Write 1M keys, 256-byte key, 1KB value value, clients 1 to 1,000",
		TestDescription: `- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 30 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.5)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.4 (Go 1.7.5)
`,
		Control: Control{
			PathPrefix:                              "/home/gyuho",
			LogPath:                                 "/home/gyuho/client-control.log",
			ClientSystemMetricsPath:                 "/home/gyuho/client-system-metrics.csv",
			ClientSystemMetricsInterpolatedPath:     "/home/gyuho/client-system-metrics-interpolated.csv",
			ClientLatencyThroughputTimeseriesPath:   "/home/gyuho/client-latency-throughput-timeseries.csv",
			ClientLatencyDistributionAllPath:        "/home/gyuho/client-latency-distribution-all.csv",
			ClientLatencyDistributionPercentilePath: "/home/gyuho/client-latency-distribution-percentile.csv",
			ClientLatencyDistributionSummaryPath:    "/home/gyuho/client-latency-distribution-summary.csv",
			ClientLatencyByKeyNumberPath:            "/home/gyuho/client-latency-by-key-number.csv",
			ServerDiskSpaceUsageSummaryPath:         "/home/gyuho/server-disk-space-usage-summary.csv",
			GoogleCloudProjectName:                  "etcd-development",
			GoogleCloudStorageKeyPath:               "config-dbtester-gcloud-key.json",
			GoogleCloudStorageKey:                   "test-key",
			GoogleCloudStorageBucketName:            "dbtester-results",
			GoogleCloudStorageSubDirectory:          "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable",
		},
		AllDatabaseIDList: []string{"etcdv3", "zookeeper", "consul"},
		DatabaseIDToTestGroup: map[string]TestGroup{
			"etcdv3": {
				DatabaseID:            "etcdv3",
				DatabaseTag:           "etcd-v3.1-go1.7.5",
				DatabaseDescription:   "etcd v3.1 (Go 1.7.5)",
				PeerIPs:               []string{"10.240.0.20", "10.240.0.21", "10.240.0.22"},
				PeerIPsString:         "10.240.0.20___10.240.0.21___10.240.0.22",
				DatabasePortToConnect: 2379,
				DatabaseEndpoints:     []string{"10.240.0.20:2379", "10.240.0.21:2379", "10.240.0.22:2379"},
				AgentPortToConnect:    3500,
				AgentEndpoints:        []string{"10.240.0.20:3500", "10.240.0.21:3500", "10.240.0.22:3500"},
				Etcdv3: Etcdv3{
					SnapCount:      100000,
					QuotaSizeBytes: 8000000000,
				},
				BenchmarkOptions: BenchmarkOptions{
					Type:                       "write",
					RequestNumber:              1000000,
					ConnectionNumber:           0,
					ClientNumber:               0,
					ConnectionClientNumbers:    []int64{1, 10, 50, 100, 300, 500, 700, 1000},
					RateLimitRequestsPerSecond: 0,
					SameKey:                    false,
					KeySizeBytes:               256,
					ValueSizeBytes:             1024,
					StaleRead:                  false,
				},
				BenchmarkSteps: BenchmarkSteps{
					Step1StartDatabase:  true,
					Step2StressDatabase: true,
					Step3StopDatabase:   true,
					Step4UploadLogs:     true,
				},
			},
			"zookeeper": {
				DatabaseID:            "zookeeper",
				DatabaseTag:           "zookeeper-r3.4.9-java8",
				DatabaseDescription:   "Zookeeper r3.4.9 (Java 8)",
				PeerIPs:               []string{"10.240.0.25", "10.240.0.27", "10.240.0.28"},
				PeerIPsString:         "10.240.0.25___10.240.0.27___10.240.0.28",
				DatabasePortToConnect: 2181,
				DatabaseEndpoints:     []string{"10.240.0.25:2181", "10.240.0.27:2181", "10.240.0.28:2181"},
				AgentPortToConnect:    3500,
				AgentEndpoints:        []string{"10.240.0.25:3500", "10.240.0.27:3500", "10.240.0.28:3500"},
				Zookeeper: Zookeeper{
					TickTime:             2000,
					InitLimit:            5,
					SyncLimit:            5,
					SnapCount:            100000,
					MaxClientConnections: 5000,
				},
				BenchmarkOptions: BenchmarkOptions{
					Type:                       "write",
					RequestNumber:              1000000,
					ConnectionNumber:           0,
					ClientNumber:               0,
					ConnectionClientNumbers:    []int64{1, 10, 50, 100, 300, 500, 700, 1000},
					RateLimitRequestsPerSecond: 0,
					SameKey:                    false,
					KeySizeBytes:               256,
					ValueSizeBytes:             1024,
					StaleRead:                  false,
				},
				BenchmarkSteps: BenchmarkSteps{
					Step1StartDatabase:  true,
					Step2StressDatabase: true,
					Step3StopDatabase:   true,
					Step4UploadLogs:     true,
				},
			},
			"consul": {
				DatabaseID:            "consul",
				DatabaseTag:           "consul-v0.7.4-go1.7.5",
				DatabaseDescription:   "Consul v0.7.4 (Go 1.7.5)",
				PeerIPs:               []string{"10.240.0.30", "10.240.0.31", "10.240.0.33"},
				PeerIPsString:         "10.240.0.30___10.240.0.31___10.240.0.33",
				DatabasePortToConnect: 8500,
				DatabaseEndpoints:     []string{"10.240.0.30:8500", "10.240.0.31:8500", "10.240.0.33:8500"},
				AgentPortToConnect:    3500,
				AgentEndpoints:        []string{"10.240.0.30:3500", "10.240.0.31:3500", "10.240.0.33:3500"},
				BenchmarkOptions: BenchmarkOptions{
					Type:                       "write",
					RequestNumber:              1000000,
					ConnectionNumber:           0,
					ClientNumber:               0,
					ConnectionClientNumbers:    []int64{1, 10, 50, 100, 300, 500, 700, 1000},
					RateLimitRequestsPerSecond: 0,
					SameKey:                    false,
					KeySizeBytes:               256,
					ValueSizeBytes:             1024,
					StaleRead:                  false,
				},
				BenchmarkSteps: BenchmarkSteps{
					Step1StartDatabase:  true,
					Step2StressDatabase: true,
					Step3StopDatabase:   true,
					Step4UploadLogs:     true,
				},
			},
		},
		DatabaseIDToTestData: map[string]TestData{
			"etcdv3": {
				DatabaseID:          "etcdv3",
				DatabaseTag:         "etcd-v3.1-go1.7.5",
				DatabaseDescription: "etcd v3.1 (Go 1.7.5)",
				PathPrefix:          "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5",

				ClientSystemMetricsInterpolatedPath:     "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-client-system-metrics-interpolated.csv",
				ClientLatencyThroughputTimeseriesPath:   "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-client-latency-throughput-timeseries.csv",
				ClientLatencyDistributionAllPath:        "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-client-latency-distribution-all.csv",
				ClientLatencyDistributionPercentilePath: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-client-latency-distribution-percentile.csv",
				ClientLatencyDistributionSummaryPath:    "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-client-latency-distribution-summary.csv",
				ClientLatencyByKeyNumberPath:            "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-client-latency-by-key-number.csv",
				ServerMemoryByKeyNumberPath:             "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-server-memory-by-key-number.csv",
				ServerReadBytesDeltaByKeyNumberPath:     "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-server-read-bytes-delta-by-key-number.csv",
				ServerWriteBytesDeltaByKeyNumberPath:    "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-server-write-bytes-delta-by-key-number.csv",
				ServerDiskSpaceUsageSummaryPath:         "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-server-disk-space-usage-summary.csv",
				ServerSystemMetricsInterpolatedPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-1-server-system-metrics-interpolated.csv",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-2-server-system-metrics-interpolated.csv",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-3-server-system-metrics-interpolated.csv",
				},
				AllAggregatedOutputPath: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/etcd-v3.1-go1.7.5-all-aggregated.csv",
			},
			"zookeeper": {
				DatabaseID:          "zookeeper",
				DatabaseTag:         "zookeeper-r3.4.9-java8",
				DatabaseDescription: "Zookeeper r3.4.9 (Java 8)",
				PathPrefix:          "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8",

				ClientSystemMetricsInterpolatedPath:     "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-client-system-metrics-interpolated.csv",
				ClientLatencyThroughputTimeseriesPath:   "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-client-latency-throughput-timeseries.csv",
				ClientLatencyDistributionAllPath:        "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-client-latency-distribution-all.csv",
				ClientLatencyDistributionPercentilePath: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-client-latency-distribution-percentile.csv",
				ClientLatencyDistributionSummaryPath:    "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-client-latency-distribution-summary.csv",
				ClientLatencyByKeyNumberPath:            "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-client-latency-by-key-number.csv",
				ServerMemoryByKeyNumberPath:             "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-server-memory-by-key-number.csv",
				ServerReadBytesDeltaByKeyNumberPath:     "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-server-read-bytes-delta-by-key-number.csv",
				ServerWriteBytesDeltaByKeyNumberPath:    "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-server-write-bytes-delta-by-key-number.csv",
				ServerDiskSpaceUsageSummaryPath:         "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-server-disk-space-usage-summary.csv",
				ServerSystemMetricsInterpolatedPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-1-server-system-metrics-interpolated.csv",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-2-server-system-metrics-interpolated.csv",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-3-server-system-metrics-interpolated.csv",
				},
				AllAggregatedOutputPath: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/zookeeper-r3.4.9-java8-all-aggregated.csv",
			},
			"consul": {
				DatabaseID:          "consul",
				DatabaseTag:         "consul-v0.7.4-go1.7.5",
				DatabaseDescription: "Consul v0.7.4 (Go 1.7.5)",
				PathPrefix:          "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5",

				ClientSystemMetricsInterpolatedPath:     "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-client-system-metrics-interpolated.csv",
				ClientLatencyThroughputTimeseriesPath:   "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-client-latency-throughput-timeseries.csv",
				ClientLatencyDistributionAllPath:        "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-client-latency-distribution-all.csv",
				ClientLatencyDistributionPercentilePath: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-client-latency-distribution-percentile.csv",
				ClientLatencyDistributionSummaryPath:    "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-client-latency-distribution-summary.csv",
				ClientLatencyByKeyNumberPath:            "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-client-latency-by-key-number.csv",
				ServerMemoryByKeyNumberPath:             "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-server-memory-by-key-number.csv",
				ServerReadBytesDeltaByKeyNumberPath:     "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-server-read-bytes-delta-by-key-number.csv",
				ServerWriteBytesDeltaByKeyNumberPath:    "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-server-write-bytes-delta-by-key-number.csv",
				ServerDiskSpaceUsageSummaryPath:         "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-server-disk-space-usage-summary.csv",
				ServerSystemMetricsInterpolatedPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-1-server-system-metrics-interpolated.csv",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-2-server-system-metrics-interpolated.csv",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-3-server-system-metrics-interpolated.csv",
				},
				AllAggregatedOutputPath: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/consul-v0.7.4-go1.7.5-all-aggregated.csv",
			},
		},
		Analyze: Analyze{
			AllAggregatedOutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/all-aggregated.csv",
			AllAggregatedOutputPathTXT: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/all-aggregated.txt",
		},
		PlotPathPrefix: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable",
		PlotList: []Plot{
			{
				Column:        "AVG-LATENCY-MS",
				XAxis:         "Second",
				YAxis:         "Latency(millisecond)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS.png",
				},
			},
			{
				Column:        "AVG-THROUGHPUT",
				XAxis:         "Second",
				YAxis:         "Throughput(Requests/Second)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT.png",
				},
			},
			{
				Column:        "AVG-VOLUNTARY-CTXT-SWITCHES",
				XAxis:         "Second",
				YAxis:         "Voluntary Context Switches",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.png",
				},
			},
			{
				Column:        "AVG-NON-VOLUNTARY-CTXT-SWITCHES",
				XAxis:         "Second",
				YAxis:         "Non-voluntary Context Switches",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.png",
				},
			},
			{
				Column:        "AVG-CPU",
				XAxis:         "Second",
				YAxis:         "CPU(%)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU.png",
				},
			},
			{
				Column:        "AVG-VMRSS-MB",
				XAxis:         "Second",
				YAxis:         "Memory(MB)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB.png",
				},
			},
			{
				Column:        "AVG-READS-COMPLETED-DELTA",
				XAxis:         "Second",
				YAxis:         "Disk Reads (Delta per Second)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.png",
				},
			},
			{
				Column:        "AVG-SECTORS-READ-DELTA",
				XAxis:         "Second",
				YAxis:         "Sectors Read (Delta per Second)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.png",
				},
			},
			{
				Column:        "AVG-WRITES-COMPLETED-DELTA",
				XAxis:         "Second",
				YAxis:         "Disk Writes (Delta per Second)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.png",
				},
			},
			{
				Column:        "AVG-SECTORS-WRITTEN-DELTA",
				XAxis:         "Second",
				YAxis:         "Sectors Written (Delta per Second)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.png",
				},
			},
			{
				Column:        "AVG-RECEIVE-BYTES-NUM-DELTA",
				XAxis:         "Second",
				YAxis:         "Network Receive(bytes) (Delta per Second)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.png",
				},
			},
			{
				Column:        "AVG-TRANSMIT-BYTES-NUM-DELTA",
				XAxis:         "Second",
				YAxis:         "Network Transmit(bytes) (Delta per Second)",
				OutputPathCSV: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.csv",
				OutputPathList: []string{
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.svg",
					"2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.png",
				},
			},
		},
		README: README{
			OutputPath: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/README.md",

			Images: []Image{
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.svg",
					Type:  "remote",
				},
				{
					Title: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA",
					Path:  "https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.svg",
					Type:  "remote",
				},
			},
		},
	}
	if !reflect.DeepEqual(cfg, expected) {
		t.Fatalf("configuration expected\n%+v\n, got\n%+v\n", expected, cfg)
	}

	req1, err := cfg.ToRequest("etcdv3", dbtesterpb.Request_Start, 0)
	if err != nil {
		t.Fatal(err)
	}
	expected1 := &dbtesterpb.Request{
		Operation:           dbtesterpb.Request_Start,
		TriggerLogUpload:    true,
		DatabaseID:          dbtesterpb.Request_etcdv3,
		DatabaseTag:         "etcd-v3.1-go1.7.5",
		PeerIPsString:       "10.240.0.20___10.240.0.21___10.240.0.22",
		IpIndex:             0,
		CurrentClientNumber: 0,
		Control: &dbtesterpb.Request_Control{
			GoogleCloudProjectName:         "etcd-development",
			GoogleCloudStorageKey:          "test-key",
			GoogleCloudStorageBucketName:   "dbtester-results",
			GoogleCloudStorageSubDirectory: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable",
		},
		Etcdv3Config: &dbtesterpb.Request_Etcdv3{
			SnapCount:      100000,
			QuotaSizeBytes: 8000000000,
		},
	}
	if !reflect.DeepEqual(req1, expected1) {
		t.Fatalf("configuration expected\n%+v\n, got\n%+v\n", expected1, req1)
	}

	req2, err := cfg.ToRequest("zookeeper", dbtesterpb.Request_Start, 2)
	if err != nil {
		t.Fatal(err)
	}
	expected2 := &dbtesterpb.Request{
		Operation:           dbtesterpb.Request_Start,
		TriggerLogUpload:    true,
		DatabaseID:          dbtesterpb.Request_zookeeper,
		DatabaseTag:         "zookeeper-r3.4.9-java8",
		PeerIPsString:       "10.240.0.25___10.240.0.27___10.240.0.28",
		IpIndex:             2,
		CurrentClientNumber: 0,
		Control: &dbtesterpb.Request_Control{
			GoogleCloudProjectName:         "etcd-development",
			GoogleCloudStorageKey:          "test-key",
			GoogleCloudStorageBucketName:   "dbtester-results",
			GoogleCloudStorageSubDirectory: "2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable",
		},
		ZookeeperConfig: &dbtesterpb.Request_Zookeeper{
			MyID:                 3,
			TickTime:             2000,
			ClientPort:           2181,
			InitLimit:            5,
			SyncLimit:            5,
			SnapCount:            100000,
			MaxClientConnections: 5000,
		},
	}
	if !reflect.DeepEqual(req2, expected2) {
		t.Fatalf("configuration expected\n%+v\n, got\n%+v\n", expected2, req2)
	}
}
