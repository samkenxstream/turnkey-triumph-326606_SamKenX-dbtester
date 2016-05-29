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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/coreos/dbtester/agent"
	"github.com/spf13/cobra"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/cheggaaa/pb"
	"github.com/coreos/etcd/clientv3"
	consulapi "github.com/hashicorp/consul/api"
)

var (
	Command = &cobra.Command{
		Use:   "control",
		Short: "Controls tests.",
		RunE:  CommandFunc,
	}
	configPath string
)

func init() {
	Command.PersistentFlags().StringVarP(&configPath, "config", "c", "", "YAML configuration file path.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	cfg, err := ReadConfig(configPath)
	if err != nil {
		return err
	}
	switch cfg.Database {
	case "etcdv2":
	case "etcdv3":
	case "zk", "zookeeper":
	case "consul":
	default:
		return fmt.Errorf("%q is not supported", cfg.Database)
	}
	if !cfg.Step2.Skip {
		switch cfg.Step2.BenchType {
		case "write":
		case "read":
		default:
			return fmt.Errorf("%q is not supported", cfg.Step2.BenchType)
		}
	}

	bts, err := ioutil.ReadFile(cfg.GoogleCloudStorageKeyPath)
	if err != nil {
		return err
	}
	cfg.GoogleCloudStorageKey = string(bts)

	cfg.PeerIPString = strings.Join(cfg.PeerIPs, "___") // protoc sorts the 'repeated' type data
	cfg.AgentEndpoints = make([]string, len(cfg.PeerIPs))
	cfg.DatabaseEndpoints = make([]string, len(cfg.PeerIPs))
	for i := range cfg.PeerIPs {
		cfg.AgentEndpoints[i] = fmt.Sprintf("%s:%d", cfg.PeerIPs[i], cfg.AgentPort)
	}
	for i := range cfg.PeerIPs {
		cfg.DatabaseEndpoints[i] = fmt.Sprintf("%s:%d", cfg.PeerIPs[i], cfg.DatabasePort)
	}

	println()
	if !cfg.Step1.Skip {
		logger.Info("step 1: starting databases...")
		if err = step1(cfg); err != nil {
			return err
		}
	}

	if !cfg.Step2.Skip {
		println()
		time.Sleep(5 * time.Second)
		logger.Info("step 2: starting tests...")
		if err = step2(cfg); err != nil {
			return err
		}
	}

	if !cfg.Step3.Skip {
		println()
		time.Sleep(5 * time.Second)
		logger.Info("step 3: stopping databases...")
		if err = step3(cfg); err != nil {
			return err
		}
	}

	return nil
}

func step1(cfg Config) error {
	req := agent.Request{}

	req.Operation = agent.Request_Start
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
	// req.EtcdCompression = cfg.EtcdCompression

	donec, errc := make(chan struct{}), make(chan error)
	for i := range cfg.PeerIPs {

		go func(i int) {
			nreq := req

			nreq.ServerIndex = uint32(i)
			nreq.ZookeeperMyID = uint32(i + 1)
			ep := cfg.AgentEndpoints[nreq.ServerIndex]

			logger.Info("sending message",
				zap.Int("index", i),
				zap.String("operation", req.Operation.String()),
				zap.String("database", req.Database.String()),
				zap.String("endpoint", ep),
			)

			conn, err := grpc.Dial(ep, grpc.WithInsecure())
			if err != nil {
				logger.Error("grpc.Dial connecting error",
					zap.Int("index", i),
					zap.String("endpoint", ep),
					zap.Err(err),
				)
				errc <- err
				return
			}

			defer conn.Close()

			cli := agent.NewTransporterClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Consul takes longer
			resp, err := cli.Transfer(ctx, &nreq)
			cancel()
			if err != nil {
				logger.Error("cli.Transfer error",
					zap.Int("index", i),
					zap.String("endpoint", ep),
					zap.Err(err),
				)
				errc <- err
				return
			}

			logger.Info("response",
				zap.Int("index", i),
				zap.String("endpoint", ep),
				zap.String("response", fmt.Sprintf("%+v", resp)),
			)
			donec <- struct{}{}
		}(i)

		time.Sleep(time.Second)
	}

	cnt := 0
	for cnt != len(cfg.PeerIPs) {
		select {
		case <-donec:
		case err := <-errc:
			return err
		}
		cnt++
	}

	return nil
}

var (
	bar     *pb.ProgressBar
	results chan result
	wg      sync.WaitGroup
)

func step2(cfg Config) error {
	var (
		valuesBytes     [][]byte
		valuesString    []string
		valueSampleSize int
	)
	if cfg.Step2.ValueTestDataPath != "" {
		fs, err := walkDir(cfg.Step2.ValueTestDataPath)
		if err != nil {
			return err
		}
		for _, elem := range fs {
			bts, err := ioutil.ReadFile(elem.path)
			if err != nil {
				return err
			}
			valuesBytes = append(valuesBytes, bts)
			valuesString = append(valuesString, string(bts))
		}
		valueSampleSize = len(valuesString)
	}

	switch cfg.Step2.BenchType {
	case "write":
		results = make(chan result)
		requests := make(chan request, cfg.Step2.Clients)
		bar = pb.New(cfg.Step2.TotalRequests)

		v := randBytes(cfg.Step2.ValueSize)
		vs := string(v)

		bar.Format("Bom !")
		bar.Start()

		var etcdClients []*clientv3.Client
		switch cfg.Database {
		case "etcdv2":
			conns := mustCreateClientsEtcdv2(cfg.DatabaseEndpoints, cfg.Step2.Connections)
			for i := range conns {
				wg.Add(1)
				go doPutEtcdv2(context.Background(), conns[i], requests)
			}

		case "etcdv3":
			etcdClients := mustCreateClientsEtcdv3(cfg.DatabaseEndpoints, etcdv3ClientCfg{
				totalConns:   cfg.Step2.Connections,
				totalClients: cfg.Step2.Clients,
				// compressionTypeTxt: cfg.EtcdCompression,
			})
			for i := range etcdClients {
				wg.Add(1)
				go doPutEtcdv3(context.Background(), etcdClients[i], requests)
			}
			defer func() {
				for i := range etcdClients {
					etcdClients[i].Close()
				}
			}()

		case "zk", "zookeeper":
			if cfg.Step2.SameKey {
				key := sameKey(cfg.Step2.KeySize)
				valueBts := randBytes(cfg.Step2.ValueSize)
				logger.Info("write started",
					zap.String("request", "PUT"),
					zap.String("key", key),
					zap.String("database", "zookeeper"),
				)
				var err error
				for i := 0; i < 7; i++ {
					conns := mustCreateConnsZk(cfg.DatabaseEndpoints, cfg.Step2.Connections)
					_, err = conns[0].Create("/"+key, valueBts, zkCreateFlags, zkCreateAcl)
					if err != nil {
						continue
					}
					for j := range conns {
						conns[j].Close()
					}
					logger.Info("write done",
						zap.String("request", "PUT"),
						zap.String("key", key),
						zap.String("database", "zookeeper"),
					)
					break
				}
				if err != nil {
					logger.Error("write error",
						zap.String("request", "PUT"),
						zap.String("key", key),
						zap.String("database", "zookeeper"),
						zap.Err(err),
					)
					os.Exit(1)
				}
			}

			conns := mustCreateConnsZk(cfg.DatabaseEndpoints, cfg.Step2.Connections)
			defer func() {
				for i := range conns {
					conns[i].Close()
				}
			}()
			for i := range conns {
				wg.Add(1)
				go doPutZk(conns[i], requests, cfg.Step2.SameKey)
			}

		case "consul":
			conns := mustCreateConnsConsul(cfg.DatabaseEndpoints, cfg.Step2.Connections)
			for i := range conns {
				wg.Add(1)
				go doPutConsul(conns[i], requests)
			}
		}

		pdoneC := printReport(results, cfg)
		go func() {
			for i := 0; i < cfg.Step2.TotalRequests; i++ {
				if cfg.Database == "etcdv3" && cfg.Step2.Etcdv3CompactionCycle > 0 && i%cfg.Step2.Etcdv3CompactionCycle == 0 {
					logger.Info("starting compaction",
						zap.Int("index", i),
						zap.String("database", "etcdv3"),
					)
					go func() {
						compactKV(etcdClients)
					}()
				}

				k := sequentialKey(cfg.Step2.KeySize, i)
				if cfg.Step2.SameKey {
					k = sameKey(cfg.Step2.KeySize)
				}
				if cfg.Step2.ValueTestDataPath != "" {
					v = valuesBytes[i%valueSampleSize]
					vs = valuesString[i%valueSampleSize]
				}

				switch cfg.Database {
				case "etcdv2":
					requests <- request{etcdv2Op: etcdv2Op{key: k, value: vs}}

				case "etcdv3":
					requests <- request{etcdv3Op: clientv3.OpPut(k, vs)}

				case "zk", "zookeeper":
					requests <- request{zkOp: zkOp{key: "/" + k, value: v}}

				case "consul":
					requests <- request{consulOp: consulOp{key: k, value: v}}
				}
				if cfg.Step2.RequestIntervalMs > 0 {
					time.Sleep(time.Duration(cfg.Step2.RequestIntervalMs) * time.Millisecond)
				}
			}
			close(requests)
		}()

		wg.Wait()

		bar.Finish()

		close(results)
		<-pdoneC

		switch cfg.Database {
		case "etcdv2":
			for k, v := range getTotalKeysEtcdv2(cfg.DatabaseEndpoints) {
				logger.Info("expected write total results",
					zap.Int("expected_total", cfg.Step2.TotalRequests),
					zap.String("database", "etcdv2"),
					zap.String("endpoint", k),
					zap.Int64("number_of_keys", v),
				)
			}

		case "etcdv3":
			for k, v := range getTotalKeysEtcdv3(cfg.DatabaseEndpoints) {
				logger.Info("expected write total results",
					zap.Int("expected_total", cfg.Step2.TotalRequests),
					zap.String("database", "etcdv3"),
					zap.String("endpoint", k),
					zap.Int64("number_of_keys", v),
				)
			}

		case "zk", "zookeeper":
			for k, v := range getTotalKeysZk(cfg.DatabaseEndpoints) {
				logger.Info("expected write total results",
					zap.Int("expected_total", cfg.Step2.TotalRequests),
					zap.String("database", "zookeeper"),
					zap.String("endpoint", k),
					zap.Int64("number_of_keys", v),
				)
			}

		case "consul":
			for k, v := range getTotalKeysConsul(cfg.DatabaseEndpoints) {
				logger.Info("expected write total results",
					zap.Int("expected_total", cfg.Step2.TotalRequests),
					zap.String("database", "consul"),
					zap.String("endpoint", k),
					zap.Int64("number_of_keys", v),
				)
			}
		}

	case "read":
		var (
			key      = sameKey(cfg.Step2.KeySize)
			valueBts = randBytes(cfg.Step2.ValueSize)
			value    = string(valueBts)
		)
		if cfg.Step2.ValueTestDataPath != "" {
			valueBts = valuesBytes[0]
			value = valuesString[0]
		}
		switch cfg.Database {
		case "etcdv2":
			logger.Info("write started",
				zap.String("request", "PUT"),
				zap.String("key", key),
				zap.String("database", "etcdv2"),
			)
			var err error
			for i := 0; i < 7; i++ {
				clients := mustCreateClientsEtcdv2(cfg.DatabaseEndpoints, cfg.Step2.Connections)
				_, err = clients[0].Set(context.Background(), key, value, nil)
				if err != nil {
					continue
				}
				logger.Info("write done",
					zap.String("request", "PUT"),
					zap.String("key", key),
					zap.String("database", "etcdv2"),
				)
				break
			}
			if err != nil {
				logger.Info("write error",
					zap.String("request", "PUT"),
					zap.String("key", key),
					zap.String("database", "etcdv2"),
					zap.Err(err),
				)
				os.Exit(1)
			}

		case "etcdv3":
			logger.Info("write started",
				zap.String("request", "PUT"),
				zap.String("key", key),
				zap.String("database", "etcdv3"),
			)
			var err error
			for i := 0; i < 7; i++ {
				clients := mustCreateClientsEtcdv3(cfg.DatabaseEndpoints, etcdv3ClientCfg{
					totalConns:   1,
					totalClients: 1,
					// compressionTypeTxt: cfg.EtcdCompression,
				})
				_, err = clients[0].Do(context.Background(), clientv3.OpPut(key, value))
				if err != nil {
					continue
				}
				logger.Info("write done",
					zap.String("request", "PUT"),
					zap.String("key", key),
					zap.String("database", "etcdv3"),
				)
				break
			}
			if err != nil {
				logger.Error("write error",
					zap.String("request", "PUT"),
					zap.String("key", key),
					zap.String("database", "etcdv3"),
					zap.Err(err),
				)
				os.Exit(1)
			}

		case "zk", "zookeeper":
			logger.Info("write started",
				zap.String("request", "PUT"),
				zap.String("key", key),
				zap.String("database", "zookeeper"),
			)
			var err error
			for i := 0; i < 7; i++ {
				conns := mustCreateConnsZk(cfg.DatabaseEndpoints, cfg.Step2.Connections)
				_, err = conns[0].Create("/"+key, valueBts, zkCreateFlags, zkCreateAcl)
				if err != nil {
					continue
				}
				for j := range conns {
					conns[j].Close()
				}
				logger.Info("write done",
					zap.String("request", "PUT"),
					zap.String("key", key),
					zap.String("database", "zookeeper"),
				)
				break
			}
			if err != nil {
				logger.Error("write error",
					zap.String("request", "PUT"),
					zap.String("key", key),
					zap.String("database", "zookeeper"),
					zap.Err(err),
				)
				os.Exit(1)
			}

		case "consul":
			logger.Info("write started",
				zap.String("request", "PUT"),
				zap.String("key", key),
				zap.String("database", "consul"),
			)
			var err error
			for i := 0; i < 7; i++ {
				clients := mustCreateConnsConsul(cfg.DatabaseEndpoints, cfg.Step2.Connections)
				_, err = clients[0].Put(&consulapi.KVPair{Key: key, Value: valueBts}, nil)
				if err != nil {
					continue
				}
				logger.Info("write done",
					zap.String("request", "PUT"),
					zap.String("key", key),
					zap.String("database", "consul"),
				)
				break
			}
			if err != nil {
				logger.Error("write error",
					zap.String("request", "PUT"),
					zap.String("key", key),
					zap.String("database", "consul"),
					zap.Err(err),
				)
				os.Exit(1)
			}
		}

		results = make(chan result)
		requests := make(chan request, cfg.Step2.Clients)
		bar = pb.New(cfg.Step2.TotalRequests)

		bar.Format("Bom !")
		bar.Start()

		switch cfg.Database {
		case "etcdv2":
			conns := mustCreateClientsEtcdv2(cfg.DatabaseEndpoints, cfg.Step2.Connections)
			for i := range conns {
				wg.Add(1)
				go doRangeEtcdv2(conns[i], requests)
			}

		case "etcdv3":
			clients := mustCreateClientsEtcdv3(cfg.DatabaseEndpoints, etcdv3ClientCfg{
				totalConns:   cfg.Step2.Connections,
				totalClients: cfg.Step2.Clients,
				// compressionTypeTxt: cfg.EtcdCompression,
			})
			for i := range clients {
				wg.Add(1)
				go doRangeEtcdv3(clients[i].KV, requests)
			}
			defer func() {
				for i := range clients {
					clients[i].Close()
				}
			}()

		case "zk", "zookeeper":
			conns := mustCreateConnsZk(cfg.DatabaseEndpoints, cfg.Step2.Connections)
			defer func() {
				for i := range conns {
					conns[i].Close()
				}
			}()
			for i := range conns {
				wg.Add(1)
				go doRangeZk(conns[i], requests)
			}

		case "consul":
			conns := mustCreateConnsConsul(cfg.DatabaseEndpoints, cfg.Step2.Connections)
			for i := range conns {
				wg.Add(1)
				go doRangeConsul(conns[i], requests)
			}
		}

		pdoneC := printReport(results, cfg)
		go func() {
			for i := 0; i < cfg.Step2.TotalRequests; i++ {
				switch cfg.Database {
				case "etcdv2":
					// serializable read by default
					requests <- request{etcdv2Op: etcdv2Op{key: key}}

				case "etcdv3":
					opts := []clientv3.OpOption{clientv3.WithRange("")}
					if cfg.Step2.LocalRead {
						opts = append(opts, clientv3.WithSerializable())
					}
					requests <- request{etcdv3Op: clientv3.OpGet(key, opts...)}

				case "zk", "zookeeper":
					// serializable read by default
					requests <- request{zkOp: zkOp{key: key}}

				case "consul":
					// serializable read by default
					requests <- request{consulOp: consulOp{key: key}}
				}
				if cfg.Step2.RequestIntervalMs > 0 {
					time.Sleep(time.Duration(cfg.Step2.RequestIntervalMs) * time.Millisecond)
				}
			}
			close(requests)
		}()

		wg.Wait()

		bar.Finish()

		close(results)
		<-pdoneC
	}

	return nil
}

func step3(cfg Config) error {
	req := agent.Request{}
	req.Operation = agent.Request_Stop

	switch cfg.Database {
	case "etcdv2":
		req.Database = agent.Request_etcdv2

	case "etcdv3":
		req.Database = agent.Request_etcdv3

	case "zk":
		cfg.Database = "zookeeper"
		req.Database = agent.Request_ZooKeeper

	case "zookeeper":
		req.Database = agent.Request_ZooKeeper

	case "consul":
		req.Database = agent.Request_Consul
	}

	donec, errc := make(chan struct{}), make(chan error)
	for i := range cfg.PeerIPs {

		go func(i int) {
			ep := cfg.AgentEndpoints[i]

			logger.Info("sending message",
				zap.Int("index", i),
				zap.String("operation", req.Operation.String()),
				zap.String("database", req.Database.String()),
				zap.String("endpoint", ep),
			)

			conn, err := grpc.Dial(ep, grpc.WithInsecure())
			if err != nil {
				logger.Error("grpc.Dial connecting error",
					zap.Int("index", i),
					zap.String("endpoint", ep),
					zap.Err(err),
				)
				errc <- err
				return
			}

			defer conn.Close()

			cli := agent.NewTransporterClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Consul takes longer
			resp, err := cli.Transfer(ctx, &req)
			cancel()
			if err != nil {
				logger.Error("cli.Transfer error",
					zap.Int("index", i),
					zap.String("endpoint", ep),
					zap.Err(err),
				)
				errc <- err
				return
			}

			logger.Info("response",
				zap.Int("index", i),
				zap.String("endpoint", ep),
				zap.String("response", fmt.Sprintf("%+v", resp)),
			)
			donec <- struct{}{}
		}(i)

		time.Sleep(time.Second)
	}

	cnt := 0
	for cnt != len(cfg.PeerIPs) {
		select {
		case <-donec:
		case err := <-errc:
			return err
		}
		cnt++
	}

	return nil
}
