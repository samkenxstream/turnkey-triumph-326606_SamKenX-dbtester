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
	"io/ioutil"
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
	"github.com/gyuho/psn/ps"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
	"google.golang.org/grpc"
)

type (
	Flags struct {
		GRPCPort string
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

	agentLogPath = filepath.Join(homeDir(), "agent.log")

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
}

func CommandFunc(cmd *cobra.Command, args []string) {
	f, err := openToAppend(agentLogPath)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("gRPC serving: %s", globalFlags.GRPCPort)
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
	log.Printf("Message from %q", r.PeerIPs)
	peerIPs := strings.Split(r.PeerIPs, "___")

	if r.Operation == Request_Start || r.Operation == Request_Restart {
		if r.WorkingDirectory == "" {
			r.WorkingDirectory = homeDir()
		}
		if !exist(r.WorkingDirectory) {
			return nil, fmt.Errorf("%s does not exist", r.WorkingDirectory)
		}
		if !filepath.HasPrefix(etcdDataDir, r.WorkingDirectory) {
			etcdDataDir = filepath.Join(r.WorkingDirectory, etcdDataDir)
		}
		if !filepath.HasPrefix(zkWorkingDir, r.WorkingDirectory) {
			zkWorkingDir = filepath.Join(r.WorkingDirectory, zkWorkingDir)
		}
		if !filepath.HasPrefix(zkDataDir, r.WorkingDirectory) {
			zkDataDir = filepath.Join(r.WorkingDirectory, zkDataDir)
		}
		if r.LogPrefix != "" {
			if !strings.HasPrefix(filepath.Base(r.DatabaseLogPath), r.LogPrefix) {
				r.DatabaseLogPath = filepath.Join(filepath.Dir(r.DatabaseLogPath), r.LogPrefix+"_"+filepath.Base(r.DatabaseLogPath))
			}
			if !strings.HasPrefix(filepath.Base(r.MonitorResultPath), r.LogPrefix) {
				r.MonitorResultPath = filepath.Join(filepath.Dir(r.MonitorResultPath), r.LogPrefix+"_"+filepath.Base(r.MonitorResultPath))
			}
		} else {
			if !filepath.HasPrefix(r.DatabaseLogPath, r.WorkingDirectory) {
				r.DatabaseLogPath = filepath.Join(r.WorkingDirectory, r.DatabaseLogPath)
			}
			if !filepath.HasPrefix(r.MonitorResultPath, r.WorkingDirectory) {
				r.MonitorResultPath = filepath.Join(r.WorkingDirectory, r.MonitorResultPath)
			}
		}
		log.Printf("Working directory: %s", r.WorkingDirectory)
		log.Printf("etcd data directory: %s", etcdDataDir)
		log.Printf("Zookeeper working directory: %s", zkWorkingDir)
		log.Printf("Zookeeper data directory: %s", zkDataDir)
		log.Printf("Database log path: %s", r.DatabaseLogPath)
		log.Printf("Monitor result path: %s", r.MonitorResultPath)
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
			grpcURLs := make([]string, clusterN)
			clientURLs := make([]string, clusterN)
			peerURLs := make([]string, clusterN)
			members := make([]string, clusterN)
			for i, u := range peerIPs {
				names[i] = fmt.Sprintf("etcd-%d", i)
				grpcURLs[i] = fmt.Sprintf("%s:2378", u)
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
				"--experimental-gRPC-addr", grpcURLs[r.EtcdServerIndex],
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

					// initialize auth
					conf, err := google.JWTConfigFromJSON(
						[]byte(r.GoogleCloudStorageJSONKey),
						storage.ScopeFullControl,
					)
					if err != nil {
						log.Warnf("error (%v) with %q", err, r.GoogleCloudStorageJSONKey)
						return
					}
					ctx := context.Background()
					aclient, err := storage.NewAdminClient(ctx, r.GoogleCloudProjectName, cloud.WithTokenSource(conf.TokenSource(ctx)))
					if err != nil {
						log.Warnf("error (%v) with %q", err, r.GoogleCloudProjectName)
					}
					defer aclient.Close()

					if err := aclient.CreateBucket(context.Background(), r.GoogleCloudStorageBucketName, nil); err != nil {
						if !strings.Contains(err.Error(), "You already own this bucket. Please select another name") {
							log.Warnf("error (%v) with %q", err, r.GoogleCloudStorageBucketName)
						}
					}

					sctx := context.Background()
					sclient, err := storage.NewClient(sctx, cloud.WithTokenSource(conf.TokenSource(sctx)))
					if err != nil {
						log.Warnf("error (%v)", err)
						return
					}
					defer sclient.Close()

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
					wc1 := sclient.Bucket(r.GoogleCloudStorageBucketName).Object(dstDatabaseLogPath).NewWriter(context.Background())
					wc1.ContentType = "text/plain"
					bts1, err := ioutil.ReadFile(srcDatabaseLogPath)
					if err != nil {
						log.Warnf("error (%v)", err)
						return
					}
					if _, err := wc1.Write(bts1); err != nil {
						log.Warnf("error (%v)", err)
						return
					}
					if err := wc1.Close(); err != nil {
						log.Warnf("error (%v)", err)
						return
					}

					log.Printf("Uploading %s", srcMonitorResultPath)
					wc2 := sclient.Bucket(r.GoogleCloudStorageBucketName).Object(dstMonitorResultPath).NewWriter(context.Background())
					wc2.ContentType = "text/plain"
					bts2, err := ioutil.ReadFile(srcMonitorResultPath)
					if err != nil {
						log.Warnf("error (%v)", err)
						return
					}
					if _, err := wc2.Write(bts2); err != nil {
						log.Warnf("error (%v)", err)
						return
					}
					if err := wc2.Close(); err != nil {
						log.Warnf("error (%v)", err)
						return
					}

					log.Printf("Uploading %s", srcAgentLogPath)
					wc3 := sclient.Bucket(r.GoogleCloudStorageBucketName).Object(dstAgentLogPath).NewWriter(context.Background())
					wc3.ContentType = "text/plain"
					bts3, err := ioutil.ReadFile(srcAgentLogPath)
					if err != nil {
						log.Warnf("error (%v)", err)
						return
					}
					if _, err := wc3.Write(bts3); err != nil {
						log.Warnf("error (%v)", err)
						return
					}
					if err := wc3.Close(); err != nil {
						log.Warnf("error (%v)", err)
						return
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
