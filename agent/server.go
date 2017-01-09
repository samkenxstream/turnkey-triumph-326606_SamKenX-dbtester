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
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/dbtester/agent/agentpb"
	"golang.org/x/net/context"
)

// implements agentpb.TransporterServer
type transporterServer struct {
	req agentpb.Request

	databaseLogFile      *os.File
	proxyDatabaseLogfile *os.File

	cmd *exec.Cmd
	pid int64

	proxyCmd *exec.Cmd
	proxyPid int64

	// trigger log uploads to cloud storage
	// this should be triggered before we shut down
	// the agent server
	uploadSig chan struct{}

	// notified after all tests finish
	notifier chan os.Signal
}

// NewServer returns a new server that implements gRPC interface.
func NewServer() agentpb.TransporterServer {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	return &transporterServer{
		uploadSig: make(chan struct{}, 1),
		notifier:  notifier,
	}
}

func (t *transporterServer) Transfer(ctx context.Context, r *agentpb.Request) (*agentpb.Response, error) {
	if r != nil {
		plog.Infof("received gRPC request %q with database %q", r.Operation, r.Database)
	}

	if r.Operation == agentpb.Request_Start {
		f, err := openToAppend(globalFlags.databaseLog)
		if err != nil {
			return nil, err
		}
		t.databaseLogFile = f

		plog.Infof("agent log path: %q", globalFlags.agentLog)
		plog.Infof("database log path: %q", globalFlags.databaseLog)
		if r.Database == agentpb.Request_zetcd || r.Database == agentpb.Request_cetcd {
			proxyLog := globalFlags.databaseLog + "-" + t.req.Database.String()
			pf, err := openToAppend(proxyLog)
			if err != nil {
				return nil, err
			}
			t.proxyDatabaseLogfile = pf
			plog.Infof("proxy-database log path: %q", proxyLog)
		}
		plog.Infof("system metrics log path: %q", globalFlags.systemMetricsLog)

		switch r.Database {
		case agentpb.Request_ZooKeeper:
			plog.Infof("Zookeeper working directory: %q", globalFlags.zkWorkDir)
			plog.Infof("Zookeeper data directory: %q", globalFlags.zkDataDir)
			plog.Infof("Zookeeper configuration path: %q", globalFlags.zkConfig)

		case agentpb.Request_etcdv2, agentpb.Request_etcdv3:
			plog.Infof("etcd executable binary path: %q", globalFlags.etcdExec)
			plog.Infof("etcd data directory: %q", globalFlags.etcdDataDir)

		case agentpb.Request_zetcd:
			plog.Infof("zetcd executable binary path: %q", globalFlags.zetcdExec)
			plog.Infof("zetcd data directory: %q", globalFlags.etcdDataDir)

		case agentpb.Request_cetcd:
			plog.Infof("cetcd executable binary path: %q", globalFlags.cetcdExec)
			plog.Infof("cetcd data directory: %q", globalFlags.etcdDataDir)

		case agentpb.Request_Consul:
			plog.Infof("Consul executable binary path: %q", globalFlags.consulExec)
			plog.Infof("Consul data directory: %q", globalFlags.consulDataDir)
		}

		// re-use configurations for next requests
		t.req = *r
	}

	switch r.Operation {
	case agentpb.Request_Start:
		switch t.req.Database {
		case agentpb.Request_etcdv2, agentpb.Request_etcdv3, agentpb.Request_zetcd, agentpb.Request_cetcd:
			if err := startEtcd(&globalFlags, t); err != nil {
				plog.Errorf("startEtcd error %v", err)
				return nil, err
			}
			switch t.req.Database {
			case agentpb.Request_zetcd:
				if err := startZetcd(&globalFlags, t); err != nil {
					plog.Errorf("startZetcd error %v", err)
					return nil, err
				}
				go func() {
					if err := t.proxyCmd.Wait(); err != nil {
						plog.Errorf("cmd.Wait %q returned error %v", t.proxyCmd.Path, err)
						return
					}
					plog.Infof("exiting %q", t.proxyCmd.Path)
				}()
			case agentpb.Request_cetcd:
				if err := startCetcd(&globalFlags, t); err != nil {
					plog.Errorf("startCetcd error %v", err)
					return nil, err
				}
				go func() {
					if err := t.proxyCmd.Wait(); err != nil {
						plog.Errorf("cmd.Wait %q returned error %v", t.proxyCmd.Path, err)
						return
					}
					plog.Infof("exiting %q", t.proxyCmd.Path)
				}()
			}
		case agentpb.Request_ZooKeeper:
			if err := startZookeeper(&globalFlags, t); err != nil {
				plog.Errorf("startZookeeper error %v", err)
				return nil, err
			}
		case agentpb.Request_Consul:
			if err := startConsul(&globalFlags, t); err != nil {
				plog.Errorf("startConsul error %v", err)
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown database %q", t.req.Database)
		}

		go func() {
			if err := t.cmd.Wait(); err != nil {
				plog.Errorf("cmd.Wait %q returned error %v", t.cmd.Path, err)
				return
			}
			plog.Infof("exiting %q", t.cmd.Path)
		}()

		if err := collectMetrics(&globalFlags, t); err != nil {
			plog.Errorf("collectMetrics error %v", err)
			return nil, err
		}

	case agentpb.Request_Stop:
		if t.cmd == nil {
			return nil, fmt.Errorf("nil command")
		}

		// to collect more monitoring data
		plog.Infof("waiting a few more seconds before stopping %q", t.cmd.Path)
		time.Sleep(3 * time.Second)

		plog.Infof("sending %q to %q [PID: %d]", syscall.SIGINT, t.cmd.Path, t.pid)
		if err := t.cmd.Process.Signal(syscall.SIGINT); err != nil {
			return nil, err
		}
		time.Sleep(3 * time.Second)
		plog.Infof("sending %q to %q [PID: %d]", syscall.SIGTERM, t.cmd.Path, t.pid)
		if err := syscall.Kill(int(t.pid), syscall.SIGTERM); err != nil {
			plog.Warningf("syscall.Kill failed with %v", err)
		}
		if t.databaseLogFile != nil {
			t.databaseLogFile.Sync()
			t.databaseLogFile.Close()
		}
		plog.Infof("stopped binary %q [PID: %d]", t.req.Database.String(), t.pid)

		if t.proxyCmd != nil {
			plog.Infof("sending %q to %q [PID: %d]", syscall.SIGINT, t.proxyCmd.Path, t.proxyPid)
			if err := t.proxyCmd.Process.Signal(syscall.SIGINT); err != nil {
				return nil, err
			}
			time.Sleep(3 * time.Second)
			plog.Infof("sending %q to %q [PID: %d]", syscall.SIGTERM, t.proxyCmd.Path, t.proxyPid)
			if err := syscall.Kill(int(t.proxyPid), syscall.SIGTERM); err != nil {
				plog.Warningf("syscall.Kill failed with %v", err)
			}
			plog.Infof("stopped binary proxy for %q [PID: %d]", t.req.Database.String(), t.pid)
		}
		if t.proxyDatabaseLogfile != nil {
			t.proxyDatabaseLogfile.Sync()
			t.proxyDatabaseLogfile.Close()
		}

		t.uploadSig <- struct{}{}

		if err := uploadLog(&globalFlags, t); err != nil {
			plog.Warningf("uploadLog error %v", err)
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Not implemented %v", r.Operation)
	}

	plog.Info("Transfer success!")
	return &agentpb.Response{Success: true}, nil
}
