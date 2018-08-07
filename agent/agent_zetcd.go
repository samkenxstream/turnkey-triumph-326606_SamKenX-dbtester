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
	"os/exec"
	"strings"

	"github.com/etcd-io/dbtester/dbtesterpb"

	"go.uber.org/zap"
)

// startZetcd starts zetcd. This assumes that etcd is already started.
func startZetcd(fs *flags, t *transporterServer) error {
	if !exist(fs.zetcdExec) {
		return fmt.Errorf("zetcd binary %q does not exist", globalFlags.zetcdExec)
	}

	peerIPs := strings.Split(t.req.PeerIPsString, "___")
	clientURLs := make([]string, len(peerIPs))
	for i, u := range peerIPs {
		clientURLs[i] = fmt.Sprintf("http://%s:2379", u)
	}

	var flags []string
	switch t.req.DatabaseID {
	case dbtesterpb.DatabaseID_zetcd__beta:
		flags = []string{
			// "-zkaddr", "0.0.0.0:2181",
			"-zkaddr", fmt.Sprintf("%s:2181", peerIPs[t.req.IPIndex]),
			"-endpoint", clientURLs[t.req.IPIndex],
		}

	default:
		return fmt.Errorf("database ID %q is not supported", t.req.DatabaseID)
	}

	flagString := strings.Join(flags, " ")

	cmd := exec.Command(fs.zetcdExec, flags...)
	cmd.Stdout = t.proxyDatabaseLogfile
	cmd.Stderr = t.proxyDatabaseLogfile
	cs := fmt.Sprintf("%s %s", cmd.Path, flagString)

	t.lg.Info("starting database", zap.String("command", cs))
	if err := cmd.Start(); err != nil {
		return err
	}
	t.proxyCmd = cmd
	t.proxyCmdWait = make(chan struct{})
	t.proxyPid = int64(cmd.Process.Pid)
	t.lg.Info("started database", zap.String("command", cs), zap.Int64("pid", t.proxyPid))

	return nil
}
