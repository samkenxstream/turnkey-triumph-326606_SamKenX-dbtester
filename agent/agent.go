// Copyright 2016 CoreOS, Inc.
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
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/coreos/dbtester/remotestorage"
	"github.com/gyuho/psn/ps"
	"github.com/spf13/cobra"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type (
	Flags struct {
		GRPCPort         string
		WorkingDirectory string
	}

	// ZookeeperConfig is zookeeper configuration.
	// http://zookeeper.apache.org/doc/trunk/zookeeperAdmin.html
	ZookeeperConfig struct {
		TickTime       int
		DataDir        string
		ClientPort     string
		InitLimit      int
		SyncLimit      int
		MaxClientCnxns int64
		SnapCount      int64
		Peers          []ZookeeperPeer
	}
	ZookeeperPeer struct {
		MyID int
		IP   string
	}
)

var (
	shell = os.Getenv("SHELL")

	agentLogPath    = "agent.log"
	databaseLogPath = "database.log"
	monitorLogPath  = "monitor.csv"

	etcdBinaryPath   = filepath.Join(os.Getenv("GOPATH"), "bin/etcd")
	consulBinaryPath = filepath.Join(os.Getenv("GOPATH"), "bin/consul")
	javaBinaryPath   = "/usr/bin/java"

	etcdToken     = "etcd_token"
	etcdDataDir   = "data.etcd"
	consulDataDir = "data.consul"

	zkWorkingDir = "zookeeper"
	zkDataDir    = "zookeeper/data.zk"
	zkConfigPath = "zookeeper.config"
	zkTemplate   = `tickTime={{.TickTime}}
dataDir={{.DataDir}}
clientPort={{.ClientPort}}
initLimit={{.InitLimit}}
syncLimit={{.SyncLimit}}
maxClientCnxns={{.MaxClientCnxns}}
snapCount={{.SnapCount}}
{{range .Peers}}server.{{.MyID}}={{.IP}}:2888:3888
{{end}}
`
	zkConfigDefault = ZookeeperConfig{
		TickTime:       2000,
		ClientPort:     "2181",
		InitLimit:      5,
		SyncLimit:      5,
		MaxClientCnxns: 60,
		Peers: []ZookeeperPeer{
			{MyID: 1, IP: ""},
			{MyID: 2, IP: ""},
			{MyID: 3, IP: ""},
		},
	}

	Command = &cobra.Command{
		Use:   "agent",
		Short: "Database agent in remote servers.",
		RunE:  CommandFunc,
	}

	globalFlags = Flags{}
)

func init() {
	if len(shell) == 0 {
		shell = "sh"
	}
	Command.PersistentFlags().StringVar(&globalFlags.GRPCPort, "agent-port", ":3500", "Port to server agent gRPC server.")
	Command.PersistentFlags().StringVar(&globalFlags.WorkingDirectory, "working-directory", homeDir(), "Working directory.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	if !exist(globalFlags.WorkingDirectory) {
		return fmt.Errorf("%s does not exist", globalFlags.WorkingDirectory)
	}
	if !filepath.HasPrefix(agentLogPath, globalFlags.WorkingDirectory) {
		agentLogPath = filepath.Join(globalFlags.WorkingDirectory, agentLogPath)
	}

	f, err := openToAppend(agentLogPath)
	if err != nil {
		return err
	}
	defer f.Close()

	logger = zap.NewJSON(
		zap.Fields(zap.String("package", "dbtester/agent")),
		zap.Output(zap.AddSync(f)),
	)
	logger.Info("started serving gRPC",
		zap.String("port", globalFlags.GRPCPort),
	)

	var (
		grpcServer = grpc.NewServer()
		sender     = NewTransporterServer()
	)
	ln, err := net.Listen("tcp", globalFlags.GRPCPort)
	if err != nil {
		return err
	}

	RegisterTransporterServer(grpcServer, sender)

	return grpcServer.Serve(ln)
}

type transporterServer struct { // satisfy TransporterServer
	req     Request
	cmd     *exec.Cmd
	logfile *os.File
	pid     int
}

var databaseStopped = make(chan struct{})

func (t *transporterServer) Transfer(ctx context.Context, r *Request) (*Response, error) {
	peerIPs := strings.Split(r.PeerIPString, "___")
	if r.Operation == Request_Start || r.Operation == Request_Restart {
		if !filepath.HasPrefix(etcdDataDir, globalFlags.WorkingDirectory) {
			etcdDataDir = filepath.Join(globalFlags.WorkingDirectory, etcdDataDir)
		}
		if !filepath.HasPrefix(consulDataDir, globalFlags.WorkingDirectory) {
			consulDataDir = filepath.Join(globalFlags.WorkingDirectory, consulDataDir)
		}
		if !filepath.HasPrefix(zkWorkingDir, globalFlags.WorkingDirectory) {
			zkWorkingDir = filepath.Join(globalFlags.WorkingDirectory, zkWorkingDir)
		}
		if !filepath.HasPrefix(zkDataDir, globalFlags.WorkingDirectory) {
			zkDataDir = filepath.Join(globalFlags.WorkingDirectory, zkDataDir)
		}
		if !filepath.HasPrefix(databaseLogPath, globalFlags.WorkingDirectory) {
			databaseLogPath = filepath.Join(globalFlags.WorkingDirectory, databaseLogPath)
		}
		if !filepath.HasPrefix(monitorLogPath, globalFlags.WorkingDirectory) {
			monitorLogPath = filepath.Join(globalFlags.WorkingDirectory, monitorLogPath)
		}

		logger.Info("received gRPC request",
			zap.String("working_directory", globalFlags.WorkingDirectory),
			zap.String("working_directory_zookeeper", zkWorkingDir),
			zap.String("data_directory_etcd", etcdDataDir),
			zap.String("data_directory_consul", consulDataDir),
			zap.String("data_directory_zookeeper", zkDataDir),
			zap.String("database_log_path", databaseLogPath),
			zap.String("monitor_log_path", monitorLogPath),
		)
	}
	if r.Operation == Request_Start {
		t.req = *r
	}

	if t.req.GoogleCloudStorageKey != "" {
		if err := toFile(t.req.GoogleCloudStorageKey, filepath.Join(globalFlags.WorkingDirectory, "gcloud-key.json")); err != nil {
			return nil, err
		}
	}

	var processPID int
	switch r.Operation {
	case Request_Start:
		switch t.req.Database {
		case Request_etcdv2, Request_etcdv3:
			_, err := os.Stat(etcdBinaryPath)
			if err != nil {
				return nil, err
			}
			if err := os.RemoveAll(etcdDataDir); err != nil {
				return nil, err
			}
			f, err := openToAppend(databaseLogPath)
			if err != nil {
				return nil, err
			}
			t.logfile = f

			clusterN := len(peerIPs)
			names := make([]string, clusterN)
			clientURLs := make([]string, clusterN)
			peerURLs := make([]string, clusterN)
			members := make([]string, clusterN)
			for i, u := range peerIPs {
				names[i] = fmt.Sprintf("etcd-%d", i+1)
				clientURLs[i] = fmt.Sprintf("http://%s:2379", u)
				peerURLs[i] = fmt.Sprintf("http://%s:2380", u)
				members[i] = fmt.Sprintf("%s=%s", names[i], peerURLs[i])
			}
			clusterStr := strings.Join(members, ",")
			flags := []string{
				"--name", names[t.req.ServerIndex],
				"--data-dir", etcdDataDir,

				"--listen-client-urls", clientURLs[t.req.ServerIndex],
				"--advertise-client-urls", clientURLs[t.req.ServerIndex],

				"--listen-peer-urls", peerURLs[t.req.ServerIndex],
				"--initial-advertise-peer-urls", peerURLs[t.req.ServerIndex],

				"--initial-cluster-token", etcdToken,
				"--initial-cluster", clusterStr,
				"--initial-cluster-state", "new",
			}
			// if t.req.EtcdCompression != "" {
			// 	if compress.ParseType(t.req.EtcdCompression) != compress.NoCompress {
			// 		flags = append(flags,
			// 			"--experimental-compression", t.req.EtcdCompression,
			// 		)
			// 	}
			// }
			flagString := strings.Join(flags, " ")

			cmd := exec.Command(etcdBinaryPath, flags...)
			cmd.Stdout = f
			cmd.Stderr = f

			cmdString := fmt.Sprintf("%s %s", cmd.Path, flagString)
			logger.Info("starting binary", zap.String("command", cmdString))
			if err := cmd.Start(); err != nil {
				return nil, err
			}
			t.cmd = cmd
			t.pid = cmd.Process.Pid
			logger.Info("started binary", zap.String("command", cmdString), zap.Int("pid", t.pid))
			processPID = t.pid
			go func() {
				if err := cmd.Wait(); err != nil {
					logger.Error("cmd.Wait returned error", zap.String("command", cmdString), zap.Err(err))
					return
				}
				logger.Info("exiting", zap.String("command", cmdString))
			}()

		case Request_ZooKeeper:
			_, err := os.Stat(javaBinaryPath)
			if err != nil {
				return nil, err
			}

			logger.Info("os.Chdir", zap.String("path", zkWorkingDir))
			if err := os.Chdir(zkWorkingDir); err != nil {
				return nil, err
			}

			logger.Info("os.MkdirAll", zap.String("path", zkDataDir))
			if err := os.MkdirAll(zkDataDir, 0777); err != nil {
				return nil, err
			}

			idFilePath := filepath.Join(zkDataDir, "myid")
			logger.Info("writing zk myid file", zap.Int("myid", int(t.req.ZookeeperMyID)), zap.String("path", idFilePath))
			if err := toFile(fmt.Sprintf("%d", t.req.ZookeeperMyID), idFilePath); err != nil {
				return nil, err
			}

			// generate zookeeper config
			zkCfg := zkConfigDefault
			zkCfg.DataDir = zkDataDir
			peers := []ZookeeperPeer{}
			for i := range peerIPs {
				peers = append(peers, ZookeeperPeer{MyID: i + 1, IP: peerIPs[i]})
			}
			zkCfg.Peers = peers
			zkCfg.MaxClientCnxns = t.req.ZookeeperMaxClientCnxns
			zkCfg.SnapCount = t.req.ZookeeperSnapCount
			tpl := template.Must(template.New("zkTemplate").Parse(zkTemplate))
			buf := new(bytes.Buffer)
			if err := tpl.Execute(buf, zkCfg); err != nil {
				return nil, err
			}
			zc := buf.String()

			configFilePath := filepath.Join(zkWorkingDir, zkConfigPath)
			logger.Info("writing zk config file", zap.String("path", configFilePath), zap.String("config", zc))
			if err := toFile(zc, configFilePath); err != nil {
				return nil, err
			}

			f, err := openToAppend(databaseLogPath)
			if err != nil {
				return nil, err
			}
			t.logfile = f

			// this changes for different releases
			flagString := `-cp zookeeper-3.4.8.jar:lib/slf4j-api-1.6.1.jar:lib/slf4j-log4j12-1.6.1.jar:lib/log4j-1.2.16.jar:conf org.apache.zookeeper.server.quorum.QuorumPeerMain`
			args := []string{shell, "-c", javaBinaryPath + " " + flagString + " " + configFilePath}

			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = f
			cmd.Stderr = f

			cmdString := fmt.Sprintf("%s %s", cmd.Path, strings.Join(args[1:], " "))
			logger.Info("starting binary", zap.String("command", cmdString))
			if err := cmd.Start(); err != nil {
				return nil, err
			}
			t.cmd = cmd
			t.pid = cmd.Process.Pid
			logger.Info("started binary", zap.Int("pid", t.pid), zap.String("command", cmdString))
			processPID = t.pid
			go func() {
				if err := cmd.Wait(); err != nil {
					logger.Error("cmd.Wait returned error", zap.String("command", cmdString), zap.Err(err))
					return
				}
				logger.Info("exiting", zap.String("command", cmdString), zap.Err(err))
			}()

		case Request_Consul:
			_, err := os.Stat(consulBinaryPath)
			if err != nil {
				return nil, err
			}
			if err := os.RemoveAll(consulDataDir); err != nil {
				return nil, err
			}
			f, err := openToAppend(databaseLogPath)
			if err != nil {
				return nil, err
			}
			t.logfile = f

			var flags []string
			switch t.req.ServerIndex {
			case 0: // leader
				flags = []string{
					"agent",
					"-server",
					"-data-dir", consulDataDir,
					"-bind", peerIPs[t.req.ServerIndex],
					"-client", peerIPs[t.req.ServerIndex],
					"-bootstrap-expect", "3",
				}

			default:
				flags = []string{
					"agent",
					"-server",
					"-data-dir", consulDataDir,
					"-bind", peerIPs[t.req.ServerIndex],
					"-client", peerIPs[t.req.ServerIndex],
					"-join", peerIPs[0],
				}
			}
			flagString := strings.Join(flags, " ")

			cmd := exec.Command(consulBinaryPath, flags...)
			cmd.Stdout = f
			cmd.Stderr = f

			cmdString := fmt.Sprintf("%s %s", cmd.Path, flagString)
			logger.Info("starting binary", zap.String("command", cmdString))
			if err := cmd.Start(); err != nil {
				return nil, err
			}
			t.cmd = cmd
			t.pid = cmd.Process.Pid
			logger.Info("started binary", zap.Int("pid", t.pid), zap.String("command", cmdString))
			processPID = t.pid
			go func() {
				if err := cmd.Wait(); err != nil {
					logger.Error("cmd.Wait returned error", zap.String("command", cmdString), zap.Err(err))
					return
				}
				logger.Info("exiting", zap.String("command", cmdString), zap.Err(err))
			}()

		default:
			return nil, fmt.Errorf("unknown database (%q)", r.Database)
		}

	case Request_Restart:
		if t.cmd == nil {
			return nil, fmt.Errorf("nil command")
		}

		logger.Info("restarting database", zap.String("database", t.req.Database.String()))
		if r.Database == Request_ZooKeeper {
			if err := os.Chdir(zkWorkingDir); err != nil {
				return nil, err
			}
		}

		f, err := openToAppend(databaseLogPath)
		if err != nil {
			return nil, err
		}
		t.logfile = f

		cmd := exec.Command(t.cmd.Path, t.cmd.Args[1:]...)
		cmd.Stdout = f
		cmd.Stderr = f

		cmdString := strings.Join(t.cmd.Args, " ")
		logger.Info("restarting binary", zap.String("command", cmdString))
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		t.cmd = cmd
		t.pid = cmd.Process.Pid
		logger.Info("restarted binary", zap.Int("pid", t.pid), zap.String("command", cmdString))
		processPID = t.pid
		go func() {
			if err := cmd.Wait(); err != nil {
				logger.Error("cmd.Wait returned error", zap.String("command", cmdString), zap.Err(err))
				return
			}
			logger.Info("exiting", zap.String("command", cmdString))
		}()

	case Request_Stop:
		time.Sleep(3 * time.Second) // wait a few more seconds to collect more monitoring data
		if t.cmd == nil {
			return nil, fmt.Errorf("nil command")
		}
		logger.Info("stopping binary", zap.String("database", t.req.Database.String()), zap.Int("pid", t.pid))
		if err := syscall.Kill(t.pid, syscall.SIGTERM); err != nil {
			return nil, err
		}
		if t.logfile != nil {
			t.logfile.Close()
		}
		logger.Info("stopped binary", zap.String("database", t.req.Database.String()), zap.Int("pid", t.pid))
		processPID = t.pid
		databaseStopped <- struct{}{}

	default:
		return nil, fmt.Errorf("Not implemented %v", r.Operation)
	}

	if r.Operation == Request_Start || r.Operation == Request_Restart {
		go func(processPID int) {
			notifier := make(chan os.Signal, 1)
			signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

			rFunc := func() error {
				pss, err := ps.List(&ps.Process{Stat: ps.Stat{Pid: int64(processPID)}})
				if err != nil {
					return err
				}

				f, err := openToAppend(monitorLogPath)
				if err != nil {
					return err
				}
				defer f.Close()

				return ps.WriteToCSV(f, pss...)
			}

			logger.Info("saving monitoring results", zap.String("database", t.req.Database.String()), zap.String("path", monitorLogPath))
			var err error
			if err = rFunc(); err != nil {
				logger.Error("monitoring error", zap.Err(err))
				return
			}

			for {
				select {
				case <-time.After(time.Second):
					if err = rFunc(); err != nil {
						logger.Error("monitoring error", zap.Err(err))
						continue
					}

				case sig := <-notifier:
					logger.Info("signal received", zap.String("signal", sig.String()))
					return

				case <-databaseStopped:
					logger.Info("stopped monitoring, uploading to storage", zap.String("name", t.req.GoogleCloudProjectName))
					u, err := remotestorage.NewGoogleCloudStorage([]byte(t.req.GoogleCloudStorageKey), t.req.GoogleCloudProjectName)
					if err != nil {
						logger.Error("remotestorage.NewGoogleCloudStorage error", zap.Err(err))
						return
					}

					// set up file names
					srcDatabaseLogPath := databaseLogPath
					dstDatabaseLogPath := filepath.Base(databaseLogPath)
					if !strings.HasPrefix(filepath.Base(databaseLogPath), t.req.TestName) {
						dstDatabaseLogPath = fmt.Sprintf("%s-%d-%s", t.req.TestName, t.req.ServerIndex+1, filepath.Base(databaseLogPath))
					}
					dstDatabaseLogPath = filepath.Join(t.req.GoogleCloudStorageSubDirectory, dstDatabaseLogPath)

					logger.Info("uploading database log", zap.String("src", srcDatabaseLogPath), zap.String("dst", dstDatabaseLogPath))
					var uerr error
					for k := 0; k < 30; k++ {
						if uerr = u.UploadFile(t.req.GoogleCloudStorageBucketName, srcDatabaseLogPath, dstDatabaseLogPath); uerr != nil {
							logger.Error("u.UploadFile error... sleep and retry...", zap.Err(uerr))
							time.Sleep(2 * time.Second)
							continue
						} else {
							break
						}
					}

					srcMonitorResultPath := monitorLogPath
					dstMonitorResultPath := filepath.Base(monitorLogPath)
					if !strings.HasPrefix(filepath.Base(monitorLogPath), t.req.TestName) {
						dstMonitorResultPath = fmt.Sprintf("%s-%d-%s", t.req.TestName, t.req.ServerIndex+1, filepath.Base(monitorLogPath))
					}
					dstMonitorResultPath = filepath.Join(t.req.GoogleCloudStorageSubDirectory, dstMonitorResultPath)

					logger.Info("uploading monitor results", zap.String("src", srcMonitorResultPath), zap.String("dst", dstMonitorResultPath))
					for k := 0; k < 30; k++ {
						if uerr = u.UploadFile(t.req.GoogleCloudStorageBucketName, srcMonitorResultPath, dstMonitorResultPath); uerr != nil {
							logger.Error("u.UploadFile error... sleep and retry...", zap.Err(uerr))
							time.Sleep(2 * time.Second)
							continue
						} else {
							break
						}
					}

					srcAgentLogPath := agentLogPath
					dstAgentLogPath := filepath.Base(agentLogPath)
					if !strings.HasPrefix(filepath.Base(agentLogPath), t.req.TestName) {
						dstAgentLogPath = fmt.Sprintf("%s-%d-%s", t.req.TestName, t.req.ServerIndex+1, filepath.Base(agentLogPath))
					}
					dstAgentLogPath = filepath.Join(t.req.GoogleCloudStorageSubDirectory, dstAgentLogPath)

					logger.Info("uploading agent logs", zap.String("src", srcAgentLogPath), zap.String("dst", dstAgentLogPath))
					for k := 0; k < 30; k++ {
						if uerr = u.UploadFile(t.req.GoogleCloudStorageBucketName, srcAgentLogPath, dstAgentLogPath); uerr != nil {
							logger.Error("u.UploadFile error... sleep and retry...", zap.Err(uerr))
							time.Sleep(2 * time.Second)
							continue
						} else {
							break
						}
					}

					return
				}
			}
		}(processPID)
	}

	logger.Info("transfer success", zap.Time("time", time.Now()))
	return &Response{Success: true}, nil
}

func NewTransporterServer() TransporterServer {
	return &transporterServer{}
}
