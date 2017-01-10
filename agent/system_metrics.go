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
	"path/filepath"
	"time"

	"github.com/gyuho/psn"
)

// startMetrics starts collecting metrics.
func startMetrics(fs *flags, t *transporterServer) error {
	if fs == nil || t == nil || t.cmd == nil {
		return fmt.Errorf("cannot find process to track (%+v, %+v)", fs, t)
	}
	plog.Infof("starting collecting metrics [database %q | PID: %d | disk device: %q | network interface: %q]",
		t.req.Database, t.pid, fs.diskDevice, fs.networkInterface)

	if err := os.RemoveAll(fs.systemMetricsLog); err != nil {
		return err
	}

	epath := filepath.Join(homeDir(), "etcd-client-num")
	if err := toFile(fmt.Sprintf("%d", t.req.ClientNum), epath); err != nil {
		return err
	}

	c := psn.NewCSV(fs.systemMetricsLog, t.pid, fs.diskDevice, fs.networkInterface, epath)
	if err := c.Add(); err != nil {
		return err
	}
	t.metricsCSV = c

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				if err := c.Add(); err != nil {
					plog.Errorf("psn.CSV.Add error (%v)", err)
					continue
				}

			case <-t.uploadSig:
				plog.Infof("upload signal received; saving CSV at %q", t.metricsCSV.FilePath)
				if err := t.metricsCSV.Save(); err != nil {
					plog.Errorf("psn.CSV.Save error %v", err)
				} else {
					plog.Infof("CSV saved at %q", t.metricsCSV.FilePath)
				}
				close(t.csvReady)
				return

			case sig := <-t.notifier:
				plog.Infof("signal received %q", sig.String())
				return
			}
		}
	}()
	return nil
}
