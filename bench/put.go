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

package bench

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	clientv2 "github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Benchmark put",

	Run: putFunc,
}

var (
	zkCreateFlags = int32(0)
	zkCreateAcl   = zk.WorldACL(zk.PermAll)

	keySize int
	valSize int

	putTotal int

	etcdCompactionCycle int64
)

func init() {
	Command.AddCommand(putCmd)
	putCmd.Flags().IntVar(&keySize, "key-size", 8, "Key size of put request")
	putCmd.Flags().IntVar(&valSize, "val-size", 8, "Value size of put request")
	putCmd.Flags().IntVar(&putTotal, "total", 10000, "Total number of put requests")
	putCmd.Flags().Int64Var(&etcdCompactionCycle, "etcd-compaction-cycle", 0, "Compact every X number of put requests. 0 means no compaction.")
}

func putFunc(cmd *cobra.Command, args []string) {
	results = make(chan result)
	requests := make(chan request, totalClients)
	bar = pb.New(putTotal)

	keys := multiRandStrings(keySize, putTotal)
	value := string(mustRandBytes(valSize))

	bar.Format("Bom !")
	bar.Start()

	var etcdClients []*clientv3.Client
	switch database {
	case "etcdv2":
		conns := mustCreateClientsEtcdv2(totalConns)
		for i := range conns {
			wg.Add(1)
			go doPutEtcdv2(context.Background(), conns[i], requests)
		}

	case "etcdv3":
		etcdClients = mustCreateClientsEtcdv3(totalClients, totalConns)
		for i := range etcdClients {
			wg.Add(1)
			go doPutEtcdv3(context.Background(), etcdClients[i], requests)
		}
		defer func() {
			for i := range etcdClients {
				etcdClients[i].Close()
			}
		}()

	case "zk":
		conns := mustCreateConnsZk(totalConns)
		defer func() {
			for i := range conns {
				conns[i].Close()
			}
		}()
		for i := range conns {
			wg.Add(1)
			go doPutZk(conns[i], requests)
		}

	case "consul":
		conns := mustCreateConnsConsul(totalConns)
		for i := range conns {
			wg.Add(1)
			go doPutConsul(conns[i], requests)
		}

	default:
		log.Fatalf("unknown database %s", database)
	}

	pdoneC := printReport(results)
	go func() {
		for i := 0; i < putTotal; i++ {
			if database == "etcd" && etcdCompactionCycle > 0 && int64(i)%etcdCompactionCycle == 0 {
				log.Printf("etcd starting compaction at %d put request", i)
				go func() {
					compactKV(etcdClients)
				}()
			}
			key := keys[i]
			switch database {
			case "etcdv2":
				requests <- request{etcdv2Op: etcdv2Op{key: key, value: value}}
			case "etcdv3":
				requests <- request{etcdv3Op: clientv3.OpPut(key, value)}
			case "zk":
				requests <- request{zkOp: zkOp{key: "/" + key, value: []byte(value)}}
			case "consul":
				requests <- request{consulOp: consulOp{key: key, value: []byte(value)}}
			}
		}
		close(requests)
	}()

	wg.Wait()

	bar.Finish()

	close(results)
	<-pdoneC

	fmt.Println("Expected Put Total:", putTotal)
	switch database {
	case "etcdv2":
		for k, v := range getTotalKeysEtcdv2(endpoints) {
			fmt.Println("Endpoint      :", k)
			fmt.Println("Number of Keys:", v)
			fmt.Println()
		}
	case "etcdv3":
		for k, v := range getTotalKeysEtcdv3(endpoints) {
			fmt.Println("Endpoint      :", k)
			fmt.Println("Number of Keys:", v)
			fmt.Println()
		}

	case "zk":
		for k, v := range getTotalKeysZk(endpoints) {
			fmt.Println("Endpoint      :", k)
			fmt.Println("Number of Keys:", v)
			fmt.Println()
		}

	case "consul":
		for k, v := range getTotalKeysConsul(endpoints) {
			fmt.Println("Endpoint      :", k)
			fmt.Println("Number of Keys:", v)
			fmt.Println()
		}
	}
}

func doPutEtcdv2(ctx context.Context, conn clientv2.KeysAPI, requests <-chan request) {
	defer wg.Done()

	for req := range requests {
		op := req.etcdv2Op
		st := time.Now()

		_, err := conn.Set(context.Background(), op.key, op.value, nil)

		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		results <- result{errStr: errStr, duration: time.Since(st), happened: time.Now()}
		bar.Increment()
	}
}

func getTotalKeysEtcdv2(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	for _, ep := range endpoints {
		rs[ep] = 0 // not supported in metrics
	}
	return rs
}

func doPutEtcdv3(ctx context.Context, client clientv3.KV, requests <-chan request) {
	defer wg.Done()

	for req := range requests {
		op := req.etcdv3Op
		st := time.Now()
		_, err := client.Do(ctx, op)

		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		results <- result{errStr: errStr, duration: time.Since(st), happened: time.Now()}
		bar.Increment()
	}
}

func getTotalKeysEtcdv3(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	for _, ep := range endpoints {
		if strings.HasPrefix(ep, "http://") {
			ep = "http://" + ep
		}
		resp, err := http.Get(ep + "/metrics")
		if err != nil {
			log.Println(err)
			rs[ep] = 0
		}
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.HasPrefix(txt, "#") {
				continue
			}
			ts := strings.SplitN(txt, " ", 2)
			fv := 0.0
			if len(ts) == 2 {
				v, err := strconv.ParseFloat(ts[1], 64)
				if err == nil {
					fv = v
				}
			}
			if ts[0] == "etcd_storage_keys_total" {
				rs[ep] = int64(fv)
				break
			}
		}
		gracefulClose(resp)
	}
	return rs
}

func doPutZk(conn *zk.Conn, requests <-chan request) {
	defer wg.Done()

	for req := range requests {
		op := req.zkOp
		st := time.Now()

		_, err := conn.Create(op.key, op.value, zkCreateFlags, zkCreateAcl)

		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		results <- result{errStr: errStr, duration: time.Since(st), happened: time.Now()}
		bar.Increment()
	}
}

func getTotalKeysZk(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	stats, ok := zk.FLWSrvr(endpoints, 5*time.Second)
	if !ok {
		log.Printf("getTotalKeysZk failed with %+v", stats)
		for _, ep := range endpoints {
			rs[ep] = 0
		}
		return rs
	}
	for i, s := range stats {
		rs[endpoints[i]] = s.NodeCount
	}
	return rs
}

func doPutConsul(conn *consulapi.KV, requests <-chan request) {
	defer wg.Done()

	for req := range requests {
		op := req.consulOp
		st := time.Now()

		_, err := conn.Put(&consulapi.KVPair{Key: op.key, Value: op.value}, nil)

		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		results <- result{errStr: errStr, duration: time.Since(st), happened: time.Now()}
		bar.Increment()
	}
}

func getTotalKeysConsul(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	for _, ep := range endpoints {
		rs[ep] = 0 // not supported in consul
	}
	return rs
}

func compactKV(clients []*clientv3.Client) {
	var curRev int64
	for _, c := range clients {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := c.KV.Get(ctx, "foo")
		cancel()
		if err != nil {
			panic(err)
		}
		curRev = resp.Header.Revision
		break
	}
	revToCompact := max(0, curRev-1000)
	for _, c := range clients {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := c.KV.Compact(ctx, revToCompact)
		cancel()
		if err != nil {
			panic(err)
		}
		break
	}
}

func max(n1, n2 int64) int64 {
	if n1 > n2 {
		return n1
	}
	return n2
}
