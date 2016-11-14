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

package control

import (
	"io/ioutil"

	"github.com/coreos/dbtester/agent"

	"gopkg.in/yaml.v2"
)

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

	// https://zookeeper.apache.org/doc/trunk/zookeeperAdmin.html
	Step1 struct {
		Skip                    bool  `yaml:"skip"`
		ZookeeperMaxClientCnxns int64 `yaml:"zookeeper_max_client_connections"`
		ZookeeperSnapCount      int64 `yaml:"zookeeper_snap_count"`
	} `yaml:"step1"`

	Step2 struct {
		Skip              bool   `yaml:"skip"`
		BenchType         string `yaml:"bench_type"`
		StaleRead         bool   `yaml:"stale_read"`
		ResultPath        string `yaml:"result_path"`
		Connections       int    `yaml:"connections"`
		Clients           int    `yaml:"clients"`
		KeySize           int    `yaml:"key_size"`
		SameKey           bool   `yaml:"same_key"`
		ValueSize         int    `yaml:"value_size"`
		TotalRequests     int    `yaml:"total_requests"`
		RequestIntervalMs int    `yaml:"request_interval_ms"`
	} `yaml:"step2"`

	Step3 struct {
		Skip       bool   `yaml:"skip"`
		ResultPath string `yaml:"result_path"`
	}
}

var (
	defaultZookeeperMaxClientCnxns int64 = 5000
	defaultZookeeperSnapCount      int64 = 100000
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

	if rs.Step1.ZookeeperMaxClientCnxns == 0 {
		rs.Step1.ZookeeperMaxClientCnxns = defaultZookeeperMaxClientCnxns
	}
	if rs.Step1.ZookeeperSnapCount == 0 {
		rs.Step1.ZookeeperSnapCount = defaultZookeeperSnapCount
	}

	return rs, nil
}

// ToRequest converts control configuration to agent RPC.
func (cfg *Config) ToRequest() agent.Request {
	req := agent.Request{}

	req.TestName = cfg.TestName
	req.GoogleCloudProjectName = cfg.GoogleCloudProjectName
	req.GoogleCloudStorageKey = cfg.GoogleCloudStorageKey
	req.GoogleCloudStorageBucketName = cfg.GoogleCloudStorageBucketName
	req.GoogleCloudStorageSubDirectory = cfg.GoogleCloudStorageSubDirectory

	switch cfg.Database {
	case "etcdv2":
		req.Database = agent.Request_etcdv2

	case "etcdv3":
		req.Database = agent.Request_etcdv3

	case "zk", "zookeeper":
		cfg.Database = "zookeeper"
		req.Database = agent.Request_ZooKeeper

	case "consul":
		req.Database = agent.Request_Consul
	}

	req.PeerIPString = cfg.PeerIPString

	req.ZookeeperMaxClientCnxns = cfg.Step1.ZookeeperMaxClientCnxns
	req.ZookeeperSnapCount = cfg.Step1.ZookeeperSnapCount

	return req
}
