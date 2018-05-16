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

// Package control controls the database agents and benchmark testers.
package control

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/coreos/dbtester"
	"github.com/coreos/dbtester/dbtesterpb"
	"github.com/coreos/dbtester/pkg/ntp"

	"github.com/coreos/etcd/pkg/netutil"
	"github.com/gyuho/linux-inspect/df"
	"github.com/gyuho/linux-inspect/inspect"
	"github.com/gyuho/linux-inspect/top"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
	dn, err := df.GetDevice("/")
	if err != nil {
		lg.Warn("cannot get disk device mounted at '/'", zap.Error(err))
	}
	nm, err := netutil.GetDefaultInterfaces()
	if err != nil {
		lg.Warn("cannot detect default network interface", zap.Error(err))
	}
	var nt string
	for k := range nm {
		nt = k
		break
	}

	ids := dbtesterpb.GetAllDatabaseIDs()
	Command.PersistentFlags().StringVar(&databaseID, "database-id", ids[0], strings.Join(ids, ", "))
	Command.PersistentFlags().StringVarP(&configPath, "config", "c", "", "YAML configuration file path.")
	Command.PersistentFlags().StringVar(&diskDevice, "disk-device", dn, "Disk device to collect disk statistics metrics from.")
	Command.PersistentFlags().StringVar(&networkInterface, "network-interface", nt, "Network interface to record in/outgoing packets.")
}

func commandFunc(cmd *cobra.Command, args []string) error {
	if !dbtesterpb.IsValidDatabaseID(databaseID) {
		return fmt.Errorf("database id %q is unknown", databaseID)
	}

	cfg, err := dbtester.ReadConfig(configPath, false)
	if err != nil {
		return err
	}
	gcfg, ok := cfg.DatabaseIDToConfigClientMachineAgentControl[databaseID]
	if !ok {
		return fmt.Errorf("%q is not found", databaseID)
	}

	if gcfg.ConfigClientMachineBenchmarkSteps.Step2StressDatabase {
		switch gcfg.ConfigClientMachineBenchmarkOptions.Type {
		case "write":
		case "read":
		case "read-oneshot":
		default:
			return fmt.Errorf("%q is not supported", gcfg.ConfigClientMachineBenchmarkOptions.Type)
		}
	}

	pid := int64(os.Getpid())
	lg.Info(
		"starting collecting system metrics",
		zap.String("system-metrics-path", cfg.ConfigClientMachineInitial.ClientSystemMetricsPath),
		zap.String("disk-device", diskDevice),
		zap.String("network-device", networkInterface),
		zap.Int64("pid", pid),
	)
	if err = os.RemoveAll(cfg.ConfigClientMachineInitial.ClientSystemMetricsPath); err != nil {
		return err
	}
	tcfg := &top.Config{
		Exec:           top.DefaultExecPath,
		IntervalSecond: 1,
		PID:            pid,
	}
	var metricsCSV *inspect.CSV
	metricsCSV, err = inspect.NewCSV(
		cfg.ConfigClientMachineInitial.ClientSystemMetricsPath,
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
					lg.Warn("inspect.CSV.Add error", zap.Error(err))
					continue
				}

			case <-donec:
				lg.Info("finishing collecting system metrics; saving CSV", zap.String("path", cfg.ConfigClientMachineInitial.ClientSystemMetricsPath))

				if err := metricsCSV.Save(); err != nil {
					lg.Warn("inspect.CSV.Save failed", zap.String("path", metricsCSV.FilePath), zap.Error(err))
				} else {
					lg.Info("saved CSV", zap.String("path", metricsCSV.FilePath))
				}

				interpolated, err := metricsCSV.Interpolate()
				if err != nil {
					lg.Fatal("inspect.CSV.Interpolate failed", zap.String("path", metricsCSV.FilePath), zap.Error(err))
				}
				interpolated.FilePath = cfg.ConfigClientMachineInitial.ClientSystemMetricsInterpolatedPath
				if err := interpolated.Save(); err != nil {
					lg.Warn("inspect.CSV.Save failed", zap.String("path", interpolated.FilePath), zap.Error(err))
				} else {
					lg.Info("saved CSV", zap.String("path", interpolated.FilePath))
				}

				close(sysdonec)

				lg.Info("finished collecting system metrics")
				return
			}
		}
	}()

	no, nerr := ntp.DefaultSync()
	lg.Info("npt update output", zap.String("output", no))
	if nerr != nil {
		lg.Warn("ntp update failed", zap.Error(nerr))
	}

	println()
	if gcfg.ConfigClientMachineBenchmarkSteps.Step1StartDatabase {
		lg.Info("step 1: starting databases...")
		if _, err = cfg.BroadcaseRequest(databaseID, dbtesterpb.Operation_Start); err != nil {
			return err
		}
	}

	if gcfg.ConfigClientMachineBenchmarkSteps.Step2StressDatabase {
		println()
		time.Sleep(5 * time.Second)
		println()
		lg.Info("step 2: starting tests...")
		if err = cfg.Stress(databaseID); err != nil {
			return err
		}
	}

	if gcfg.ConfigClientMachineBenchmarkSteps.Step3StopDatabase {
		println()
		time.Sleep(5 * time.Second)
		println()
		lg.Info("step 3: stopping tests...")
		var idxToResp map[int]dbtesterpb.Response
		for i := 0; i < 5; i++ {
			idxToResp, err = cfg.BroadcaseRequest(databaseID, dbtesterpb.Operation_Stop)
			if err != nil {
				lg.Warn("STOP failed", zap.Int("i", i), zap.Error(err))
				time.Sleep(300 * time.Millisecond)
				continue
			}
			break
		}
		for idx := range gcfg.AgentEndpoints {
			lg.Info("stop response", zap.String("response", fmt.Sprintf("%+v", idxToResp[idx])))
		}

		println()
		time.Sleep(time.Second)
		println()
		lg.Info("step 3: saving responses...")
		if err = cfg.SaveDiskSpaceUsageSummary(databaseID, idxToResp); err != nil {
			return err
		}
	}

	close(donec)
	<-sysdonec

	if gcfg.ConfigClientMachineBenchmarkSteps.Step4UploadLogs {
		println()
		time.Sleep(3 * time.Second)
		println()
		lg.Info("step 4: uploading logs...")
		if err = cfg.UploadToGoogle(databaseID, cfg.ConfigClientMachineInitial.LogPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.ConfigClientMachineInitial.ClientSystemMetricsPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.ConfigClientMachineInitial.ClientSystemMetricsInterpolatedPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.ConfigClientMachineInitial.ClientLatencyThroughputTimeseriesPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.ConfigClientMachineInitial.ClientLatencyDistributionAllPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.ConfigClientMachineInitial.ClientLatencyDistributionPercentilePath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.ConfigClientMachineInitial.ClientLatencyDistributionSummaryPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.ConfigClientMachineInitial.ClientLatencyByKeyNumberPath); err != nil {
			return err
		}
		if err = cfg.UploadToGoogle(databaseID, cfg.ConfigClientMachineInitial.ServerDiskSpaceUsageSummaryPath); err != nil {
			return err
		}
	}

	lg.Info("all done!")
	return nil
}
