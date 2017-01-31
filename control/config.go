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

package control

import (
	"fmt"
	"io/ioutil"

	"github.com/coreos/dbtester/agent/agentpb"

	"gopkg.in/yaml.v2"
)

// Config configures dbtester control clients.
type Config struct {
	Database                       string `yaml:"database"`
	TestName                       string `yaml:"test_name"`
	GoogleCloudProjectName         string `yaml:"google_cloud_project_name"`
	GoogleCloudStorageKey          string
	GoogleCloudStorageKeyPath      string `yaml:"google_cloud_storage_key_path"`
	GoogleCloudStorageBucketName   string `yaml:"google_cloud_storage_bucket_name"`
	GoogleCloudStorageSubDirectory string `yaml:"google_cloud_storage_sub_directory"`

	PeerIPs      []string `yaml:"peer_ips"`
	PeerIPString string
	AgentPort    int `yaml:"agent_port"`
	DatabasePort int `yaml:"database_port"`

	AgentEndpoints    []string
	DatabaseEndpoints []string

	Log                               string `yaml:"log"`
	DataLatencyDistributionSummary    string `yaml:"data_latency_distribution_summary"`
	DataLatencyDistributionPercentile string `yaml:"data_latency_distribution_percentile"`
	DataLatencyDistributionAll        string `yaml:"data_latency_distribution_all"`
	DataLatencyThroughputTimeseries   string `yaml:"data_latency_throughput_timeseries"`

	// https://zookeeper.apache.org/doc/trunk/zookeeperAdmin.html
	Step1 struct {
		SkipStartDatabase       bool  `yaml:"skip_start_database"`
		EtcdSnapCount           int64 `yaml:"etcd_snap_count"`
		EtcdQuotaSizeBytes      int64 `yaml:"etcd_quota_size_bytes"`
		ZookeeperSnapCount      int64 `yaml:"zookeeper_snap_count"`
		ZookeeperMaxClientCnxns int64 `yaml:"zookeeper_max_client_connections"`
	} `yaml:"step1"`

	Step2 struct {
		SkipStressDatabase bool   `yaml:"skip_stress_database"`
		BenchType          string `yaml:"bench_type"`
		StaleRead          bool   `yaml:"stale_read"`
		Connections        int    `yaml:"connections"`
		Clients            int    `yaml:"clients"`
		ConnectionsClients []int  `yaml:"connections_clients"`
		KeySize            int    `yaml:"key_size"`
		SameKey            bool   `yaml:"same_key"`
		ValueSize          int    `yaml:"value_size"`
		TotalRequests      int    `yaml:"total_requests"`
		RequestsPerSecond  int    `yaml:"requests_per_second"`
	} `yaml:"step2"`

	Step3 struct {
		Action string `yaml:"action"`
	}
}

var (
	defaultEtcdSnapCount           int64 = 100000
	defaultZookeeperSnapCount      int64 = 100000
	defaultZookeeperMaxClientCnxns int64 = 5000
)

// ReadConfig reads control configuration file.
func ReadConfig(fpath string) (Config, error) {
	bts, err := ioutil.ReadFile(fpath)
	if err != nil {
		return Config{}, err
	}
	rs := Config{}
	if err := yaml.Unmarshal(bts, &rs); err != nil {
		return Config{}, err
	}

	if rs.Step2.Connections != rs.Step2.Clients {
		switch rs.Database {
		case "etcdv2":
			return Config{}, fmt.Errorf("connected %d != clients %d", rs.Step2.Connections, rs.Step2.Clients)

		case "etcdv3":

		case "zookeeper":
			return Config{}, fmt.Errorf("connected %d != clients %d", rs.Step2.Connections, rs.Step2.Clients)

		case "zetcd":
			return Config{}, fmt.Errorf("connected %d != clients %d", rs.Step2.Connections, rs.Step2.Clients)

		case "consul":
			return Config{}, fmt.Errorf("connected %d != clients %d", rs.Step2.Connections, rs.Step2.Clients)

		case "cetcd":
			return Config{}, fmt.Errorf("connected %d != clients %d", rs.Step2.Connections, rs.Step2.Clients)
		}
	}

	if rs.Step1.EtcdSnapCount == 0 {
		rs.Step1.EtcdSnapCount = defaultEtcdSnapCount
	}
	if rs.Step1.ZookeeperSnapCount == 0 {
		rs.Step1.ZookeeperSnapCount = defaultZookeeperSnapCount
	}
	if rs.Step1.ZookeeperMaxClientCnxns == 0 {
		rs.Step1.ZookeeperMaxClientCnxns = defaultZookeeperMaxClientCnxns
	}

	return rs, nil
}

// ToRequest converts control configuration to agent RPC.
func (cfg *Config) ToRequest() agentpb.Request {
	req := agentpb.Request{}

	req.TestName = cfg.TestName
	req.GoogleCloudProjectName = cfg.GoogleCloudProjectName
	req.GoogleCloudStorageKey = cfg.GoogleCloudStorageKey
	req.GoogleCloudStorageBucketName = cfg.GoogleCloudStorageBucketName
	req.GoogleCloudStorageSubDirectory = cfg.GoogleCloudStorageSubDirectory

	switch cfg.Database {
	case "etcdv2":
		req.Database = agentpb.Request_etcdv2

	case "etcdv3":
		req.Database = agentpb.Request_etcdv3

	case "zookeeper":
		req.Database = agentpb.Request_ZooKeeper

	case "zetcd":
		req.Database = agentpb.Request_zetcd

	case "consul":
		req.Database = agentpb.Request_Consul

	case "cetcd":
		req.Database = agentpb.Request_cetcd
	}

	req.PeerIPString = cfg.PeerIPString

	req.EtcdSnapCount = cfg.Step1.EtcdSnapCount
	req.EtcdQuotaSizeBytes = cfg.Step1.EtcdQuotaSizeBytes
	req.ZookeeperSnapCount = cfg.Step1.ZookeeperSnapCount
	req.ZookeeperMaxClientCnxns = cfg.Step1.ZookeeperMaxClientCnxns

	req.ClientNum = int64(cfg.Step2.Clients)

	return req
}
