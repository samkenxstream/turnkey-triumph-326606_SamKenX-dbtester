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
	"math"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/coreos/dbtester/agent/agentpb"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/report"
	consulapi "github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
)

func step2StressDatabase(cfg Config) error {
	vals, err := newValues(cfg)
	if err != nil {
		return err
	}

	switch cfg.Step2.BenchType {
	case "write":
		plog.Println("write generateReport is started...")

		// fixed number of client numbers
		if len(cfg.Step2.ConnectionsClients) == 0 {
			h, done := newWriteHandlers(cfg)
			reqGen := func(inflightReqs chan<- request) { generateWrites(cfg, 0, vals, inflightReqs) }
			generateReport(cfg, h, done, reqGen)

		} else {
			// variable client numbers
			rs := assignRequest(cfg.Step2.ConnectionsClients, cfg.Step2.TotalRequests)

			var stats []report.Stats
			reqCompleted := 0
			for i := 0; i < len(rs); i++ {
				copied := cfg
				copied.Step2.Connections = cfg.Step2.ConnectionsClients[i]
				copied.Step2.Clients = cfg.Step2.ConnectionsClients[i]
				copied.Step2.TotalRequests = rs[i]

				go func() {
					plog.Infof("signaling agent with client number %d", copied.Step2.Clients)
					if _, err := bcastReq(copied, agentpb.Request_Heartbeat); err != nil {
						plog.Panic(err)
					}
				}()

				h, done := newWriteHandlers(copied)
				reqGen := func(inflightReqs chan<- request) { generateWrites(copied, reqCompleted, vals, inflightReqs) }
				b := newBenchmark(copied.Step2.TotalRequests, copied.Step2.Clients, h, done, reqGen)

				// wait until rs[i] requests are finished
				// do not end reports yet
				b.startRequests()
				b.waitRequestsEnd()

				plog.Print("finishing reports...")
				now := time.Now()
				b.finishReports()
				plog.Printf("finished reports... took %v", time.Since(now))

				reqCompleted += rs[i]
				stats = append(stats, b.stats)
			}
			plog.Info("combining all reports")

			tsToClientN := make(map[int64]int, cfg.Step2.TotalRequests)
			combined := report.Stats{ErrorDist: make(map[string]int)}
			for i, st := range stats {
				combined.AvgTotal += st.AvgTotal
				combined.Total += st.Total
				combined.Lats = append(combined.Lats, st.Lats...)
				combined.TimeSeries = append(combined.TimeSeries, st.TimeSeries...)

				clientsN := cfg.Step2.ConnectionsClients[i]
				for _, v := range st.TimeSeries {
					tsToClientN[v.Timestamp] = clientsN
				}

				for k, v := range st.ErrorDist {
					if _, ok := combined.ErrorDist[k]; !ok {
						combined.ErrorDist[k] = v
					} else {
						combined.ErrorDist[k] += v
					}
				}
			}

			combined.Average = combined.AvgTotal / float64(len(combined.Lats))
			combined.RPS = float64(len(combined.Lats)) / combined.Total.Seconds()
			plog.Printf("got total %d data points and total %f seconds (RPS %f)", len(combined.Lats), combined.Total.Seconds(), combined.RPS)

			for i := range combined.Lats {
				dev := combined.Lats[i] - combined.Average
				combined.Stddev += dev * dev
			}
			combined.Stddev = math.Sqrt(combined.Stddev / float64(len(combined.Lats)))

			sort.Float64s(combined.Lats)
			if len(combined.Lats) > 0 {
				combined.Fastest = combined.Lats[0]
				combined.Slowest = combined.Lats[len(combined.Lats)-1]
			}

			plog.Info("combined all reports")
			printStats(combined)
			saveAllStats(cfg, combined, tsToClientN)
		}

		plog.Println("write generateReport is finished...")

		plog.Println("checking total keys on", cfg.DatabaseEndpoints)
		var totalKeysFunc func([]string) map[string]int64
		switch cfg.Database {
		case "etcdv2":
			totalKeysFunc = getTotalKeysEtcdv2
		case "etcdv3":
			totalKeysFunc = getTotalKeysEtcdv3
		case "zookeeper", "zetcd":
			totalKeysFunc = getTotalKeysZk
		case "consul", "cetcd":
			totalKeysFunc = getTotalKeysConsul
		}
		for k, v := range totalKeysFunc(cfg.DatabaseEndpoints) {
			plog.Infof("expected write total results [expected_total: %d | database: %q | endpoint: %q | number_of_keys: %d]",
				cfg.Step2.TotalRequests, cfg.Database, k, v)
		}

	case "read":
		key, value := sameKey(cfg.Step2.KeySize), vals.strings[0]

		switch cfg.Database {
		case "etcdv2":
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, cfg.Database)
			var err error
			for i := 0; i < 7; i++ {
				clients := mustCreateClientsEtcdv2(cfg.DatabaseEndpoints, cfg.Step2.Connections)
				_, err = clients[0].Set(context.Background(), key, value, nil)
				if err != nil {
					continue
				}
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, cfg.Database)
				break
			}
			if err != nil {
				plog.Errorf("write error [request: PUT | key: %q | database: %q]", key, cfg.Database)
				os.Exit(1)
			}

		case "etcdv3":
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, cfg.Database)
			var err error
			for i := 0; i < 7; i++ {
				clients := mustCreateClientsEtcdv3(cfg.DatabaseEndpoints, etcdv3ClientCfg{
					totalConns:   1,
					totalClients: 1,
				})
				_, err = clients[0].Do(context.Background(), clientv3.OpPut(key, value))
				if err != nil {
					continue
				}
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, cfg.Database)
				break
			}
			if err != nil {
				plog.Errorf("write error [request: PUT | key: %q | database: %q]", key, cfg.Database)
				os.Exit(1)
			}

		case "zookeeper", "zetcd":
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, cfg.Database)
			var err error
			for i := 0; i < 7; i++ {
				conns := mustCreateConnsZk(cfg.DatabaseEndpoints, cfg.Step2.Connections)
				_, err = conns[0].Create("/"+key, vals.bytes[0], zkCreateFlags, zkCreateAcl)
				if err != nil {
					continue
				}
				for j := range conns {
					conns[j].Close()
				}
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, cfg.Database)
				break
			}
			if err != nil {
				plog.Errorf("write error [request: PUT | key: %q | database: %q]", key, cfg.Database)
				os.Exit(1)
			}

		case "consul", "cetcd":
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, cfg.Database)
			var err error
			for i := 0; i < 7; i++ {
				clients := mustCreateConnsConsul(cfg.DatabaseEndpoints, cfg.Step2.Connections)
				_, err = clients[0].Put(&consulapi.KVPair{Key: key, Value: vals.bytes[0]}, nil)
				if err != nil {
					continue
				}
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, cfg.Database)
				break
			}
			if err != nil {
				plog.Errorf("write done [request: PUT | key: %q | database: %q]", key, cfg.Database)
				os.Exit(1)
			}
		}

		h, done := newReadHandlers(cfg)
		reqGen := func(inflightReqs chan<- request) { generateReads(cfg, key, inflightReqs) }
		generateReport(cfg, h, done, reqGen)
		plog.Println("read generateReport is finished...")

	case "read-oneshot":
		key, value := sameKey(cfg.Step2.KeySize), vals.strings[0]
		plog.Infof("writing key for read-oneshot [key: %q | database: %q]", key, cfg.Database)
		var err error
		switch cfg.Database {
		case "etcdv2":
			clients := mustCreateClientsEtcdv2(cfg.DatabaseEndpoints, 1)
			_, err = clients[0].Set(context.Background(), key, value, nil)

		case "etcdv3":
			clients := mustCreateClientsEtcdv3(cfg.DatabaseEndpoints, etcdv3ClientCfg{
				totalConns:   1,
				totalClients: 1,
			})
			_, err = clients[0].Do(context.Background(), clientv3.OpPut(key, value))
			clients[0].Close()

		case "zookeeper", "zetcd":
			conns := mustCreateConnsZk(cfg.DatabaseEndpoints, 1)
			_, err = conns[0].Create("/"+key, vals.bytes[0], zkCreateFlags, zkCreateAcl)
			conns[0].Close()

		case "consul", "cetcd":
			clients := mustCreateConnsConsul(cfg.DatabaseEndpoints, 1)
			_, err = clients[0].Put(&consulapi.KVPair{Key: key, Value: vals.bytes[0]}, nil)
		}
		if err != nil {
			plog.Errorf("write error on read-oneshot (%v)", err)
			os.Exit(1)
		}

		h := newReadOneshotHandlers(cfg)
		reqGen := func(inflightReqs chan<- request) { generateReads(cfg, key, inflightReqs) }
		generateReport(cfg, h, nil, reqGen)
		plog.Println("read-oneshot generateReport is finished...")
	}

	return nil
}

func newReadHandlers(cfg Config) (rhs []ReqHandler, done func()) {
	rhs = make([]ReqHandler, cfg.Step2.Clients)
	switch cfg.Database {
	case "etcdv2":
		conns := mustCreateClientsEtcdv2(cfg.DatabaseEndpoints, cfg.Step2.Connections)
		for i := range conns {
			rhs[i] = newGetEtcd2(conns[i])
		}
	case "etcdv3":
		clients := mustCreateClientsEtcdv3(cfg.DatabaseEndpoints, etcdv3ClientCfg{
			totalConns:   cfg.Step2.Connections,
			totalClients: cfg.Step2.Clients,
		})
		for i := range clients {
			rhs[i] = newGetEtcd3(clients[i].KV)
		}
		done = func() {
			for i := range clients {
				clients[i].Close()
			}
		}
	case "zookeeper", "zetcd":
		conns := mustCreateConnsZk(cfg.DatabaseEndpoints, cfg.Step2.Connections)
		for i := range conns {
			rhs[i] = newGetZK(conns[i])
		}
		done = func() {
			for i := range conns {
				conns[i].Close()
			}
		}
	case "consul", "cetcd":
		conns := mustCreateConnsConsul(cfg.DatabaseEndpoints, cfg.Step2.Connections)
		for i := range conns {
			rhs[i] = newGetConsul(conns[i])
		}
	}
	return rhs, done
}

func newWriteHandlers(cfg Config) (rhs []ReqHandler, done func()) {
	rhs = make([]ReqHandler, cfg.Step2.Clients)
	switch cfg.Database {
	case "etcdv2":
		conns := mustCreateClientsEtcdv2(cfg.DatabaseEndpoints, cfg.Step2.Connections)
		for i := range conns {
			rhs[i] = newPutEtcd2(conns[i])
		}
	case "etcdv3":
		etcdClients := mustCreateClientsEtcdv3(cfg.DatabaseEndpoints, etcdv3ClientCfg{
			totalConns:   cfg.Step2.Connections,
			totalClients: cfg.Step2.Clients,
		})
		for i := range etcdClients {
			rhs[i] = newPutEtcd3(etcdClients[i])
		}
		done = func() {
			for i := range etcdClients {
				etcdClients[i].Close()
			}
		}
	case "zookeeper", "zetcd":
		if cfg.Step2.SameKey {
			key := sameKey(cfg.Step2.KeySize)
			valueBts := randBytes(cfg.Step2.ValueSize)
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, cfg.Database)
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
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, cfg.Database)
				break
			}
			if err != nil {
				plog.Errorf("write error [request: PUT | key: %q | database: %q]", key, cfg.Database)
				os.Exit(1)
			}
		}

		conns := mustCreateConnsZk(cfg.DatabaseEndpoints, cfg.Step2.Connections)
		for i := range conns {
			if cfg.Step2.SameKey {
				rhs[i] = newPutOverwriteZK(conns[i])
			} else {
				rhs[i] = newPutCreateZK(conns[i])
			}
		}
		done = func() {
			for i := range conns {
				conns[i].Close()
			}
		}
	case "consul", "cetcd":
		conns := mustCreateConnsConsul(cfg.DatabaseEndpoints, cfg.Step2.Connections)
		for i := range conns {
			rhs[i] = newPutConsul(conns[i])
		}
	}

	for k := range rhs {
		if rhs[k] == nil {
			plog.Panicf("%d-th write handler is nil (out of %d)", k, len(rhs))
		}
	}
	return
}

func newReadOneshotHandlers(cfg Config) []ReqHandler {
	rhs := make([]ReqHandler, cfg.Step2.Clients)
	switch cfg.Database {
	case "etcdv2":
		for i := range rhs {
			rhs[i] = func(ctx context.Context, req *request) error {
				conns := mustCreateClientsEtcdv2(cfg.DatabaseEndpoints, 1)
				return newGetEtcd2(conns[0])(ctx, req)
			}
		}
	case "etcdv3":
		for i := range rhs {
			rhs[i] = func(ctx context.Context, req *request) error {
				conns := mustCreateClientsEtcdv3(cfg.DatabaseEndpoints, etcdv3ClientCfg{
					totalConns:   1,
					totalClients: 1,
				})
				defer conns[0].Close()
				return newGetEtcd3(conns[0])(ctx, req)
			}
		}
	case "zookeeper", "zetcd":
		for i := range rhs {
			rhs[i] = func(ctx context.Context, req *request) error {
				conns := mustCreateConnsZk(cfg.DatabaseEndpoints, cfg.Step2.Connections)
				defer conns[0].Close()
				return newGetZK(conns[0])(ctx, req)
			}
		}
	case "consul", "cetcd":
		for i := range rhs {
			rhs[i] = func(ctx context.Context, req *request) error {
				conns := mustCreateConnsConsul(cfg.DatabaseEndpoints, 1)
				return newGetConsul(conns[0])(ctx, req)
			}
		}
	}
	return rhs
}

func generateReads(cfg Config, key string, inflightReqs chan<- request) {
	defer close(inflightReqs)

	var rateLimiter *rate.Limiter
	if cfg.Step2.RequestsPerSecond > 0 {
		rateLimiter = rate.NewLimiter(rate.Limit(cfg.Step2.RequestsPerSecond), cfg.Step2.RequestsPerSecond)
	}

	for i := 0; i < cfg.Step2.TotalRequests; i++ {
		if rateLimiter != nil {
			rateLimiter.Wait(context.TODO())
		}

		switch cfg.Database {
		case "etcdv2":
			// serializable read by default
			inflightReqs <- request{etcdv2Op: etcdv2Op{key: key}}

		case "etcdv3":
			opts := []clientv3.OpOption{clientv3.WithRange("")}
			if cfg.Step2.StaleRead {
				opts = append(opts, clientv3.WithSerializable())
			}
			inflightReqs <- request{etcdv3Op: clientv3.OpGet(key, opts...)}

		case "zookeeper", "zetcd":
			op := zkOp{key: key}
			if cfg.Step2.StaleRead {
				op.staleRead = true
			}
			inflightReqs <- request{zkOp: op}

		case "consul", "cetcd":
			op := consulOp{key: key}
			if cfg.Step2.StaleRead {
				op.staleRead = true
			}
			inflightReqs <- request{consulOp: op}
		}
	}
}

func generateWrites(cfg Config, startIdx int, vals values, inflightReqs chan<- request) {
	var rateLimiter *rate.Limiter
	if cfg.Step2.RequestsPerSecond > 0 {
		rateLimiter = rate.NewLimiter(rate.Limit(cfg.Step2.RequestsPerSecond), cfg.Step2.RequestsPerSecond)
	}

	var wg sync.WaitGroup
	defer func() {
		close(inflightReqs)
		wg.Wait()
	}()

	for i := 0; i < cfg.Step2.TotalRequests; i++ {
		k := sequentialKey(cfg.Step2.KeySize, i+startIdx)
		if cfg.Step2.SameKey {
			k = sameKey(cfg.Step2.KeySize)
		}

		v := vals.bytes[i%vals.sampleSize]
		vs := vals.strings[i%vals.sampleSize]

		if rateLimiter != nil {
			rateLimiter.Wait(context.TODO())
		}

		switch cfg.Database {
		case "etcdv2":
			inflightReqs <- request{etcdv2Op: etcdv2Op{key: k, value: vs}}
		case "etcdv3":
			inflightReqs <- request{etcdv3Op: clientv3.OpPut(k, vs)}
		case "zookeeper", "zetcd":
			inflightReqs <- request{zkOp: zkOp{key: "/" + k, value: v}}
		case "consul", "cetcd":
			inflightReqs <- request{consulOp: consulOp{key: k, value: v}}
		}
	}
}
