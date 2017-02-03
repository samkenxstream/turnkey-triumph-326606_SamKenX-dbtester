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
	"context"
	"time"

	"github.com/coreos/dbtester/agent/agentpb"
	"google.golang.org/grpc"
)

func bcastReq(cfg Config, op agentpb.Request_Operation) (map[int]agentpb.Response, error) {
	req := cfg.ToRequest()
	req.Operation = op

	type result struct {
		idx int
		r   agentpb.Response
	}
	donec, errc := make(chan result), make(chan error)
	for i := range cfg.PeerIPs {
		go func(i int) {
			if resp, err := sendReq(cfg.AgentEndpoints[i], req, i); err != nil {
				errc <- err
			} else {
				donec <- result{idx: i, r: *resp}
			}
		}(i)
		time.Sleep(time.Second)
	}

	im := make(map[int]agentpb.Response)

	var errs []error
	for cnt := 0; cnt != len(cfg.PeerIPs); cnt++ {
		select {
		case rs := <-donec:
			im[rs.idx] = rs.r
		case err := <-errc:
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return im, nil
}

func sendReq(ep string, req agentpb.Request, i int) (*agentpb.Response, error) {
	req.ServerIndex = uint32(i)
	req.ZookeeperMyID = uint32(i + 1)

	plog.Infof("sending message [index: %d | operation: %q | database: %q | endpoint: %q]", i, req.Operation, req.Database, ep)

	conn, err := grpc.Dial(ep, grpc.WithInsecure())
	if err != nil {
		plog.Errorf("grpc.Dial connecting error (%v) [index: %d | endpoint: %q]", err, i, ep)
		return nil, err
	}

	defer conn.Close()

	cli := agentpb.NewTransporterClient(conn)

	// give enough timeout
	// e.g. uploading logs takes longer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	resp, err := cli.Transfer(ctx, &req)
	cancel()
	if err != nil {
		plog.Errorf("cli.Transfer error (%v) [index: %d | endpoint: %q]", err, i, ep)
		return nil, err
	}

	plog.Infof("got response [index: %d | endpoint: %q | response: %+v]", i, ep, resp)
	return resp, nil
}

func step1StartDatabase(cfg Config) error {
	_, err := bcastReq(cfg, agentpb.Request_Start)
	return err
}
