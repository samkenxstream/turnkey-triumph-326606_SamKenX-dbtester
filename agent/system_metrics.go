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
	"time"

	"github.com/gyuho/psn"
)

// collectMetrics starts collecting metrics.
func collectMetrics(fs *flags, t *transporterServer) error {
	if fs == nil || t == nil || t.cmd == nil {
		return fmt.Errorf("cannot find process to track (%+v, %+v)", fs, t)
	}
	plog.Infof("starting collecting metrics [database %q | PID: %d | disk device: %q | network interface: %q]",
		t.req.Database, t.pid, fs.diskDevice, fs.networkInterface)

	if err := os.RemoveAll(fs.systemMetricsLog); err != nil {
		return err
	}

	c := psn.NewCSV(fs.systemMetricsLog, t.pid, fs.diskDevice, fs.networkInterface)
	if err := c.Add(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				if err := c.Add(); err != nil {
					plog.Errorf("psn.CSV.Add error (%v)", err)
					continue
				}

			case <-t.uploadSig:
				plog.Info("upload signal received; returning")
				return

			case sig := <-notifier:
				plog.Infof("signal received %q", sig.String())
				return
			}
		}
	}()
	return nil
}
