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
	"os"
	"time"

	"github.com/coreos/dbtester"
	"github.com/coreos/dbtester/dbtesterpb"
	"github.com/coreos/dbtester/pkg/netutil"
	"github.com/coreos/dbtester/pkg/ntp"
	"github.com/gyuho/psn"
	"github.com/spf13/cobra"
)

// Command implements 'control' command.
var Command = &cobra.Command{
	Use:   "control",
	Short: "Controls tests.",
	RunE:  commandFunc,
}

var databaseID string
var configPath string
var diskDevice string
var networkInterface string

func init() {
	dn, err := psn.GetDevice("/")
	if err != nil {
		plog.Warningf("cannot get disk device mounted at '/' (%v)", err)
	}
	nm, err := netutil.GetDefaultInterfaces()
	if err != nil {
		plog.Warningf("cannot detect default network interface (%v)", err)
	}
	var nt string
	for k := range nm {
		nt = k
		break
	}

	Command.PersistentFlags().StringVar(&databaseID, "database-id", "etcdv3", "etcdv2, etcdv3, zookeeper, consul, zetcd, cetcd.")
	Command.PersistentFlags().StringVarP(&configPath, "config", "c", "", "YAML configuration file path.")
	Command.PersistentFlags().StringVar(&diskDevice, "disk-device", dn, "Disk device to collect disk statistics metrics from.")
	Command.PersistentFlags().StringVar(&networkInterface, "network-interface", nt, "Network interface to record in/outgoing packets.")
}

func commandFunc(cmd *cobra.Command, args []string) error {
	switch databaseID {
	case "etcdv2":
	case "etcdv3":
	case "zookeeper":
	case "zetcd":
	case "consul":
	case "cetcd":
	default:
		return fmt.Errorf("%q is not supported", databaseID)
	}

	cfg, err := dbtester.ReadConfig(configPath, false)
	if err != nil {
		return err
	}

	gcfg, ok := cfg.DatabaseIDToTestGroup[databaseID]
	if !ok {
		return fmt.Errorf("%q is not found", databaseID)
	}

	if gcfg.BenchmarkSteps.Step2StressDatabase {
		switch gcfg.BenchmarkOptions.Type {
		case "write":
		case "read":
		case "read-oneshot":
		default:
			return fmt.Errorf("%q is not supported", gcfg.BenchmarkOptions.Type)
		}
	}

	pid := int64(os.Getpid())
	plog.Infof("starting collecting system metrics at %q [disk device: %q | network interface: %q | PID: %d]", cfg.Control.ClientSystemMetricsPath, diskDevice, networkInterface, pid)
	if err = os.RemoveAll(cfg.Control.ClientSystemMetricsPath); err != nil {
		return err
	}
	tcfg := &psn.TopConfig{
		Exec:           psn.DefaultTopPath,
		IntervalSecond: 1,
		PID:            pid,
	}
	var metricsCSV *psn.CSV
	metricsCSV, err = psn.NewCSV(
		cfg.Control.ClientSystemMetricsPath,
		pid,
		diskDevice,
		networkInterface,
		"",
		tcfg,
	)
	if err = metricsCSV.Add(); err != nil {
		return err
	}

	donec, sysdonec := make(chan struct{}), make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				if err := metricsCSV.Add(); err != nil {
					plog.Errorf("psn.CSV.Add error (%v)", err)
					continue
				}

			case <-donec:
				plog.Infof("finishing collecting system metrics; saving CSV at %q", cfg.Control.ClientSystemMetricsPath)

				if err := metricsCSV.Save(); err != nil {
					plog.Errorf("psn.CSV.Save(%q) error %v", metricsCSV.FilePath, err)
				} else {
					plog.Infof("CSV saved at %q", metricsCSV.FilePath)
				}

				interpolated, err := metricsCSV.Interpolate()
				if err != nil {
					plog.Fatalf("psn.CSV.Interpolate(%q) failed with %v", metricsCSV.FilePath, err)
				}
				interpolated.FilePath = cfg.Control.ClientSystemMetricsInterpolatedPath
				if err := interpolated.Save(); err != nil {
					plog.Errorf("psn.CSV.Save(%q) error %v", interpolated.FilePath, err)
				} else {
					plog.Infof("CSV saved at %q", interpolated.FilePath)
				}

				close(sysdonec)
				plog.Infof("finished collecting system metrics")
				return
			}
		}
	}()

	no, nerr := ntp.DefaultSync()
	plog.Infof("npt update output: %q", no)
	plog.Infof("npt update error: %v", nerr)

	println()
	if gcfg.BenchmarkSteps.Step1StartDatabase {
		plog.Info("step 1: starting databases...")
		if _, err = cfg.BroadcaseRequest(databaseID, dbtesterpb.Request_Start); err != nil {
			return err
		}
	}

	if gcfg.BenchmarkSteps.Step2StressDatabase {
		println()
		time.Sleep(5 * time.Second)
		println()
		plog.Info("step 2: starting tests...")
		if err = cfg.Stress(databaseID); err != nil {
			return err
		}
	}

	if gcfg.BenchmarkSteps.Step3StopDatabase {
		println()
		time.Sleep(5 * time.Second)
		println()
		plog.Info("step 3: stopping tests...")
		var idxToResp map[int]dbtesterpb.Response
		for i := 0; i < 5; i++ {
			idxToResp, err = cfg.BroadcaseRequest(databaseID, dbtesterpb.Request_Stop)
			if err != nil {
				plog.Warningf("#%d: STOP failed at %v", i, err)
				time.Sleep(300 * time.Millisecond)
				continue
			}
			break
		}
		for idx := range gcfg.AgentEndpoints {
			plog.Infof("stop response: %+v", idxToResp[idx])
		}

		println()
		time.Sleep(time.Second)
		println()
		plog.Info("step 3: saving responses...")
		if err = cfg.SaveDiskSpaceUsageSummary(databaseID, idxToResp); err != nil {
			return err
		}
	}

	close(donec)
	<-sysdonec

	if gcfg.BenchmarkSteps.Step4UploadLogs {
		println()
		time.Sleep(3 * time.Second)
		println()
		plog.Info("step 4: uploading logs...")
		if err = cfg.UploadToGoogle(databaseID, cfg.Control.LogPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.Control.ClientSystemMetricsPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.Control.ClientSystemMetricsInterpolatedPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.Control.ClientLatencyThroughputTimeseriesPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.Control.ClientLatencyDistributionAllPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.Control.ClientLatencyDistributionPercentilePath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.Control.ClientLatencyDistributionSummaryPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.Control.ClientLatencyByKeyNumberPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.Control.ServerDiskSpaceUsageSummaryPath); err != nil {
			return err
		}
	}

	plog.Info("all done!")
	return nil
}
