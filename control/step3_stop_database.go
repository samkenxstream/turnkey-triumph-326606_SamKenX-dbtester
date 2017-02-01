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
	"fmt"
	"time"

	"github.com/coreos/dbtester/agent/agentpb"
)

func step3StopDatabase(cfg Config) (map[int]agentpb.Response, error) {
	switch cfg.Step3.Action {
	case "stop":
		plog.Info("step 3: stopping databases...")
		var rm map[int]agentpb.Response
		var err error
		for i := 0; i < 5; i++ {
			if rm, err = bcastReq(cfg, agentpb.Request_Stop); err != nil {
				plog.Warningf("STOP failed at %s", cfg.PeerIPs[i])
				time.Sleep(300 * time.Millisecond)
				continue
			}
			break
		}
		return rm, err

	default:
		return nil, fmt.Errorf("unknown %q", cfg.Step3.Action)
	}
}
