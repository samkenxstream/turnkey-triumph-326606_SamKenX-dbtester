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
	"net"
	"net/http"
	"strings"
	"time"

	clientv2 "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func mustCreateClientsEtcdv2(endpoints []string, total int64) []clientv2.KeysAPI {
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

type etcdv2Op struct {
	key   string
	value string
}

func newPutEtcd2(conn clientv2.KeysAPI) ReqHandler {
	return func(ctx context.Context, req *request) error {
		op := req.etcdv2Op
		_, err := conn.Set(context.Background(), op.key, op.value, nil)
		return err
	}
}

func newGetEtcd2(conn clientv2.KeysAPI) ReqHandler {
	return func(ctx context.Context, req *request) error {
		_, err := conn.Get(ctx, req.etcdv2Op.key, nil)
		return err
	}
}

func getTotalKeysEtcdv2(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	for _, ep := range endpoints {
		rs[ep] = 0 // not supported in metrics
	}
	return rs
}
