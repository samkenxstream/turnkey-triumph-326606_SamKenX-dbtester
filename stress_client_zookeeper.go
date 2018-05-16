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
	"errors"
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var (
	zkCreateFlags = int32(0)
	zkCreateACL   = zk.WorldACL(zk.PermAll)
)

type zkOp struct {
	key       string
	value     []byte
	staleRead bool
}

func mustCreateConnsZk(endpoints []string, total int64) []*zk.Conn {
	zks := make([]*zk.Conn, total)
	for i := range zks {
		endpoint := endpoints[dialTotal%len(endpoints)]
		dialTotal++
		conn, _, err := zk.Connect([]string{endpoint}, time.Second)
		if err != nil {
			panic(err)
		}
		zks[i] = conn
	}
	return zks
}

func newPutCreateZK(conn *zk.Conn) ReqHandler {
	return func(ctx context.Context, req *request) error {
		op := req.zkOp
		_, err := conn.Create(op.key, op.value, zkCreateFlags, zkCreateACL)
		return err
	}
}

func newPutOverwriteZK(conn *zk.Conn) ReqHandler {
	// samekey
	return func(ctx context.Context, req *request) error {
		op := req.zkOp
		_, err := conn.Set(op.key, op.value, int32(-1))
		return err
	}
}

func newGetZK(conn *zk.Conn) ReqHandler {
	return func(ctx context.Context, req *request) error {
		errt := ""
		if !req.zkOp.staleRead {
			_, err := conn.Sync("/" + req.zkOp.key)
			if err != nil {
				errt += err.Error()
			}
		}
		_, _, err := conn.Get("/" + req.zkOp.key)
		if err != nil {
			if errt != "" {
				errt += "; "
			}
			errt += fmt.Sprintf("%q while getting %q", err.Error(), "/"+req.zkOp.key)
		}
		if errt != "" {
			return errors.New(errt)
		}
		return nil
	}
}

func getTotalKeysZk(lg *zap.Logger, endpoints []string) map[string]int64 {
	rs := make(map[string]int64)
	stats, ok := zk.FLWSrvr(endpoints, 5*time.Second)
	if !ok {
		lg.Sugar().Infof("getTotalKeysZk failed with %+v", stats)
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
