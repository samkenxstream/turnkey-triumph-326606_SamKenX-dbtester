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
)

// startConsul starts Consul.
func startConsul(fs *flags, t *transporterServer) error {
	if !exist(fs.consulExec) {
		return fmt.Errorf("Consul binary %q does not exist", globalFlags.consulExec)
	}

	if err := os.RemoveAll(fs.consulDataDir); err != nil {
		return err
	}

	peerIPs := strings.Split(t.req.PeerIPsString, "___")

	var flags []string
	switch t.req.IPIndex {
	case 0: // leader
		flags = []string{
			"agent",
			"-server",
			"-data-dir", fs.consulDataDir,
			"-bind", peerIPs[t.req.IPIndex],
			"-client", peerIPs[t.req.IPIndex],
			"-bootstrap-expect", "3",
		}

	default:
		flags = []string{
			"agent",
			"-server",
			"-data-dir", fs.consulDataDir,
			"-bind", peerIPs[t.req.IPIndex],
			"-client", peerIPs[t.req.IPIndex],
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
		return err
	}
	t.cmd = cmd
	t.cmdWait = make(chan struct{})
	t.pid = int64(cmd.Process.Pid)
	plog.Infof("started database %q (PID: %d)", cs, t.pid)

	return nil
}
