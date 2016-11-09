// Copyright 2016 CoreOS, Inc.
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
	"errors"

	"fmt"

	clientv2 "github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/samuel/go-zookeeper/zk"
	"golang.org/x/net/context"
)

type ReqHandler func(ctx context.Context, req *request) error

func newPutEtcd2(conn clientv2.KeysAPI) ReqHandler {
	return func(ctx context.Context, req *request) error {
		op := req.etcdv2Op
		_, err := conn.Set(context.Background(), op.key, op.value, nil)
		return err
	}
}

func newPutEtcd3(conn clientv3.KV) ReqHandler {
	return func(ctx context.Context, req *request) error {
		_, err := conn.Do(ctx, req.etcdv3Op)
		return err
	}
}

func newPutOverwriteZK(conn *zk.Conn) ReqHandler {
	return func(ctx context.Context, req *request) error {
		op := req.zkOp
		_, err := conn.Set(op.key, op.value, int32(-1))
		return err
	}
}

func newPutCreateZK(conn *zk.Conn) ReqHandler {
	// samekey
	return func(ctx context.Context, req *request) error {
		op := req.zkOp
		_, err := conn.Create(op.key, op.value, zkCreateFlags, zkCreateAcl)
		return err
	}
}

func newPutConsul(conn *consulapi.KV) ReqHandler {
	return func(ctx context.Context, req *request) error {
		op := req.consulOp
		_, err := conn.Put(&consulapi.KVPair{Key: op.key, Value: op.value}, nil)
		return err
	}
}

func newGetEtcd2(conn clientv2.KeysAPI) ReqHandler {
	return func(ctx context.Context, req *request) error {
		_, err := conn.Get(ctx, req.etcdv2Op.key, nil)
		return err
	}
}

func newGetEtcd3(conn clientv3.KV) ReqHandler {
	return func(ctx context.Context, req *request) error {
		_, err := conn.Do(ctx, req.etcdv3Op)
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

func newGetConsul(conn *consulapi.KV) ReqHandler {
	return func(ctx context.Context, req *request) error {
		_, _, err := conn.Get(req.consulOp.key, &consulapi.QueryOptions{AllowStale: req.consulOp.staleRead})
		return err
	}
}
