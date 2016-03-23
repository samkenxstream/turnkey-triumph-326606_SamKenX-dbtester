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

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
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

	keySpaceSize int
	seqKeys      bool

	etcdCompactionCycle int64
)

func init() {
	Command.AddCommand(putCmd)
	putCmd.Flags().IntVar(&keySize, "key-size", 8, "Key size of put request")
	putCmd.Flags().IntVar(&valSize, "val-size", 8, "Value size of put request")
	putCmd.Flags().IntVar(&putTotal, "total", 10000, "Total number of put requests")
	putCmd.Flags().IntVar(&keySpaceSize, "key-space-size", 1, "Maximum possible keys")
	putCmd.Flags().BoolVar(&seqKeys, "sequential-keys", false, "Use sequential keys")
	putCmd.Flags().Int64Var(&etcdCompactionCycle, "etcd-compaction-cycle", 0, "Compact every X number of put requests. 0 means no compaction.")
}

type request struct {
	etcdOp   clientv3.Op
	etcd2Op  etcd2Op
	zkOp     zkOp
	consulOp consulOp
}

type etcd2Op struct {
	key   string
	value string
}

type zkOp struct {
	key   string
	value []byte
}

type consulOp struct {
	key   string
	value []byte
}

func putFunc(cmd *cobra.Command, args []string) {
	if keySpaceSize <= 0 {
		fmt.Fprintf(os.Stderr, "expected positive --key-space-size, got (%v)", keySpaceSize)
		os.Exit(1)
	}

	results = make(chan result)
	requests := make(chan request, totalClients)
	bar = pb.New(putTotal)

	k, v := make([]byte, keySize), string(mustRandBytes(valSize))
	keys := multiRandBytes(keySize, putTotal)

	bar.Format("Bom !")
	bar.Start()

	var etcdClients []*clientv3.Client
	switch database {
	case "etcd":
		etcdClients = mustCreateClients(totalClients, totalConns)
		for i := range etcdClients {
			wg.Add(1)
			go doPut(context.Background(), etcdClients[i], requests)
		}
		defer func() {
			for i := range etcdClients {
				etcdClients[i].Close()
			}
		}()

	case "etcd2":
		conns := mustCreateClientsEtcd2(totalConns)
		for i := range conns {
			wg.Add(1)
			go doPutEtcd2(context.Background(), conns[i], requests)
		}

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
			if seqKeys {
				binary.PutVarint(k, int64(i%keySpaceSize))
			} else {
				k = keys[i]
			}
			switch database {
			case "etcd":
				requests <- request{etcdOp: clientv3.OpPut(string(k), v)}
			case "zk":
				requests <- request{zkOp: zkOp{key: "/" + string(k), value: []byte(v)}}
			}
		}
		close(requests)
	}()

	wg.Wait()

	bar.Finish()

	close(results)
	<-pdoneC
}

func doPut(ctx context.Context, client clientv3.KV, requests <-chan request) {
	defer wg.Done()

	for req := range requests {
		op := req.etcdOp
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

func doPutEtcd2(ctx context.Context, conn clientv2.KeysAPI, requests <-chan request) {
	defer wg.Done()

	for req := range requests {
		op := req.etcd2Op
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
