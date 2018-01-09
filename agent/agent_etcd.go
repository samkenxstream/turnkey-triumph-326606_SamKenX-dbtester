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

	"github.com/coreos/dbtester/dbtesterpb"
)

// startEtcd starts etcd v3.
func startEtcd(fs *flags, t *transporterServer) error {
	if !exist(fs.etcdExec) {
		return fmt.Errorf("etcd binary %q does not exist", globalFlags.etcdExec)
	}

	if err := os.RemoveAll(fs.etcdDataDir); err != nil {
		return err
	}

	peerIPs := strings.Split(t.req.PeerIPsString, "___")

	names := make([]string, len(peerIPs))
	clientURLs := make([]string, len(peerIPs))
	peerURLs := make([]string, len(peerIPs))
	members := make([]string, len(peerIPs))
	for i, u := range peerIPs {
		names[i] = fmt.Sprintf("etcd-%d", i+1)
		clientURLs[i] = fmt.Sprintf("http://%s:2379", u)
		peerURLs[i] = fmt.Sprintf("http://%s:2380", u)
		members[i] = fmt.Sprintf("%s=%s", names[i], peerURLs[i])
	}

	var flags []string
	switch t.req.DatabaseID {
	case dbtesterpb.DatabaseID_etcd__tip,
		dbtesterpb.DatabaseID_etcd__v3_2,
		dbtesterpb.DatabaseID_etcd__v3_3:
		flags = []string{
			"--name", names[t.req.IPIndex],
			"--data-dir", fs.etcdDataDir,
			"--quota-backend-bytes", fmt.Sprintf("%d", t.req.Flag_Etcd_Tip.QuotaSizeBytes),

			"--snapshot-count", fmt.Sprintf("%d", t.req.Flag_Etcd_Tip.SnapshotCount),

			"--listen-client-urls", clientURLs[t.req.IPIndex],
			"--advertise-client-urls", clientURLs[t.req.IPIndex],

			"--listen-peer-urls", peerURLs[t.req.IPIndex],
			"--initial-advertise-peer-urls", peerURLs[t.req.IPIndex],

			"--initial-cluster-token", "mytoken",
			"--initial-cluster", strings.Join(members, ","),
			"--initial-cluster-state", "new",
		}

	default:
		return fmt.Errorf("database ID %q is not supported", t.req.DatabaseID)
	}

	flagString := strings.Join(flags, " ")

	cmd := exec.Command(fs.etcdExec, flags...)
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
