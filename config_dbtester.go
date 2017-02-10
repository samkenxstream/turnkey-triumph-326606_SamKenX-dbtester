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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/coreos/dbtester/dbtesterpb"
	"gopkg.in/yaml.v2"
)

// Config configures dbtester control clients.
type Config struct {
	TestTitle       string `yaml:"test_title"`
	TestDescription string `yaml:"test_description"`

	Control               `yaml:"control"`
	AllDatabaseIDList     []string             `yaml:"all_database_id_list"`
	DatabaseIDToTestGroup map[string]TestGroup `yaml:"datatbase_id_to_test_group"`
	DatabaseIDToTestData  map[string]TestData  `yaml:"datatbase_id_to_test_data"`

	Analyze `yaml:"analyze"`

	PlotPathPrefix string `yaml:"plot_path_prefix"`
	PlotList       []Plot `yaml:"plot_list"`
	README         `yaml:"readme"`
}

// Control represents common control options on client machine.
type Control struct {
	PathPrefix                              string `yaml:"path_prefix"`
	LogPath                                 string `yaml:"log_path"`
	ClientSystemMetricsPath                 string `yaml:"client_system_metrics_path"`
	ClientSystemMetricsInterpolatedPath     string `yaml:"client_system_metrics_interpolated_path"`
	ClientLatencyThroughputTimeseriesPath   string `yaml:"client_latency_throughput_timeseries_path"`
	ClientLatencyDistributionAllPath        string `yaml:"client_latency_distribution_all_path"`
	ClientLatencyDistributionPercentilePath string `yaml:"client_latency_distribution_percentile_path"`
	ClientLatencyDistributionSummaryPath    string `yaml:"client_latency_distribution_summary_path"`
	ClientLatencyByKeyNumberPath            string `yaml:"client_latency_by_key_number_path"`
	ServerDiskSpaceUsageSummaryPath         string `yaml:"server_disk_space_usage_summary_path"`

	GoogleCloudProjectName         string `yaml:"google_cloud_project_name"`
	GoogleCloudStorageKeyPath      string `yaml:"google_cloud_storage_key_path"`
	GoogleCloudStorageKey          string
	GoogleCloudStorageBucketName   string `yaml:"google_cloud_storage_bucket_name"`
	GoogleCloudStorageSubDirectory string `yaml:"google_cloud_storage_sub_directory"`
}

// TestGroup specifies database test group.
type TestGroup struct {
	DatabaseID          string
	DatabaseDescription string `yaml:"database_description"`
	DatabaseTag         string

	PeerIPs       []string `yaml:"peer_ips"`
	PeerIPsString string

	DatabasePortToConnect int `yaml:"database_port_to_connect"`
	DatabaseEndpoints     []string

	AgentPortToConnect int `yaml:"agent_port_to_connect"`
	AgentEndpoints     []string

	// database-specific flags to start
	Etcdv2    `yaml:"etcdv2"`
	Etcdv3    `yaml:"etcdv3"`
	Zookeeper `yaml:"zookeeper"`
	Consul    `yaml:"consul"`
	Zetcd     `yaml:"zetcd"`
	Cetcd     `yaml:"cetcd"`

	// benchmark options
	BenchmarkOptions `yaml:"benchmark_options"`
	BenchmarkSteps   `yaml:"benchmark_steps"`
}

// BenchmarkOptions specifies the benchmark options.
type BenchmarkOptions struct {
	Type string `yaml:"type"`

	RequestNumber              int64   `yaml:"request_number"`
	ConnectionNumber           int64   `yaml:"connection_number"`
	ClientNumber               int64   `yaml:"client_number"`
	ConnectionClientNumbers    []int64 `yaml:"connection_client_numbers"`
	RateLimitRequestsPerSecond int64   `yaml:"rate_limit_requests_per_second"`

	// for writes, reads
	SameKey        bool  `yaml:"same_key"`
	KeySizeBytes   int64 `yaml:"key_size_bytes"`
	ValueSizeBytes int64 `yaml:"value_size_bytes"`

	// for reads
	StaleRead bool `yaml:"stale_read"`
}

// BenchmarkSteps specifies the benchmark workflow.
type BenchmarkSteps struct {
	Step1StartDatabase  bool `yaml:"step1_start_database"`
	Step2StressDatabase bool `yaml:"step2_stress_database"`
	Step3StopDatabase   bool `yaml:"step3_stop_database"`
	Step4UploadLogs     bool `yaml:"step4_upload_logs"`
}

// TestData defines raw data to import.
type TestData struct {
	DatabaseID          string
	DatabaseTag         string
	DatabaseDescription string

	PathPrefix                              string   `yaml:"path_prefix"`
	ClientSystemMetricsInterpolatedPath     string   `yaml:"client_system_metrics_interpolated_path"`
	ClientLatencyThroughputTimeseriesPath   string   `yaml:"client_latency_throughput_timeseries_path"`
	ClientLatencyDistributionAllPath        string   `yaml:"client_latency_distribution_all_path"`
	ClientLatencyDistributionPercentilePath string   `yaml:"client_latency_distribution_percentile_path"`
	ClientLatencyDistributionSummaryPath    string   `yaml:"client_latency_distribution_summary_path"`
	ClientLatencyByKeyNumberPath            string   `yaml:"client_latency_by_key_number_path"`
	ServerDiskSpaceUsageSummaryPath         string   `yaml:"server_disk_space_usage_summary_path"`
	ServerMemoryByKeyNumberPath             string   `yaml:"server_memory_by_key_number_path"`
	ServerSystemMetricsInterpolatedPathList []string `yaml:"server_system_metrics_interpolated_path_list"`
	AllAggregatedOutputPath                 string   `yaml:"all_aggregated_output_path"`
}

// Analyze defines analyze config.
type Analyze struct {
	AllAggregatedOutputPathCSV string `yaml:"all_aggregated_output_path_csv"`
	AllAggregatedOutputPathTXT string `yaml:"all_aggregated_output_path_txt"`
}

// Plot defines plot configuration.
type Plot struct {
	Column         string   `yaml:"column"`
	XAxis          string   `yaml:"x_axis"`
	YAxis          string   `yaml:"y_axis"`
	OutputPathCSV  string   `yaml:"output_path_csv"`
	OutputPathList []string `yaml:"output_path_list"`
}

// README defines how to write README.
type README struct {
	OutputPath string  `yaml:"output_path"`
	Images     []Image `yaml:"images"`
}

// Image defines image data.
type Image struct {
	Title string `yaml:"title"`
	Path  string `yaml:"path"`
	Type  string `yaml:"type"`
}

// ReadConfig reads control configuration file.
func ReadConfig(fpath string, analyze bool) (*Config, error) {
	bts, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	if err := yaml.Unmarshal(bts, &cfg); err != nil {
		return nil, err
	}

	if cfg.Control.PathPrefix != "" {
		cfg.Control.LogPath = filepath.Join(cfg.Control.PathPrefix, cfg.Control.LogPath)
		cfg.Control.ClientSystemMetricsPath = filepath.Join(cfg.Control.PathPrefix, cfg.Control.ClientSystemMetricsPath)
		cfg.Control.ClientSystemMetricsInterpolatedPath = filepath.Join(cfg.Control.PathPrefix, cfg.Control.ClientSystemMetricsInterpolatedPath)
		cfg.Control.ClientLatencyThroughputTimeseriesPath = filepath.Join(cfg.Control.PathPrefix, cfg.Control.ClientLatencyThroughputTimeseriesPath)
		cfg.Control.ClientLatencyDistributionAllPath = filepath.Join(cfg.Control.PathPrefix, cfg.Control.ClientLatencyDistributionAllPath)
		cfg.Control.ClientLatencyDistributionPercentilePath = filepath.Join(cfg.Control.PathPrefix, cfg.Control.ClientLatencyDistributionPercentilePath)
		cfg.Control.ClientLatencyDistributionSummaryPath = filepath.Join(cfg.Control.PathPrefix, cfg.Control.ClientLatencyDistributionSummaryPath)
		cfg.Control.ClientLatencyByKeyNumberPath = filepath.Join(cfg.Control.PathPrefix, cfg.Control.ClientLatencyByKeyNumberPath)
		cfg.Control.ServerDiskSpaceUsageSummaryPath = filepath.Join(cfg.Control.PathPrefix, cfg.Control.ServerDiskSpaceUsageSummaryPath)
	}

	for databaseID, group := range cfg.DatabaseIDToTestGroup {
		group.DatabaseID = databaseID
		group.DatabaseTag = MakeTag(group.DatabaseDescription)
		group.PeerIPsString = strings.Join(group.PeerIPs, "___")
		group.DatabaseEndpoints = make([]string, len(group.PeerIPs))
		group.AgentEndpoints = make([]string, len(group.PeerIPs))
		for j := range group.PeerIPs {
			group.DatabaseEndpoints[j] = fmt.Sprintf("%s:%d", group.PeerIPs[j], group.DatabasePortToConnect)
			group.AgentEndpoints[j] = fmt.Sprintf("%s:%d", group.PeerIPs[j], group.AgentPortToConnect)
		}
		cfg.DatabaseIDToTestGroup[databaseID] = group
	}

	for databaseID, testdata := range cfg.DatabaseIDToTestData {
		testdata.PathPrefix = strings.TrimSpace(testdata.PathPrefix)
		testdata.DatabaseID = databaseID
		testdata.DatabaseTag = cfg.DatabaseIDToTestGroup[databaseID].DatabaseTag
		testdata.DatabaseDescription = cfg.DatabaseIDToTestGroup[databaseID].DatabaseDescription

		if testdata.PathPrefix != "" {
			testdata.ClientSystemMetricsInterpolatedPath = testdata.PathPrefix + "-" + testdata.ClientSystemMetricsInterpolatedPath
			testdata.ClientLatencyThroughputTimeseriesPath = testdata.PathPrefix + "-" + testdata.ClientLatencyThroughputTimeseriesPath
			testdata.ClientLatencyDistributionAllPath = testdata.PathPrefix + "-" + testdata.ClientLatencyDistributionAllPath
			testdata.ClientLatencyDistributionPercentilePath = testdata.PathPrefix + "-" + testdata.ClientLatencyDistributionPercentilePath
			testdata.ClientLatencyDistributionSummaryPath = testdata.PathPrefix + "-" + testdata.ClientLatencyDistributionSummaryPath
			testdata.ClientLatencyByKeyNumberPath = testdata.PathPrefix + "-" + testdata.ClientLatencyByKeyNumberPath
			testdata.ServerDiskSpaceUsageSummaryPath = testdata.PathPrefix + "-" + testdata.ServerDiskSpaceUsageSummaryPath
			testdata.ServerMemoryByKeyNumberPath = testdata.PathPrefix + "-" + testdata.ServerMemoryByKeyNumberPath
			for i := range testdata.ServerSystemMetricsInterpolatedPathList {
				testdata.ServerSystemMetricsInterpolatedPathList[i] = testdata.PathPrefix + "-" + testdata.ServerSystemMetricsInterpolatedPathList[i]
			}
			testdata.AllAggregatedOutputPath = testdata.PathPrefix + "-" + testdata.AllAggregatedOutputPath
		}

		cfg.DatabaseIDToTestData[databaseID] = testdata
	}

	for databaseID, group := range cfg.DatabaseIDToTestGroup {
		if databaseID != "etcdv3" && group.BenchmarkOptions.ConnectionNumber != group.BenchmarkOptions.ClientNumber {
			return nil, fmt.Errorf("%q got connected %d != clients %d", databaseID, group.BenchmarkOptions.ConnectionNumber, group.BenchmarkOptions.ClientNumber)
		}
	}

	var (
		defaultEtcdSnapCount                 int64 = 100000
		defaultZookeeperSnapCount            int64 = 100000
		defaultZookeeperTickTime             int64 = 2000
		defaultZookeeperInitLimit            int64 = 5
		defaultZookeeperSyncLimit            int64 = 5
		defaultZookeeperMaxClientConnections int64 = 5000
	)
	if v, ok := cfg.DatabaseIDToTestGroup["etcdv3"]; ok {
		if v.Etcdv3.SnapCount == 0 {
			v.Etcdv3.SnapCount = defaultEtcdSnapCount
		}
		cfg.DatabaseIDToTestGroup["etcdv3"] = v
	}
	if v, ok := cfg.DatabaseIDToTestGroup["zookeeper"]; ok {
		if v.Zookeeper.TickTime == 0 {
			v.Zookeeper.TickTime = defaultZookeeperTickTime
		}
		if v.Zookeeper.InitLimit == 0 {
			v.Zookeeper.InitLimit = defaultZookeeperInitLimit
		}
		if v.Zookeeper.SyncLimit == 0 {
			v.Zookeeper.SyncLimit = defaultZookeeperSyncLimit
		}
		if v.Zookeeper.SnapCount == 0 {
			v.Zookeeper.SnapCount = defaultZookeeperSnapCount
		}
		if v.Zookeeper.MaxClientConnections == 0 {
			v.Zookeeper.MaxClientConnections = defaultZookeeperMaxClientConnections
		}
		cfg.DatabaseIDToTestGroup["zookeeper"] = v
	}

	if cfg.Control.GoogleCloudStorageKeyPath != "" && !analyze {
		bts, err = ioutil.ReadFile(cfg.Control.GoogleCloudStorageKeyPath)
		if err != nil {
			return nil, err
		}
		cfg.Control.GoogleCloudStorageKey = string(bts)
	}

	for i := range cfg.PlotList {
		cfg.PlotList[i].OutputPathCSV = filepath.Join(cfg.PlotPathPrefix, cfg.PlotList[i].Column+".csv")
		cfg.PlotList[i].OutputPathList = make([]string, 2)
		cfg.PlotList[i].OutputPathList[0] = filepath.Join(cfg.PlotPathPrefix, cfg.PlotList[i].Column+".svg")
		cfg.PlotList[i].OutputPathList[1] = filepath.Join(cfg.PlotPathPrefix, cfg.PlotList[i].Column+".png")
	}

	return &cfg, nil
}

// MakeTag converts database scription to database tag.
func MakeTag(desc string) string {
	s := strings.ToLower(desc)
	s = strings.Replace(s, "go ", "go", -1)
	s = strings.Replace(s, "java ", "java", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	return strings.Replace(s, " ", "-", -1)
}

// ToRequest converts configuration to 'dbtesterpb.Request'.
func (cfg *Config) ToRequest(databaseID string, op dbtesterpb.Request_Operation, idx int) (req *dbtesterpb.Request, err error) {
	gcfg, ok := cfg.DatabaseIDToTestGroup[databaseID]
	if !ok {
		err = fmt.Errorf("%q is not defined", databaseID)
		return
	}

	req = &dbtesterpb.Request{
		Operation:           op,
		TriggerLogUpload:    gcfg.BenchmarkSteps.Step4UploadLogs,
		DatabaseID:          dbtesterpb.Request_Database(dbtesterpb.Request_Database_value[databaseID]),
		DatabaseTag:         gcfg.DatabaseTag,
		PeerIPsString:       gcfg.PeerIPsString,
		IpIndex:             uint32(idx),
		CurrentClientNumber: gcfg.BenchmarkOptions.ClientNumber,
		Control: &dbtesterpb.Request_Control{
			GoogleCloudProjectName:         cfg.Control.GoogleCloudProjectName,
			GoogleCloudStorageKey:          cfg.Control.GoogleCloudStorageKey,
			GoogleCloudStorageBucketName:   cfg.Control.GoogleCloudStorageBucketName,
			GoogleCloudStorageSubDirectory: cfg.Control.GoogleCloudStorageSubDirectory,
		},
	}

	switch req.DatabaseID {
	case dbtesterpb.Request_etcdv2:

	case dbtesterpb.Request_etcdv3:
		req.Etcdv3Config = &dbtesterpb.Request_Etcdv3{
			SnapCount:      gcfg.Etcdv3.SnapCount,
			QuotaSizeBytes: gcfg.Etcdv3.QuotaSizeBytes,
		}

	case dbtesterpb.Request_zookeeper:
		req.ZookeeperConfig = &dbtesterpb.Request_Zookeeper{
			MyID:                 uint32(idx + 1),
			TickTime:             gcfg.Zookeeper.TickTime,
			ClientPort:           int64(gcfg.DatabasePortToConnect),
			InitLimit:            gcfg.Zookeeper.InitLimit,
			SyncLimit:            gcfg.Zookeeper.SyncLimit,
			SnapCount:            gcfg.Zookeeper.SnapCount,
			MaxClientConnections: gcfg.Zookeeper.MaxClientConnections,
		}

	case dbtesterpb.Request_consul:
	case dbtesterpb.Request_zetcd:
	case dbtesterpb.Request_cetcd:

	default:
		err = fmt.Errorf("unknown %v", req.DatabaseID)
		return
	}

	return
}
