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
	consulapi "github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
)

type consulOp struct {
	key       string
	value     []byte
	staleRead bool
}

func mustCreateConnsConsul(endpoints []string, total int64) []*consulapi.KV {
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

func newPutConsul(conn *consulapi.KV) ReqHandler {
	return func(ctx context.Context, req *request) error {
		op := req.consulOp
		_, err := conn.Put(&consulapi.KVPair{Key: op.key, Value: op.value}, nil)
		return err
	}
}

func newGetConsul(conn *consulapi.KV) ReqHandler {
	return func(ctx context.Context, req *request) error {
		opt := &consulapi.QueryOptions{}
		if req.consulOp.staleRead {
			opt.AllowStale = true
			opt.RequireConsistent = false
		}
		if !req.consulOp.staleRead {
			opt.AllowStale = false
			opt.RequireConsistent = true
		}
		_, _, err := conn.Get(req.consulOp.key, opt)
		return err
	}
}

func getTotalKeysConsul(endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	for _, ep := range endpoints {
		rs[ep] = 0 // not supported in consul
	}
	return rs
}
