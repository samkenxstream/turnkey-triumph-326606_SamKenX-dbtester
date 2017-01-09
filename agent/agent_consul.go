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

package agent

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/coreos/dbtester/agent/agentpb"
)

// startConsul starts Consul.
func startConsul(fs *flags, t *transporterServer, req *agentpb.Request) (*exec.Cmd, error) {
	if !exist(fs.consulExec) {
		return nil, fmt.Errorf("Consul binary %q does not exist", globalFlags.consulExec)
	}

	if err := os.RemoveAll(fs.consulDataDir); err != nil {
		return nil, err
	}

	peerIPs := strings.Split(req.PeerIPString, "___")

	var flags []string
	switch t.req.ServerIndex {
	case 0: // leader
		flags = []string{
			"agent",
			"-server",
			"-data-dir", fs.consulDataDir,
			"-bind", peerIPs[t.req.ServerIndex],
			"-client", peerIPs[t.req.ServerIndex],
			"-bootstrap-expect", "3",
		}

	default:
		flags = []string{
			"agent",
			"-server",
			"-data-dir", fs.consulDataDir,
			"-bind", peerIPs[t.req.ServerIndex],
			"-client", peerIPs[t.req.ServerIndex],
			"-join", peerIPs[0],
		}
	}
	flagString := strings.Join(flags, " ")

	cmd := exec.Command(fs.consulExec, flags...)
	cmd.Stdout = t.databaseLogFile
	cmd.Stderr = t.databaseLogFile
	cs := fmt.Sprintf("%s %s", cmd.Path, flagString)

	plog.Infof("starting database %q", cs)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	t.cmd = cmd
	t.pid = int64(cmd.Process.Pid)
	plog.Infof("started database %q (PID: %d)", cs, t.pid)

	return cmd, nil
}
