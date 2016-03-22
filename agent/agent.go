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

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/dbtester/remotestorage"
	"github.com/gyuho/psn/ps"
	"github.com/spf13/cobra"
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
		PreAllocSize   int64
		MaxClientCnxns int64
		Peers          []ZookeeperPeer
	}
	ZookeeperPeer struct {
		MyID int
		IP   string
	}
)

var (
	shell = os.Getenv("SHELL")

	agentLogPath = "agent.log"

	etcdBinaryPath = filepath.Join(os.Getenv("GOPATH"), "bin/etcd")
	etcdToken      = "etcd_token"

	etcdDataDir  = "data.etcd"
	zkWorkingDir = "zookeeper"
	zkDataDir    = "zookeeper/data.zk"

	zkTemplate = `tickTime={{.TickTime}}
dataDir={{.DataDir}}
clientPort={{.ClientPort}}
initLimit={{.InitLimit}}
syncLimit={{.SyncLimit}}
preAllocSize={{.PreAllocSize}}
maxClientCnxns={{.MaxClientCnxns}}
{{range .Peers}}server.{{.MyID}}={{.IP}}:2888:3888
{{end}}
`
	zkConfig = ZookeeperConfig{
		TickTime:       2000,
		ClientPort:     "2181",
		InitLimit:      5,
		SyncLimit:      5,
		PreAllocSize:   65536 * 1024,
		MaxClientCnxns: 60,
		Peers: []ZookeeperPeer{
			{MyID: 1, IP: "10.240.0.12"},
			{MyID: 2, IP: "10.240.0.13"},
			{MyID: 3, IP: "10.240.0.14"},
		},
	}

	Command = &cobra.Command{
		Use:   "agent",
		Short: "Database agent in remote servers.",
		Run:   CommandFunc,
	}
	globalFlags = Flags{}
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	if len(shell) == 0 {
		shell = "sh"
	}
	Command.PersistentFlags().StringVar(&globalFlags.GRPCPort, "agent-port", ":3500", "Port to server agent gRPC server.")
	Command.PersistentFlags().StringVar(&globalFlags.WorkingDirectory, "working-directory", homeDir(), "Working directory.")
}

func CommandFunc(cmd *cobra.Command, args []string) {
	if !exist(globalFlags.WorkingDirectory) {
		log.Fatalf("%s does not exist", globalFlags.WorkingDirectory)
	}
	if !filepath.HasPrefix(agentLogPath, globalFlags.WorkingDirectory) {
		agentLogPath = filepath.Join(globalFlags.WorkingDirectory, agentLogPath)
	}

	f, err := openToAppend(agentLogPath)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("gRPC serving %s", globalFlags.GRPCPort)
	var (
		grpcServer = grpc.NewServer()
		sender     = NewTransporterServer()
	)
	ln, err := net.Listen("tcp", globalFlags.GRPCPort)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	RegisterTransporterServer(grpcServer, sender)

	if err := grpcServer.Serve(ln); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

type transporterServer struct { // satisfy TransporterServer
	req     Request
	cmd     *exec.Cmd
	logfile *os.File
	pid     int
}

var databaseStopped = make(chan struct{})

func (t *transporterServer) Transfer(ctx context.Context, r *Request) (*Response, error) {
	peerIPs := strings.Split(r.PeerIPs, "___")
	if r.Operation == Request_Start || r.Operation == Request_Restart {
		if !filepath.HasPrefix(etcdDataDir, globalFlags.WorkingDirectory) {
			etcdDataDir = filepath.Join(globalFlags.WorkingDirectory, etcdDataDir)
		}
		if !filepath.HasPrefix(zkWorkingDir, globalFlags.WorkingDirectory) {
			zkWorkingDir = filepath.Join(globalFlags.WorkingDirectory, zkWorkingDir)
		}
		if !filepath.HasPrefix(zkDataDir, globalFlags.WorkingDirectory) {
			zkDataDir = filepath.Join(globalFlags.WorkingDirectory, zkDataDir)
		}
		if r.LogPrefix != "" {
			if !strings.HasPrefix(filepath.Base(r.DatabaseLogPath), r.LogPrefix) {
				r.DatabaseLogPath = filepath.Join(filepath.Dir(r.DatabaseLogPath), r.LogPrefix+"_"+filepath.Base(r.DatabaseLogPath))
			}
			if !strings.HasPrefix(filepath.Base(r.MonitorResultPath), r.LogPrefix) {
				r.MonitorResultPath = filepath.Join(filepath.Dir(r.MonitorResultPath), r.LogPrefix+"_"+filepath.Base(r.MonitorResultPath))
			}
		}
		if !filepath.HasPrefix(r.DatabaseLogPath, globalFlags.WorkingDirectory) {
			r.DatabaseLogPath = filepath.Join(globalFlags.WorkingDirectory, r.DatabaseLogPath)
		}
		if !filepath.HasPrefix(r.MonitorResultPath, globalFlags.WorkingDirectory) {
			r.MonitorResultPath = filepath.Join(globalFlags.WorkingDirectory, r.MonitorResultPath)
		}

		log.Printf("working directory: %s", globalFlags.WorkingDirectory)
		log.Printf("etcd data directory: %s", etcdDataDir)
		log.Printf("zookeeper working directory: %s", zkWorkingDir)
		log.Printf("zookeeper data directory: %s", zkDataDir)
		log.Printf("database log path: %s", r.DatabaseLogPath)
		log.Printf("monitor result path: %s", r.MonitorResultPath)
	}
	t.req = *r

	var processPID int
	switch r.Operation {
	case Request_Start:
		switch r.Database {
		case Request_etcd:
			_, err := os.Stat(etcdBinaryPath)
			if err != nil {
				return nil, err
			}

			if err := os.RemoveAll(etcdDataDir); err != nil {
				return nil, err
			}

			f, err := openToAppend(r.DatabaseLogPath)
			if err != nil {
				return nil, err
			}
			t.logfile = f

			// generate flags from etcd server name
			clusterN := len(peerIPs)
			names := make([]string, clusterN)
			clientURLs := make([]string, clusterN)
			peerURLs := make([]string, clusterN)
			members := make([]string, clusterN)
			for i, u := range peerIPs {
				names[i] = fmt.Sprintf("etcd-%d", i)
				clientURLs[i] = fmt.Sprintf("http://%s:2379", u)
				peerURLs[i] = fmt.Sprintf("http://%s:2380", u)
				members[i] = fmt.Sprintf("%s=%s", names[i], peerURLs[i])
			}
			clusterStr := strings.Join(members, ",")
			flags := []string{
				"--name", fmt.Sprintf("etcd-%d", r.EtcdServerIndex),
				"--data-dir", etcdDataDir,

				"--listen-client-urls", clientURLs[r.EtcdServerIndex],
				"--advertise-client-urls", clientURLs[r.EtcdServerIndex],

				"--listen-peer-urls", peerURLs[r.EtcdServerIndex],
				"--initial-advertise-peer-urls", peerURLs[r.EtcdServerIndex],

				"--initial-cluster-token", etcdToken,
				"--initial-cluster", clusterStr,
				"--initial-cluster-state", "new",

				"--experimental-v3demo",
			}
			flagString := strings.Join(flags, " ")

			cmd := exec.Command(etcdBinaryPath, flags...)
			cmd.Stdout = f
			cmd.Stderr = f
			log.Printf("Starting: %s %s", cmd.Path, flagString)
			if err := cmd.Start(); err != nil {
				return nil, err
			}
			t.cmd = cmd
			t.pid = cmd.Process.Pid
			log.Printf("Started: %s [PID: %d]", cmd.Path, t.pid)
			processPID = t.pid
			go func() {
				if err := cmd.Wait(); err != nil {
					log.Printf("Start(%s) cmd.Wait returned %v", cmd.Path, err)
					return
				}
				log.Printf("Exiting %s", cmd.Path)
			}()

		case Request_ZooKeeper:
			_, err := os.Stat("/usr/bin/java")
			if err != nil {
				return nil, err
			}

			log.Printf("os.Chdir: %s", zkWorkingDir)
			if err := os.Chdir(zkWorkingDir); err != nil {
				return nil, err
			}

			log.Printf("os.MkdirAll: %s", zkDataDir)
			if err := os.MkdirAll(zkDataDir, 0777); err != nil {
				return nil, err
			}

			idFilePath := filepath.Join(zkDataDir, "myid")
			log.Printf("Writing %d to %s", r.ZookeeperMyID, idFilePath)
			if err := toFile(fmt.Sprintf("%d", r.ZookeeperMyID), idFilePath); err != nil {
				return nil, err
			}

			// generate zookeeper config
			zkCfg := zkConfig
			zkCfg.DataDir = zkDataDir
			peers := []ZookeeperPeer{}
			for i := range peerIPs {
				peers = append(peers, ZookeeperPeer{MyID: i + 1, IP: peerIPs[i]})
			}
			zkCfg.Peers = peers
			zkCfg.PreAllocSize = r.ZookeeperPreAllocSize
			zkCfg.MaxClientCnxns = r.ZookeeperMaxClientCnxns
			tpl := template.Must(template.New("zkTemplate").Parse(zkTemplate))
			buf := new(bytes.Buffer)
			if err := tpl.Execute(buf, zkCfg); err != nil {
				return nil, err
			}
			zc := buf.String()

			configFilePath := filepath.Join(zkWorkingDir, "zookeeper.config")
			log.Printf("Writing %q to %s", zc, configFilePath)
			if err := toFile(zc, configFilePath); err != nil {
				return nil, err
			}

			f, err := openToAppend(r.DatabaseLogPath)
			if err != nil {
				return nil, err
			}
			t.logfile = f

			// this changes for different releases
			flagString := `-cp zookeeper-3.4.8.jar:lib/slf4j-api-1.6.1.jar:lib/slf4j-log4j12-1.6.1.jar:lib/log4j-1.2.16.jar:conf org.apache.zookeeper.server.quorum.QuorumPeerMain`
			args := []string{shell, "-c", "/usr/bin/java " + flagString + " " + configFilePath}

			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = f
			cmd.Stderr = f
			log.Printf("Starting: %s %s", cmd.Path, strings.Join(args[1:], " "))
			if err := cmd.Start(); err != nil {
				return nil, err
			}
			t.cmd = cmd
			t.pid = cmd.Process.Pid
			log.Printf("Started: %s [PID: %d]", cmd.Path, t.pid)
			processPID = t.pid
			go func() {
				if err := cmd.Wait(); err != nil {
					log.Printf("Start(%s) cmd.Wait returned %v", cmd.Path, err)
					return
				}
				log.Printf("Exiting %s", cmd.Path)
			}()

		default:
			return nil, fmt.Errorf("unknown database (%q)", r.Database)
		}

	case Request_Restart:
		if t.cmd == nil {
			return nil, fmt.Errorf("nil command")
		}

		log.Printf("Restarting %s", t.req.Database)
		if r.Database == Request_ZooKeeper {
			if err := os.Chdir(zkWorkingDir); err != nil {
				return nil, err
			}
		}

		f, err := openToAppend(r.DatabaseLogPath)
		if err != nil {
			return nil, err
		}
		t.logfile = f

		cmd := exec.Command(t.cmd.Path, t.cmd.Args[1:]...)
		cmd.Stdout = f
		cmd.Stderr = f
		log.Printf("Restarting: %s", strings.Join(t.cmd.Args, " "))
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		t.cmd = cmd
		t.pid = cmd.Process.Pid
		log.Printf("Restarted: %s [PID: %d]", cmd.Path, t.pid)
		processPID = t.pid
		go func() {
			if err := cmd.Wait(); err != nil {
				log.Printf("Restart(%s) cmd.Wait returned %v", cmd.Path, err)
				return
			}
			log.Printf("Exiting %s", cmd.Path)
		}()

	case Request_Stop:
		if t.cmd == nil {
			return nil, fmt.Errorf("nil command")
		}
		log.Printf("Stopping %s [PID: %d]", t.req.Database, t.pid)
		ps.Kill(os.Stdout, false, ps.Process{Stat: ps.Stat{Pid: int64(t.pid)}})
		if t.logfile != nil {
			t.logfile.Close()
		}
		log.Printf("Stopped: %s [PID: %d]", t.cmd.Path, t.pid)
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

				f, err := openToAppend(r.MonitorResultPath)
				if err != nil {
					return err
				}
				defer f.Close()

				return ps.WriteToCSV(f, pss...)
			}

			log.Printf("%s monitor saved at %s", r.Database, r.MonitorResultPath)
			var err error
			if err = rFunc(); err != nil {
				log.Warningln("error:", err)
				return
			}

		escape:
			for {
				select {
				case <-time.After(time.Second):
					if err = rFunc(); err != nil {
						log.Warnf("Monitoring error %v", err)
						break escape
					}
				case sig := <-notifier:
					log.Printf("Received %v", sig)
					return
				case <-databaseStopped:
					log.Println("Monitoring stopped. Uploading data to cloud storage...")

					if err := toFile(r.GoogleCloudStorageJSONKey, filepath.Join(globalFlags.WorkingDirectory, "key.json")); err != nil {
						log.Warnf("error (%v)", err)
						return
					}
					u, err := remotestorage.NewGoogleCloudStorage([]byte(r.GoogleCloudStorageJSONKey), r.GoogleCloudProjectName)
					if err != nil {
						log.Warnf("error (%v)", err)
						return
					}

					// set up file names
					srcDatabaseLogPath := r.DatabaseLogPath
					dstDatabaseLogPath := filepath.Base(r.DatabaseLogPath)
					if !strings.HasPrefix(filepath.Base(r.DatabaseLogPath), r.LogPrefix) {
						dstDatabaseLogPath = r.LogPrefix + fmt.Sprintf("_%d_", r.EtcdServerIndex) + filepath.Base(r.DatabaseLogPath)
					}

					srcMonitorResultPath := r.MonitorResultPath
					dstMonitorResultPath := filepath.Base(r.MonitorResultPath)
					if !strings.HasPrefix(filepath.Base(r.MonitorResultPath), r.LogPrefix) {
						dstMonitorResultPath = r.LogPrefix + fmt.Sprintf("_%d_", r.EtcdServerIndex) + filepath.Base(r.MonitorResultPath)
					}

					srcAgentLogPath := agentLogPath
					dstAgentLogPath := filepath.Base(agentLogPath)
					if !strings.HasPrefix(filepath.Base(agentLogPath), r.LogPrefix) {
						dstAgentLogPath = r.LogPrefix + fmt.Sprintf("_%d_", r.EtcdServerIndex) + filepath.Base(agentLogPath)
					}

					log.Printf("Uploading %s", srcDatabaseLogPath)
					if err := u.UploadFile(r.GoogleCloudStorageBucketName, srcDatabaseLogPath, dstDatabaseLogPath); err != nil {
						log.Fatal(err)
					}

					log.Printf("Uploading %s", srcMonitorResultPath)
					if err := u.UploadFile(r.GoogleCloudStorageBucketName, srcMonitorResultPath, dstMonitorResultPath); err != nil {
						log.Fatal(err)
					}

					log.Printf("Uploading %s", srcAgentLogPath)
					if err := u.UploadFile(r.GoogleCloudStorageBucketName, srcAgentLogPath, dstAgentLogPath); err != nil {
						log.Fatal(err)
					}

					return
				}
			}
		}(processPID)
	}

	log.Printf("Transfer success!")
	return &Response{Success: true}, nil
}

func NewTransporterServer() TransporterServer {
	return &transporterServer{}
}
