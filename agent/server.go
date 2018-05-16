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

	"github.com/coreos/dbtester/dbtesterpb"
	"github.com/coreos/dbtester/pkg/fileinspect"

	"github.com/gyuho/linux-inspect/inspect"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// implements dbtesterpb.TransporterServer
type transporterServer struct {
	lg  *zap.Logger
	req dbtesterpb.Request

	databaseLogFile      *os.File
	proxyDatabaseLogfile *os.File
	clientNumPath        string

	// cmd is the main process that's running the database
	cmd *exec.Cmd
	// cmdWait channel is closed
	// after database process is closed
	cmdWait chan struct{}

	pid int64

	proxyCmd     *exec.Cmd
	proxyCmdWait chan struct{}
	proxyPid     int64

	metricsCSV *inspect.CSV

	// trigger log uploads to cloud storage
	// this should be triggered before we shut down
	// the agent server
	uploadSig chan struct{}
	csvReady  chan struct{}

	// notified after all tests finish
	notifier chan os.Signal
}

// NewServer returns a new server that implements gRPC interface.
func NewServer(lg *zap.Logger) dbtesterpb.TransporterServer {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	return &transporterServer{
		lg:            lg,
		clientNumPath: globalFlags.clientNumPath,
		uploadSig:     make(chan struct{}, 1),
		csvReady:      make(chan struct{}),
		notifier:      notifier,
	}
}

func (t *transporterServer) Transfer(ctx context.Context, req *dbtesterpb.Request) (*dbtesterpb.Response, error) {
	if req != nil {
		t.lg.Info(
			"received gRPC request",
			zap.String("operation", req.Operation.String()),
			zap.String("database-id", req.DatabaseID.String()),
			zap.Int64("client-number", req.CurrentClientNumber),
		)
	}

	if req.Operation == dbtesterpb.Operation_Start {
		f, err := openToAppend(globalFlags.databaseLog)
		if err != nil {
			return nil, err
		}
		t.databaseLogFile = f
		t.lg.Info("created database log file", zap.String("path", globalFlags.databaseLog))

		if req.DatabaseID == dbtesterpb.DatabaseID_zetcd__beta || req.DatabaseID == dbtesterpb.DatabaseID_cetcd__beta {
			proxyLog := globalFlags.databaseLog + "-" + t.req.DatabaseID.String()
			pf, err := openToAppend(proxyLog)
			if err != nil {
				return nil, err
			}
			t.proxyDatabaseLogfile = pf
			t.lg.Info("created database proxy log file", zap.String("path", proxyLog))
		}

		switch req.DatabaseID {
		case dbtesterpb.DatabaseID_etcd__other,
			dbtesterpb.DatabaseID_etcd__tip,
			dbtesterpb.DatabaseID_etcd__v3_2,
			dbtesterpb.DatabaseID_etcd__v3_3:
			t.lg.Info(
				"requested on etcd",
				zap.String("executable-binary-path", globalFlags.etcdExec),
				zap.String("data-directory", globalFlags.etcdDataDir),
			)

		case dbtesterpb.DatabaseID_zookeeper__r3_5_3_beta:
			t.lg.Info(
				"requested on Zookeeper",
				zap.String("working-directory", globalFlags.zkWorkDir),
				zap.String("data-directory", globalFlags.zkDataDir),
				zap.String("configuration-file", globalFlags.zkConfig),
			)

		case dbtesterpb.DatabaseID_consul__v1_0_2:
			t.lg.Info(
				"requested on Consul",
				zap.String("executable-binary-path", globalFlags.consulExec),
				zap.String("data-directory", globalFlags.consulDataDir),
			)

		case dbtesterpb.DatabaseID_zetcd__beta:
			t.lg.Info(
				"requested on zetcd",
				zap.String("executable-binary-path", globalFlags.zetcdExec),
				zap.String("data-directory", globalFlags.etcdDataDir),
			)

		case dbtesterpb.DatabaseID_cetcd__beta:
			t.lg.Info(
				"requested on cetcd",
				zap.String("executable-binary-path", globalFlags.cetcdExec),
				zap.String("data-directory", globalFlags.etcdDataDir),
			)
		}

		// re-use configurations for next requests
		t.req = *req
	}
	if req.Operation == dbtesterpb.Operation_Heartbeat {
		t.req.CurrentClientNumber = req.CurrentClientNumber
	}

	var diskSpaceUsageBytes int64
	switch req.Operation {
	case dbtesterpb.Operation_Start:
		switch t.req.DatabaseID {
		case dbtesterpb.DatabaseID_etcd__other,
			dbtesterpb.DatabaseID_etcd__tip,
			dbtesterpb.DatabaseID_etcd__v3_2,
			dbtesterpb.DatabaseID_etcd__v3_3,
			dbtesterpb.DatabaseID_zetcd__beta,
			dbtesterpb.DatabaseID_cetcd__beta:
			if err := startEtcd(&globalFlags, t); err != nil {
				return nil, err
			}
			switch t.req.DatabaseID {
			case dbtesterpb.DatabaseID_zetcd__beta:
				if err := startZetcd(&globalFlags, t); err != nil {
					return nil, err
				}
				go func() {
					defer close(t.proxyCmdWait)
					if err := t.proxyCmd.Wait(); err != nil {
						t.lg.Warn("zetcd t.proxyCmd.Wait() returned error", zap.Error(err))
						return
					}
					t.lg.Info("exiting zetcd", zap.String("executable-path", t.proxyCmd.Path))
				}()

			case dbtesterpb.DatabaseID_cetcd__beta:
				if err := startCetcd(&globalFlags, t); err != nil {
					return nil, err
				}
				go func() {
					defer close(t.proxyCmdWait)
					if err := t.proxyCmd.Wait(); err != nil {
						t.lg.Warn("cetcd t.proxyCmd.Wait() returned error", zap.Error(err))
						return
					}
					t.lg.Info("exiting cetcd", zap.String("executable-path", t.proxyCmd.Path))
				}()
			}

		case dbtesterpb.DatabaseID_zookeeper__r3_5_3_beta:
			if err := startZookeeper(&globalFlags, t); err != nil {
				return nil, err
			}

		case dbtesterpb.DatabaseID_consul__v1_0_2:
			if err := startConsul(&globalFlags, t); err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("unknown database %q", t.req.DatabaseID)
		}

		go func() {
			defer close(t.cmdWait)
			if err := t.cmd.Wait(); err != nil {
				t.lg.Warn("t.cmd.Wait() returned error", zap.Error(err))
				return
			}
			t.lg.Info("exiting", zap.String("executable-path", t.cmd.Path))
		}()
		if err := startMetrics(&globalFlags, t); err != nil {
			return nil, err
		}

	case dbtesterpb.Operation_Stop:
		if t.cmd == nil {
			return nil, fmt.Errorf("nil command")
		}

		// to collect more monitoring data
		t.lg.Info("waiting a few more seconds before stopping", zap.String("executable-path", t.cmd.Path))
		time.Sleep(3 * time.Second)

		// TODO: https://github.com/coreos/dbtester/issues/330
		t.lg.Info("sending", zap.String("syscall", syscall.SIGINT.String()), zap.Int64("pid", t.pid), zap.String("executable-path", t.cmd.Path))
		if err := t.cmd.Process.Signal(syscall.SIGINT); err != nil {
			t.lg.Warn("syscall.SIGINT failed", zap.Error(err))

			time.Sleep(3 * time.Second)
			t.lg.Info("sending", zap.String("syscall", syscall.SIGTERM.String()), zap.Int64("pid", t.pid), zap.String("executable-path", t.cmd.Path))
			if err := syscall.Kill(int(t.pid), syscall.SIGTERM); err != nil {
				t.lg.Warn("syscall.Kill failed", zap.Error(err))
			}
		}

		time.Sleep(time.Second)
		<-t.cmdWait

		if t.databaseLogFile != nil {
			t.databaseLogFile.Sync()
			t.databaseLogFile.Close()
		}
		t.lg.Info("stopped", zap.String("database", t.req.DatabaseID.String()), zap.Int64("pid", t.pid))

		if t.proxyCmd != nil {
			t.lg.Info("sending", zap.String("syscall", syscall.SIGINT.String()), zap.Int64("pid", t.proxyPid), zap.String("executable-path", t.proxyCmd.Path))
			if err := t.proxyCmd.Process.Signal(syscall.SIGINT); err != nil {
				t.lg.Warn("syscall.SIGINT failed", zap.Error(err))

				time.Sleep(3 * time.Second)
				t.lg.Info("sending", zap.String("syscall", syscall.SIGTERM.String()), zap.Int64("pid", t.proxyPid), zap.String("executable-path", t.proxyCmd.Path))
				if err := syscall.Kill(int(t.proxyPid), syscall.SIGTERM); err != nil {
					t.lg.Warn("syscall.Kill failed", zap.Error(err))
				}
			}

			<-t.proxyCmdWait
			t.lg.Info("stopped", zap.String("database", t.req.DatabaseID.String()), zap.Int64("pid", t.proxyPid))

			if t.proxyDatabaseLogfile != nil {
				t.proxyDatabaseLogfile.Sync()
				t.proxyDatabaseLogfile.Close()
			}
		}

		t.uploadSig <- struct{}{}
		<-t.csvReady

		if t.req.TriggerLogUpload {
			if err := uploadLog(&globalFlags, t); err != nil {
				return nil, err
			}
		}

		dbs, err := measureDatabasSize(globalFlags, req.DatabaseID)
		if err != nil {
			return nil, err
		}
		diskSpaceUsageBytes = dbs

	case dbtesterpb.Operation_Heartbeat:
		t.lg.Info("overwriting clients number", zap.Int64("number", t.req.CurrentClientNumber), zap.String("number-path", t.clientNumPath))
		if err := toFile(fmt.Sprintf("%d", t.req.CurrentClientNumber), t.clientNumPath); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Not implemented %v", req.Operation)
	}

	t.lg.Info("Transfer success!")
	return &dbtesterpb.Response{Success: true, DiskSpaceUsageBytes: diskSpaceUsageBytes}, nil
}

func measureDatabasSize(flg flags, rdb dbtesterpb.DatabaseID) (int64, error) {
	switch rdb {
	case dbtesterpb.DatabaseID_etcd__other,
		dbtesterpb.DatabaseID_etcd__tip,
		dbtesterpb.DatabaseID_etcd__v3_2,
		dbtesterpb.DatabaseID_etcd__v3_3,
		dbtesterpb.DatabaseID_cetcd__beta,
		dbtesterpb.DatabaseID_zetcd__beta:
		return fileinspect.Size(flg.etcdDataDir)

	case dbtesterpb.DatabaseID_zookeeper__r3_5_3_beta:
		return fileinspect.Size(flg.zkDataDir)

	case dbtesterpb.DatabaseID_consul__v1_0_2:
		return fileinspect.Size(flg.consulDataDir)

	default:
		return 0, fmt.Errorf("uknown %q", rdb)
	}
}
