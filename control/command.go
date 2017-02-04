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
	"strings"
	"time"

	"github.com/coreos/dbtester/pkg/ntp"
	"github.com/coreos/etcd/pkg/netutil"
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

	println()
	if !cfg.Step1.SkipStartDatabase {
		no, nerr := ntp.DefaultSync()
		plog.Infof("npt update output: %q", no)
		plog.Infof("npt update error: %v", nerr)

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
