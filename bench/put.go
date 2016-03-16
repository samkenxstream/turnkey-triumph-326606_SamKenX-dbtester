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
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/cheggaaa/pb"
	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	v3 "github.com/coreos/etcd/clientv3"
	"github.com/samuel/go-zookeeper/zk"
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
	uniqueKeys   bool
)

func init() {
	Command.AddCommand(putCmd)
	putCmd.Flags().IntVar(&keySize, "key-size", 8, "Key size of put request")
	putCmd.Flags().IntVar(&valSize, "val-size", 8, "Value size of put request")
	putCmd.Flags().IntVar(&putTotal, "total", 10000, "Total number of put requests")
	putCmd.Flags().IntVar(&keySpaceSize, "key-space-size", 1, "Maximum possible keys")
	putCmd.Flags().BoolVar(&seqKeys, "sequential-keys", false, "Use sequential keys")
	putCmd.Flags().BoolVarP(&uniqueKeys, "unique-keys", "u", false, "Use unique keys (do not duplicate with sequential-keys)")
}

type request struct {
	etcdOp v3.Op
	zkOp   zkOp
}

type zkOp struct {
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

	switch database {
	case "etcd":
		clients := mustCreateClients(totalClients, totalConns)
		for i := range clients {
			wg.Add(1)
			go doPut(context.Background(), clients[i], requests)
		}
		defer func() {
			for i := range clients {
				clients[i].Close()
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
			go doPutZk(context.Background(), conns[i], requests)
		}
	default:
		log.Fatalf("unknown database %s", database)
	}

	pdoneC := printReport(results)

	go func() {
		for i := 0; i < putTotal; i++ {
			if seqKeys {
				binary.PutVarint(k, int64(i%keySpaceSize))
			} else if uniqueKeys {
				k = keys[i]
			} else {
				binary.PutVarint(k, int64(rand.Intn(keySpaceSize)))
			}
			switch database {
			case "etcd":
				requests <- request{etcdOp: v3.OpPut(string(k), v)}
			case "zk":
				requests <- request{zkOp: zkOp{key: string(k), value: []byte(v)}}
			}
		}
		close(requests)
	}()

	wg.Wait()

	bar.Finish()

	close(results)
	<-pdoneC
}

func doPut(ctx context.Context, client v3.KV, requests <-chan request) {
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

func doPutZk(ctx context.Context, conn *zk.Conn, requests <-chan request) {
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

func max(n1, n2 int64) int64 {
	if n1 > n2 {
		return n1
	}
	return n2
}
