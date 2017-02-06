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
	"math"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/coreos/dbtester/dbtesterpb"
	"github.com/coreos/dbtester/pkg/report"
	"github.com/coreos/etcd/clientv3"
	consulapi "github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
)

// Stress stresses the database.
func (cfg *Config) Stress(databaseID string) error {
	gcfg, ok := cfg.DatabaseIDToTestGroup[databaseID]
	if !ok {
		return fmt.Errorf("%q does not exist", databaseID)
	}

	vals, err := newValues(gcfg)
	if err != nil {
		return err
	}

	switch gcfg.BenchmarkOptions.Type {
	case "write":
		plog.Println("write generateReport is started...")

		// fixed number of client numbers
		if len(gcfg.BenchmarkOptions.ConnectionClientNumbers) == 0 {
			h, done := newWriteHandlers(gcfg)
			reqGen := func(inflightReqs chan<- request) { generateWrites(gcfg, 0, vals, inflightReqs) }
			cfg.generateReport(gcfg, h, done, reqGen)

		} else {
			// variable client numbers
			rs := assignRequest(gcfg.BenchmarkOptions.ConnectionClientNumbers, gcfg.BenchmarkOptions.RequestNumber)

			var stats []report.Stats
			reqCompleted := int64(0)
			for i := 0; i < len(rs); i++ {
				copied := gcfg
				copied.BenchmarkOptions.ConnectionNumber = gcfg.BenchmarkOptions.ConnectionClientNumbers[i]
				copied.BenchmarkOptions.ClientNumber = gcfg.BenchmarkOptions.ConnectionClientNumbers[i]
				copied.BenchmarkOptions.RequestNumber = rs[i]
				ncfg := *cfg
				ncfg.DatabaseIDToTestGroup[databaseID] = copied

				go func() {
					plog.Infof("signaling agent with client number %d", copied.BenchmarkOptions.ClientNumber)
					if _, err := (&ncfg).BroadcaseRequest(databaseID, dbtesterpb.Request_Heartbeat); err != nil {
						plog.Panic(err)
					}
				}()

				h, done := newWriteHandlers(copied)
				reqGen := func(inflightReqs chan<- request) { generateWrites(copied, reqCompleted, vals, inflightReqs) }
				b := newBenchmark(copied.BenchmarkOptions.RequestNumber, copied.BenchmarkOptions.ClientNumber, h, done, reqGen)

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

			tsToClientN := make(map[int64]int64, gcfg.BenchmarkOptions.RequestNumber)
			combined := report.Stats{ErrorDist: make(map[string]int)}
			for i, st := range stats {
				combined.AvgTotal += st.AvgTotal
				combined.Total += st.Total
				combined.Lats = append(combined.Lats, st.Lats...)
				combined.TimeSeries = append(combined.TimeSeries, st.TimeSeries...)
				//
				// Need to handle duplicate unix second timestamps when two ranges are merged.
				// This can happen when the following run happens within the same unix timesecond,
				// since finishing up the previous report and restarting the next range of requests
				// with different number of clients takes only 100+/- ms.
				//
				// For instance, we have the following raw data:
				//
				//   unix-second, client-number, throughput
				//   1486389257,       700,         30335  === ending of previous combined.TimeSeries
				//   1486389258,      "700",        23188  === ending of previous combined.TimeSeries
				//   1486389258,       1000,         5739  === beginning of current st.TimeSeries
				//
				// And the line below will overwrite the 'client-number' as:
				//
				//   unix-second, client-number, throughput
				//   1486389257,       700,        30335  === ending of previous combined.TimeSeries
				//   1486389258,      "1000",      23188  === ending of previous combined.TimeSeries
				//   1486389258,       1000,        5739  === beginning of current st.TimeSeries
				//
				// So now we have two duplicate unix time seconds.
				//
				clientsN := gcfg.BenchmarkOptions.ConnectionClientNumbers[i]
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

			// handle duplicate unix seconds around boundaries
			sec2dp := make(map[int64]report.DataPoint)
			for _, tss := range combined.TimeSeries {
				v, ok := sec2dp[tss.Timestamp]
				if !ok {
					sec2dp[tss.Timestamp] = tss
				}

				// two datapoints share the time unix second
				if v.MinLatency > tss.MinLatency {
					v.MinLatency = tss.MinLatency
				}
				if v.MaxLatency < tss.MaxLatency {
					v.MaxLatency = tss.MaxLatency
				}
				v.AvgLatency = (v.AvgLatency + tss.AvgLatency) / time.Duration(2)
				v.ThroughPut += tss.ThroughPut
				sec2dp[tss.Timestamp] = v
			}
			var fts report.TimeSeries
			for _, dp := range sec2dp {
				fts = append(fts, dp)
			}
			sort.Sort(report.TimeSeries(fts))
			combined.TimeSeries = fts

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
			cfg.saveAllStats(gcfg, combined, tsToClientN)
		}

		plog.Println("write generateReport is finished...")

		plog.Println("checking total keys on", gcfg.DatabaseEndpoints)
		var totalKeysFunc func([]string) map[string]int64
		switch gcfg.DatabaseID {
		case "etcdv2":
			totalKeysFunc = getTotalKeysEtcdv2
		case "etcdv3":
			totalKeysFunc = getTotalKeysEtcdv3
		case "zookeeper", "zetcd":
			totalKeysFunc = getTotalKeysZk
		case "consul", "cetcd":
			totalKeysFunc = getTotalKeysConsul
		}
		for k, v := range totalKeysFunc(gcfg.DatabaseEndpoints) {
			plog.Infof("expected write total results [expected_total: %d | database: %q | endpoint: %q | number_of_keys: %d]",
				gcfg.BenchmarkOptions.RequestNumber, gcfg.DatabaseID, k, v)
		}

	case "read":
		key, value := sameKey(gcfg.BenchmarkOptions.KeySizeBytes), vals.strings[0]

		switch gcfg.DatabaseID {
		case "etcdv2":
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
			var err error
			for i := 0; i < 7; i++ {
				clients := mustCreateClientsEtcdv2(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
				_, err = clients[0].Set(context.Background(), key, value, nil)
				if err != nil {
					continue
				}
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				break
			}
			if err != nil {
				plog.Errorf("write error [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				os.Exit(1)
			}

		case "etcdv3":
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
			var err error
			for i := 0; i < 7; i++ {
				clients := mustCreateClientsEtcdv3(gcfg.DatabaseEndpoints, etcdv3ClientCfg{
					totalConns:   1,
					totalClients: 1,
				})
				_, err = clients[0].Do(context.Background(), clientv3.OpPut(key, value))
				if err != nil {
					continue
				}
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				break
			}
			if err != nil {
				plog.Errorf("write error [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				os.Exit(1)
			}

		case "zookeeper", "zetcd":
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
			var err error
			for i := 0; i < 7; i++ {
				conns := mustCreateConnsZk(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
				_, err = conns[0].Create("/"+key, vals.bytes[0], zkCreateFlags, zkCreateACL)
				if err != nil {
					continue
				}
				for j := range conns {
					conns[j].Close()
				}
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				break
			}
			if err != nil {
				plog.Errorf("write error [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				os.Exit(1)
			}

		case "consul", "cetcd":
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
			var err error
			for i := 0; i < 7; i++ {
				clients := mustCreateConnsConsul(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
				_, err = clients[0].Put(&consulapi.KVPair{Key: key, Value: vals.bytes[0]}, nil)
				if err != nil {
					continue
				}
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				break
			}
			if err != nil {
				plog.Errorf("write done [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				os.Exit(1)
			}
		}

		h, done := newReadHandlers(gcfg)
		reqGen := func(inflightReqs chan<- request) { generateReads(gcfg, key, inflightReqs) }
		cfg.generateReport(gcfg, h, done, reqGen)
		plog.Println("read generateReport is finished...")

	case "read-oneshot":
		key, value := sameKey(gcfg.BenchmarkOptions.KeySizeBytes), vals.strings[0]
		plog.Infof("writing key for read-oneshot [key: %q | database: %q]", key, gcfg.DatabaseID)
		var err error
		switch gcfg.DatabaseID {
		case "etcdv2":
			clients := mustCreateClientsEtcdv2(gcfg.DatabaseEndpoints, 1)
			_, err = clients[0].Set(context.Background(), key, value, nil)

		case "etcdv3":
			clients := mustCreateClientsEtcdv3(gcfg.DatabaseEndpoints, etcdv3ClientCfg{
				totalConns:   1,
				totalClients: 1,
			})
			_, err = clients[0].Do(context.Background(), clientv3.OpPut(key, value))
			clients[0].Close()

		case "zookeeper", "zetcd":
			conns := mustCreateConnsZk(gcfg.DatabaseEndpoints, 1)
			_, err = conns[0].Create("/"+key, vals.bytes[0], zkCreateFlags, zkCreateACL)
			conns[0].Close()

		case "consul", "cetcd":
			clients := mustCreateConnsConsul(gcfg.DatabaseEndpoints, 1)
			_, err = clients[0].Put(&consulapi.KVPair{Key: key, Value: vals.bytes[0]}, nil)
		}
		if err != nil {
			plog.Errorf("write error on read-oneshot (%v)", err)
			os.Exit(1)
		}

		h := newReadOneshotHandlers(gcfg)
		reqGen := func(inflightReqs chan<- request) { generateReads(gcfg, key, inflightReqs) }
		cfg.generateReport(gcfg, h, nil, reqGen)
		plog.Println("read-oneshot generateReport is finished...")
	}

	return nil
}

func newReadHandlers(gcfg TestGroup) (rhs []ReqHandler, done func()) {
	rhs = make([]ReqHandler, gcfg.BenchmarkOptions.ClientNumber)
	switch gcfg.DatabaseID {
	case "etcdv2":
		conns := mustCreateClientsEtcdv2(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
		for i := range conns {
			rhs[i] = newGetEtcd2(conns[i])
		}
	case "etcdv3":
		clients := mustCreateClientsEtcdv3(gcfg.DatabaseEndpoints, etcdv3ClientCfg{
			totalConns:   gcfg.BenchmarkOptions.ConnectionNumber,
			totalClients: gcfg.BenchmarkOptions.ClientNumber,
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
		conns := mustCreateConnsZk(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
		for i := range conns {
			rhs[i] = newGetZK(conns[i])
		}
		done = func() {
			for i := range conns {
				conns[i].Close()
			}
		}
	case "consul", "cetcd":
		conns := mustCreateConnsConsul(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
		for i := range conns {
			rhs[i] = newGetConsul(conns[i])
		}
	}
	return rhs, done
}

func newWriteHandlers(gcfg TestGroup) (rhs []ReqHandler, done func()) {
	rhs = make([]ReqHandler, gcfg.BenchmarkOptions.ClientNumber)
	switch gcfg.DatabaseID {
	case "etcdv2":
		conns := mustCreateClientsEtcdv2(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
		for i := range conns {
			rhs[i] = newPutEtcd2(conns[i])
		}
	case "etcdv3":
		etcdClients := mustCreateClientsEtcdv3(gcfg.DatabaseEndpoints, etcdv3ClientCfg{
			totalConns:   gcfg.BenchmarkOptions.ConnectionNumber,
			totalClients: gcfg.BenchmarkOptions.ClientNumber,
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
		if gcfg.BenchmarkOptions.SameKey {
			key := sameKey(gcfg.BenchmarkOptions.KeySizeBytes)
			valueBts := randBytes(gcfg.BenchmarkOptions.ValueSizeBytes)
			plog.Infof("write started [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
			var err error
			for i := 0; i < 7; i++ {
				conns := mustCreateConnsZk(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
				_, err = conns[0].Create("/"+key, valueBts, zkCreateFlags, zkCreateACL)
				if err != nil {
					continue
				}
				for j := range conns {
					conns[j].Close()
				}
				plog.Infof("write done [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				break
			}
			if err != nil {
				plog.Errorf("write error [request: PUT | key: %q | database: %q]", key, gcfg.DatabaseID)
				os.Exit(1)
			}
		}

		conns := mustCreateConnsZk(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
		for i := range conns {
			if gcfg.BenchmarkOptions.SameKey {
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
		conns := mustCreateConnsConsul(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
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

func newReadOneshotHandlers(gcfg TestGroup) []ReqHandler {
	rhs := make([]ReqHandler, gcfg.BenchmarkOptions.ClientNumber)
	switch gcfg.DatabaseID {
	case "etcdv2":
		for i := range rhs {
			rhs[i] = func(ctx context.Context, req *request) error {
				conns := mustCreateClientsEtcdv2(gcfg.DatabaseEndpoints, 1)
				return newGetEtcd2(conns[0])(ctx, req)
			}
		}
	case "etcdv3":
		for i := range rhs {
			rhs[i] = func(ctx context.Context, req *request) error {
				conns := mustCreateClientsEtcdv3(gcfg.DatabaseEndpoints, etcdv3ClientCfg{
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
				conns := mustCreateConnsZk(gcfg.DatabaseEndpoints, gcfg.BenchmarkOptions.ConnectionNumber)
				defer conns[0].Close()
				return newGetZK(conns[0])(ctx, req)
			}
		}
	case "consul", "cetcd":
		for i := range rhs {
			rhs[i] = func(ctx context.Context, req *request) error {
				conns := mustCreateConnsConsul(gcfg.DatabaseEndpoints, 1)
				return newGetConsul(conns[0])(ctx, req)
			}
		}
	}
	return rhs
}

func generateReads(gcfg TestGroup, key string, inflightReqs chan<- request) {
	defer close(inflightReqs)

	var rateLimiter *rate.Limiter
	if gcfg.BenchmarkOptions.RateLimitRequestsPerSecond > 0 {
		rateLimiter = rate.NewLimiter(
			rate.Limit(gcfg.BenchmarkOptions.RateLimitRequestsPerSecond),
			int(gcfg.BenchmarkOptions.RateLimitRequestsPerSecond),
		)
	}

	for i := int64(0); i < gcfg.BenchmarkOptions.RequestNumber; i++ {
		if rateLimiter != nil {
			rateLimiter.Wait(context.TODO())
		}

		switch gcfg.DatabaseID {
		case "etcdv2":
			// serializable read by default
			inflightReqs <- request{etcdv2Op: etcdv2Op{key: key}}

		case "etcdv3":
			opts := []clientv3.OpOption{clientv3.WithRange("")}
			if gcfg.BenchmarkOptions.StaleRead {
				opts = append(opts, clientv3.WithSerializable())
			}
			inflightReqs <- request{etcdv3Op: clientv3.OpGet(key, opts...)}

		case "zookeeper", "zetcd":
			op := zkOp{key: key}
			if gcfg.BenchmarkOptions.StaleRead {
				op.staleRead = true
			}
			inflightReqs <- request{zkOp: op}

		case "consul", "cetcd":
			op := consulOp{key: key}
			if gcfg.BenchmarkOptions.StaleRead {
				op.staleRead = true
			}
			inflightReqs <- request{consulOp: op}
		}
	}
}

func generateWrites(gcfg TestGroup, startIdx int64, vals values, inflightReqs chan<- request) {
	var rateLimiter *rate.Limiter
	if gcfg.BenchmarkOptions.RateLimitRequestsPerSecond > 0 {
		rateLimiter = rate.NewLimiter(
			rate.Limit(gcfg.BenchmarkOptions.RateLimitRequestsPerSecond),
			int(gcfg.BenchmarkOptions.RateLimitRequestsPerSecond),
		)
	}

	var wg sync.WaitGroup
	defer func() {
		close(inflightReqs)
		wg.Wait()
	}()

	for i := int64(0); i < gcfg.BenchmarkOptions.RequestNumber; i++ {
		k := sequentialKey(gcfg.BenchmarkOptions.KeySizeBytes, i+startIdx)
		if gcfg.BenchmarkOptions.SameKey {
			k = sameKey(gcfg.BenchmarkOptions.KeySizeBytes)
		}

		v := vals.bytes[i%int64(vals.sampleSize)]
		vs := vals.strings[i%int64(vals.sampleSize)]

		if rateLimiter != nil {
			rateLimiter.Wait(context.TODO())
		}

		switch gcfg.DatabaseID {
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
