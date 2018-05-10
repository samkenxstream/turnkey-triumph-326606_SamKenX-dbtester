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
	"context"
	"fmt"
	"time"

	"github.com/coreos/dbtester/dbtesterpb"

	"google.golang.org/grpc"
)

// BroadcaseRequest sends request to all endpoints.
func (cfg *Config) BroadcaseRequest(databaseID string, op dbtesterpb.Operation) (map[int]dbtesterpb.Response, error) {
	gcfg, ok := cfg.DatabaseIDToConfigClientMachineAgentControl[databaseID]
	if !ok {
		return nil, fmt.Errorf("database id %q does not exist", databaseID)
	}

	type result struct {
		idx int
		r   dbtesterpb.Response
	}
	donec, errc := make(chan result), make(chan error)
	for i := range gcfg.AgentEndpoints {
		req, err := cfg.ToRequest(databaseID, op, i)
		if err != nil {
			return nil, err
		}
		ep := gcfg.AgentEndpoints[i]

		go func(i int, ep string, req *dbtesterpb.Request) {
			plog.Infof("sending message [index: %d | operation: %q | database: %q | endpoint: %q]", i, op, req.DatabaseID, ep)

			conn, err := grpc.Dial(ep, grpc.WithInsecure())
			if err != nil {
				plog.Errorf("grpc.Dial connecting error (%v) [index: %d | endpoint: %q]", err, i, ep)
				errc <- fmt.Errorf("%v (%q)", err, ep)
				return
			}
			defer conn.Close()

			// give enough timeout
			// e.g. uploading logs takes longer
			cli := dbtesterpb.NewTransporterClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			resp, err := cli.Transfer(ctx, req)
			cancel()
			if err != nil {
				plog.Errorf("cli.Transfer error (%v) [index: %d | endpoint: %q]", err, i, ep)
				errc <- fmt.Errorf("%v (%q)", err, ep)
				return
			}

			plog.Infof("got response [index: %d | endpoint: %q | response: %+v]", i, ep, resp)
			donec <- result{idx: i, r: *resp}
		}(i, ep, req)

		time.Sleep(time.Second)
	}

	im := make(map[int]dbtesterpb.Response)
	var errs []error
	for cnt := 0; cnt != len(gcfg.AgentEndpoints); cnt++ {
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
