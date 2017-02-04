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
	"io/ioutil"
	"os"
	"strings"
	"time"

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
	Command.PersistentFlags().StringVarP(&configPath, "config", "c", "", "YAML configuration file path.")
	Command.PersistentFlags().StringVar(&diskDevice, "disk-device", dn, "Disk device to collect disk statistics metrics from.")
	Command.PersistentFlags().StringVar(&networkInterface, "network-interface", nt, "Network interface to record in/outgoing packets.")
}

func commandFunc(cmd *cobra.Command, args []string) error {
	cfg, err := ReadConfig(configPath)
	if err != nil {
		return err
	}
	switch cfg.Database {
	case "etcdv2":
	case "etcdv3":
	case "zookeeper":
	case "zetcd":
	case "consul":
	case "cetcd":
	default:
		return fmt.Errorf("%q is not supported", cfg.Database)
	}

	if !cfg.Step2.SkipStressDatabase {
		switch cfg.Step2.BenchType {
		case "write":
		case "read":
		case "read-oneshot":
		default:
			return fmt.Errorf("%q is not supported", cfg.Step2.BenchType)
		}
	}

	if cfg.Step4.UploadLogs {
		bts, err := ioutil.ReadFile(cfg.Step4.GoogleCloudStorageKeyPath)
		if err != nil {
			return err
		}
		cfg.Step4.GoogleCloudStorageKey = string(bts)
	}

	pid := int64(os.Getpid())
	plog.Infof("starting collecting system metrics at %q [disk device: %q | network interface: %q | PID: %d]", cfg.ClientSystemMetrics, diskDevice, networkInterface, pid)
	if err = os.RemoveAll(cfg.ClientSystemMetrics); err != nil {
		return err
	}
	tcfg := &psn.TopConfig{
		Exec:           psn.DefaultTopPath,
		IntervalSecond: 1,
		PID:            pid,
	}
	var metricsCSV *psn.CSV
	metricsCSV, err = psn.NewCSV(
		cfg.ClientSystemMetrics,
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
				plog.Infof("finishing collecting system metrics; saving CSV at %q", cfg.ClientSystemMetrics)

				if err := metricsCSV.Save(); err != nil {
					plog.Errorf("psn.CSV.Save(%q) error %v", metricsCSV.FilePath, err)
				} else {
					plog.Infof("CSV saved at %q", metricsCSV.FilePath)
				}

				interpolated, err := metricsCSV.Interpolate()
				if err != nil {
					plog.Fatalf("psn.CSV.Interpolate(%q) failed with %v", metricsCSV.FilePath, err)
				}
				interpolated.FilePath = cfg.ClientSystemMetricsInterpolated
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

	// protoc sorts the 'repeated' type data
	// encode in string to enforce ordering of IPs
	cfg.PeerIPString = strings.Join(cfg.PeerIPs, "___")
	cfg.AgentEndpoints = make([]string, len(cfg.PeerIPs))
	cfg.DatabaseEndpoints = make([]string, len(cfg.PeerIPs))
	for i := range cfg.PeerIPs {
		cfg.AgentEndpoints[i] = fmt.Sprintf("%s:%d", cfg.PeerIPs[i], cfg.AgentPort)
	}
	for i := range cfg.PeerIPs {
		cfg.DatabaseEndpoints[i] = fmt.Sprintf("%s:%d", cfg.PeerIPs[i], cfg.DatabasePort)
	}

	no, nerr := ntp.DefaultSync()
	plog.Infof("npt update output: %q", no)
	plog.Infof("npt update error: %v", nerr)

	println()
	if !cfg.Step1.SkipStartDatabase {
		plog.Info("step 1: starting databases...")
		if err = step1StartDatabase(cfg); err != nil {
			return err
		}
	}

	if !cfg.Step2.SkipStressDatabase {
		println()
		time.Sleep(5 * time.Second)
		plog.Info("step 2: starting tests...")
		if err = step2StressDatabase(cfg); err != nil {
			return err
		}
	}

	println()
	time.Sleep(5 * time.Second)
	idxToResponse, err := step3StopDatabase(cfg)
	if err != nil {
		plog.Warning(err)
	}
	for idx := range cfg.AgentEndpoints {
		plog.Infof("stop response: %+v", idxToResponse[idx])
	}

	println()
	time.Sleep(time.Second)
	saveDatasizeSummary(cfg, idxToResponse)

	close(donec)
	<-sysdonec

	if cfg.Step4.UploadLogs {
		println()
		time.Sleep(3 * time.Second)
		if err := step4UploadLogs(cfg); err != nil {
			return err
		}
	}

	plog.Info("all done!")
	return nil
}
