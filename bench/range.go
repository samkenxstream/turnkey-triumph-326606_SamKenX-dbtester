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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/cheggaaa/pb"
	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	v3 "github.com/coreos/etcd/clientv3"
	"github.com/samuel/go-zookeeper/zk"
)

// rangeCmd represents the range command
var rangeCmd = &cobra.Command{
	Use:   "range key [end-range]",
	Short: "Benchmark range",

	Run: rangeFunc,
}

var (
	rangeTotal       int
	rangeConsistency string
	singleKey        bool
)

func init() {
	Command.AddCommand(rangeCmd)
	rangeCmd.Flags().IntVar(&rangeTotal, "total", 10000, "Total number of range requests")
	rangeCmd.Flags().StringVar(&rangeConsistency, "consistency", "l", "Linearizable(l) or Serializable(s)")
	rangeCmd.Flags().BoolVar(&singleKey, "single-key", false, "'true' to get only one single key (automatic put before test)")
}

func rangeFunc(cmd *cobra.Command, args []string) {
	if singleKey { // write 'foo'
		switch database {
		case "etcd":
			fmt.Println("PUT 'foo' to etcd")
			clients := mustCreateClients(1, 1)
			_, err := clients[0].Do(context.Background(), v3.OpPut("foo", "bar"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println("Done with PUT 'foo' to etcd")
		case "zk":
			fmt.Println("PUT 'foo' to zookeeper")
			conn := mustCreateConnsZk(1)
			_, err := conn[0].Create("foo", []byte("bar"), zkCreateFlags, zkCreateAcl)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println("Done with PUT 'foo' to zookeeper")
		}
	} else if len(args) == 0 || len(args) > 2 {
		fmt.Fprintln(os.Stderr, cmd.Usage())
		os.Exit(1)
	}

	var k, end string
	if singleKey {
		k = "foo"
	} else {
		k = args[0]
		if len(args) == 2 {
			end = args[1]
		}
	}

	if rangeConsistency == "l" {
		fmt.Println("bench with linearizable range")
	} else if rangeConsistency == "s" {
		fmt.Println("bench with serializable range")
	} else {
		fmt.Fprintln(os.Stderr, cmd.Usage())
		os.Exit(1)
	}

	results = make(chan result)
	requests := make(chan request, totalClients)
	bar = pb.New(rangeTotal)

	bar.Format("Bom !")
	bar.Start()

	switch database {
	case "etcd":
		clients := mustCreateClients(totalClients, totalConns)
		for i := range clients {
			wg.Add(1)
			go doRange(clients[i].KV, requests)
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
			go doRangeZk(conns[i], requests)
		}
	default:
		log.Fatalf("unknown database %s", database)
	}

	pdoneC := printReport(results)
	go func() {
		for i := 0; i < rangeTotal; i++ {
			switch database {
			case "etcd":
				opts := []v3.OpOption{v3.WithRange(end)}
				if rangeConsistency == "s" {
					opts = append(opts, v3.WithSerializable())
				}
				requests <- request{etcdOp: v3.OpGet(k, opts...)}
			case "zk":
				requests <- request{zkOp: zkOp{key: k}}
			}
		}
		close(requests)
	}()

	wg.Wait()

	bar.Finish()

	close(results)
	<-pdoneC
}

func doRange(client v3.KV, requests <-chan request) {
	defer wg.Done()

	for req := range requests {
		op := req.etcdOp
		st := time.Now()
		_, err := client.Do(context.Background(), op)

		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		results <- result{errStr: errStr, duration: time.Since(st), happened: time.Now()}
		bar.Increment()
	}
}

func doRangeZk(conn *zk.Conn, requests <-chan request) {
	defer wg.Done()

	for req := range requests {
		op := req.zkOp
		st := time.Now()
		_, _, err := conn.Get(op.key)

		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		results <- result{errStr: errStr, duration: time.Since(st), happened: time.Now()}
		bar.Increment()
	}
}
