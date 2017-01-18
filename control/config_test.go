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
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	c, err := ReadConfig("config_test.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if c.Database != "etcdv3" {
		t.Fatalf("unexpected %s", c.Database)
	}
	if c.TestName != "01-etcdv3" {
		t.Fatalf("unexpected %s", c.TestName)
	}
	if !reflect.DeepEqual(c.PeerIPs, []string{"10.240.0.13", "10.240.0.14", "10.240.0.15"}) {
		t.Fatalf("unexpected %s", c.PeerIPs)
	}
	if c.AgentPort != 3500 {
		t.Fatalf("unexpected %d", c.AgentPort)
	}
	if c.DatabasePort != 2379 {
		t.Fatalf("unexpected %d", c.DatabasePort)
	}
	if c.GoogleCloudProjectName != "etcd-development" {
		t.Fatalf("unexpected %s", c.GoogleCloudProjectName)
	}
	if c.GoogleCloudStorageKeyPath != "$HOME/gcloud-key.json" {
		t.Fatalf("unexpected %s", c.GoogleCloudStorageKeyPath)
	}
	if c.GoogleCloudStorageBucketName != "dbtester-results" {
		t.Fatalf("unexpected %s", c.GoogleCloudStorageBucketName)
	}
	if c.GoogleCloudStorageSubDirectory != "2016041501" {
		t.Fatalf("unexpected %s", c.GoogleCloudStorageSubDirectory)
	}

	if c.Log != "control.log" {
		t.Fatalf("unexpected %v", c.Log)
	}
	if c.DataLatencyDistributionSummary != "data-latency-distribution-summary.csv" {
		t.Fatalf("unexpected %s", c.DataLatencyDistributionSummary)
	}
	if c.DataLatencyDistributionPercentile != "data-latency-distribution-percentile.csv" {
		t.Fatalf("unexpected %s", c.DataLatencyDistributionPercentile)
	}
	if c.DataLatencyDistributionAll != "data-latency-distribution-all.csv" {
		t.Fatalf("unexpected %s", c.DataLatencyDistributionAll)
	}
	if c.DataLatencyThroughputTimeseries != "data-latency-throughput-timeseries.csv" {
		t.Fatalf("unexpected %s", c.DataLatencyThroughputTimeseries)
	}

	if c.Step1.SkipStartDatabase {
		t.Fatalf("unexpected %v", c.Step1.SkipStartDatabase)
	}
	if c.Step1.EtcdSnapCount != 100000 {
		t.Fatalf("unexpected %d", c.Step1.EtcdSnapCount)
	}
	if c.Step1.ZookeeperSnapCount != 100000 {
		t.Fatalf("unexpected %d", c.Step1.ZookeeperSnapCount)
	}
	if c.Step1.ZookeeperMaxClientCnxns != 5000 {
		t.Fatalf("unexpected %d", c.Step1.ZookeeperMaxClientCnxns)
	}

	if c.Step2.SkipStressDatabase {
		t.Fatalf("unexpected %v", c.Step2.SkipStressDatabase)
	}
	if c.Step2.BenchType != "write" {
		t.Fatalf("unexpected %s", c.Step2.BenchType)
	}
	if c.Step2.Clients != 1 {
		t.Fatalf("unexpected %d", c.Step2.Clients)
	}
	if c.Step2.Connections != 1 {
		t.Fatalf("unexpected %d", c.Step2.Connections)
	}
	if !reflect.DeepEqual(c.Step2.ConnectionsClients, []int{1, 10, 50, 100, 300, 500, 700, 1000, 1500, 2000, 2500}) {
		t.Fatalf("unexpected %d", c.Step2.ConnectionsClients)
	}
	if c.Step2.KeySize != 8 {
		t.Fatalf("unexpected %d", c.Step2.KeySize)
	}
	if !c.Step2.SameKey {
		t.Fatalf("unexpected %v", c.Step2.SameKey)
	}
	if !c.Step2.StaleRead {
		t.Fatalf("unexpected %v", c.Step2.StaleRead)
	}
	if c.Step2.TotalRequests != 2000000 {
		t.Fatalf("unexpected %d", c.Step2.TotalRequests)
	}
	if c.Step2.RequestsPerSecond != 100 {
		t.Fatalf("unexpected %d", c.Step2.RequestsPerSecond)
	}

	if c.Step3.Action != "stop" {
		t.Fatalf("unexpected %v", c.Step3.Action)
	}
}
