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
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	clientv2 "github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	zkCreateFlags = int32(0)
	zkCreateAcl   = zk.WorldACL(zk.PermAll)
)

type request struct {
	etcdv2Op etcdv2Op
	etcdv3Op clientv3.Op
	zkOp     zkOp
	consulOp consulOp
}

type etcdv2Op struct {
	key   string
	value string
}

type zkOp struct {
	key       string
	value     []byte
	staleRead bool
}

type consulOp struct {
	key       string
	value     []byte
	staleRead bool
}

var (
	// dialTotal counts the number of mustCreateConn calls so that endpoint
	// connections can be handed out in round-robin order
	dialTotal int
)

func mustCreateConnEtcdv3(endpoints []string) *clientv3.Client {
	endpoint := endpoints[dialTotal%len(endpoints)]
	dialTotal++
	cfg := clientv3.Config{
		Endpoints: []string{endpoint},
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial error: %v\n", err)
		os.Exit(1)
	}
	return client
}

type etcdv3ClientCfg struct {
	totalConns   int
	totalClients int
}

func mustCreateClientsEtcdv3(endpoints []string, cfg etcdv3ClientCfg) []*clientv3.Client {
	conns := make([]*clientv3.Client, cfg.totalConns)
	for i := range conns {
		conns[i] = mustCreateConnEtcdv3(endpoints)
	}

	clients := make([]*clientv3.Client, cfg.totalClients)
	for i := range clients {
		clients[i] = conns[i%int(cfg.totalConns)]
	}
	return clients
}

func mustCreateClientsEtcdv2(endpoints []string, total int) []clientv2.KeysAPI {
	cks := make([]clientv2.KeysAPI, total)
	for i := range cks {
		endpoint := endpoints[dialTotal%len(endpoints)]
		dialTotal++

		if !strings.HasPrefix(endpoint, "http://") {
			endpoint = "http://" + endpoint
		}

		tr := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}
		cfg := clientv2.Config{
			Endpoints:               []string{endpoint},
			Transport:               tr,
			HeaderTimeoutPerRequest: time.Second,
		}
		c, err := clientv2.New(cfg)
		if err != nil {
			plog.Fatal(err)
		}
		kapi := clientv2.NewKeysAPI(c)

		cks[i] = kapi
	}
	return cks
}

func mustCreateConnsZk(endpoints []string, total int) []*zk.Conn {
	zks := make([]*zk.Conn, total)
	for i := range zks {
		endpoint := endpoints[dialTotal%len(endpoints)]
		dialTotal++
		conn, _, err := zk.Connect([]string{endpoint}, time.Second)
		if err != nil {
			plog.Fatal(err)
		}
		zks[i] = conn
	}
	return zks
}

func mustCreateConnsConsul(endpoints []string, total int) []*consulapi.KV {
	css := make([]*consulapi.KV, total)
	for i := range css {
		endpoint := endpoints[dialTotal%len(endpoints)]
		dialTotal++

		dcfg := consulapi.DefaultConfig()
		dcfg.Address = endpoint // x.x.x.x:8500
		cli, err := consulapi.NewClient(dcfg)
		if err != nil {
			plog.Fatal(err)
		}

		css[i] = cli.KV()
	}
	return css
}
func getTotalKeysEtcdv2(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	for _, ep := range endpoints {
		rs[ep] = 0 // not supported in metrics
	}
	return rs
}

func getTotalKeysEtcdv3(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	for _, ep := range endpoints {
		if !strings.HasPrefix(ep, "http://") {
			ep = "http://" + ep
		}

		plog.Println("GET", ep+"/metrics")
		resp, err := http.Get(ep + "/metrics")
		if err != nil {
			plog.Println(err)
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
			if ts[0] == "etcd_debugging_mvcc_keys_total" {
				rs[ep] = int64(fv)
				break
			}
		}
		gracefulClose(resp)
	}

	plog.Println("getTotalKeysEtcdv3", rs)
	return rs
}

func getTotalKeysZk(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	stats, ok := zk.FLWSrvr(endpoints, 5*time.Second)
	if !ok {
		plog.Printf("getTotalKeysZk failed with %+v", stats)
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

func getTotalKeysConsul(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	for _, ep := range endpoints {
		rs[ep] = 0 // not supported in consul
	}
	return rs
}
