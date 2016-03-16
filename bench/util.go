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
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	// dialTotal counts the number of mustCreateConn calls so that endpoint
	// connections can be handed out in round-robin order
	dialTotal int
)

func mustCreateConnsZk(total uint) []*zk.Conn {
	endpoint := endpoints[dialTotal%len(endpoints)]
	dialTotal++
	zks := make([]*zk.Conn, total)
	for i := range zks {
		conn, _, err := zk.Connect([]string{endpoint}, time.Second)
		if err != nil {
			log.Fatal(err)
		}
		zks[i] = conn
	}
	return zks
}

func mustCreateConn() *clientv3.Client {
	endpoint := endpoints[dialTotal%len(endpoints)]
	dialTotal++
	cfg := clientv3.Config{Endpoints: []string{endpoint}}
	client, err := clientv3.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial error: %v\n", err)
		os.Exit(1)
	}
	return client
}

func mustCreateClients(totalClients, totalConns uint) []*clientv3.Client {
	fmt.Println("mustCreateClients")
	conns := make([]*clientv3.Client, totalConns)
	for i := range conns {
		conns[i] = mustCreateConn()
	}

	clients := make([]*clientv3.Client, totalClients)
	for i := range clients {
		clients[i] = conns[i%int(totalConns)]
	}
	return clients
}

func mustRandBytes(n int) []byte {
	rb := make([]byte, n)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate value: %v\n", err)
		os.Exit(1)
	}
	return rb
}

func multiRandBytes(bytesN, sliceN int) [][]byte {
	m := make(map[string]struct{})
	var rs [][]byte
	for len(rs) != sliceN {
		b := mustRandBytes(bytesN)
		if _, ok := m[string(b)]; !ok {
			rs = append(rs, b)
			m[string(b)] = struct{}{}
		}
	}
	return rs
}
