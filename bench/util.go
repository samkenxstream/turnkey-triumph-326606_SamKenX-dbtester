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
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mrand "math/rand"

	clientv2 "github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	// dialTotal counts the number of mustCreateConn calls so that endpoint
	// connections can be handed out in round-robin order
	dialTotal int
)

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

func mustCreateClientsEtcd2(total uint) []clientv2.KeysAPI {
	cks := make([]clientv2.KeysAPI, total)
	for i := range cks {
		endpoint := endpoints[dialTotal%len(endpoints)]
		dialTotal++

		if !strings.HasPrefix(endpoint, "http://") {
			endpoint = "http://" + endpoint
		}
		cfg := clientv2.Config{
			Endpoints:               []string{endpoint},
			Transport:               clientv2.DefaultTransport,
			HeaderTimeoutPerRequest: time.Second,
		}
		c, err := clientv2.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		kapi := clientv2.NewKeysAPI(c)

		cks[i] = kapi
	}
	return cks
}

func mustCreateConnsZk(total uint) []*zk.Conn {
	zks := make([]*zk.Conn, total)
	for i := range zks {
		endpoint := endpoints[dialTotal%len(endpoints)]
		dialTotal++
		conn, _, err := zk.Connect([]string{endpoint}, time.Second)
		if err != nil {
			log.Fatal(err)
		}
		zks[i] = conn
	}
	return zks
}

func mustCreateConnsConsul(total uint) []*consulapi.KV {
	css := make([]*consulapi.KV, total)
	for i := range css {
		endpoint := endpoints[dialTotal%len(endpoints)]
		dialTotal++

		dcfg := consulapi.DefaultConfig()
		dcfg.Address = endpoint // x.x.x.x:8500
		cli, err := consulapi.NewClient(dcfg)
		if err != nil {
			log.Fatal(err)
		}

		css[i] = cli.KV()
	}
	return css
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

func randBytes(bytesN int) []byte {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	src := mrand.NewSource(time.Now().UnixNano())
	b := make([]byte, bytesN)
	for i, cache, remain := bytesN-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}

func multiRandBytes(bytesN, sliceN int) [][]byte {
	m := make(map[string]struct{})
	var rs [][]byte
	for len(rs) != sliceN {
		b := randBytes(bytesN)
		if _, ok := m[string(b)]; !ok {
			rs = append(rs, b)
			m[string(b)] = struct{}{}
		}
	}
	return rs
}

func toFile(txt, fpath string) error {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
	}
	defer f.Close()
	if _, err := f.WriteString(txt); err != nil {
		return err
	}
	return nil
}

func toMillisecond(d time.Duration) float64 {
	return d.Seconds() * 1000
}
