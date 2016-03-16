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
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/gyuho/psn/ps"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type (
	Flags struct {
		GRPCPort        string
		WorkingDir      string
		Monitor         bool
		MonitorInterval time.Duration
		MonitorFilePath string
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

	etcdBinaryPath = filepath.Join(os.Getenv("GOPATH"), "bin/etcd")
	etcdToken      = "etcd_token"

	dataLogPath  = "data.log"
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

func homeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func init() {
	if len(shell) == 0 {
		shell = "sh"
	}
	Command.PersistentFlags().StringVarP(&globalFlags.GRPCPort, "agent-port", "p", ":3500", "Port to server agent gRPC server.")
	Command.PersistentFlags().StringVarP(&globalFlags.WorkingDir, "working-dir", "d", homeDir(), "Working directory to store data.")
	Command.PersistentFlags().BoolVar(&globalFlags.Monitor, "monitor", false, "Periodically records resource usage.")
	Command.PersistentFlags().DurationVar(&globalFlags.MonitorInterval, "monitor-interval", time.Second, "Resource monitor interval.")
	Command.PersistentFlags().StringVar(&globalFlags.MonitorFilePath, "monitor-path", "monitor.csv", "File path to store monitor results.")
}

func CommandFunc(cmd *cobra.Command, args []string) {
	dataLogPath = filepath.Join(globalFlags.WorkingDir, dataLogPath)
	etcdDataDir = filepath.Join(globalFlags.WorkingDir, etcdDataDir)
	zkWorkingDir = filepath.Join(globalFlags.WorkingDir, zkWorkingDir)
	zkDataDir = filepath.Join(globalFlags.WorkingDir, zkDataDir)

	log.Printf("gRPC has started serving at  %s\n", globalFlags.GRPCPort)
	log.Printf("data log path:               %s\n", dataLogPath)
	log.Printf("etcd data directory:         %s\n", etcdDataDir)
	log.Printf("Zookeeper working directory: %s\n", zkWorkingDir)
	log.Printf("Zookeeper data directory:    %s\n", zkDataDir)

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

	log.Printf("gRPC is now serving at %s\n", globalFlags.GRPCPort)
	if globalFlags.Monitor {
		log.Printf(" As soon as database started, it will monitor every %v...\n", globalFlags.MonitorInterval)
	}
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

func (t *transporterServer) Transfer(ctx context.Context, r *Request) (*Response, error) {
	log.Printf("Received message for peer %q", r.PeerIPs)
	peerIPs := strings.Split(r.PeerIPs, "___")
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

			f, err := openToAppend(dataLogPath)
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

			// args := []string{shell, "-c", etcdBinaryPath + " " + flagString}
			// cmd := exec.Command(args[0], args[1:]...)

			cmd := exec.Command(etcdBinaryPath, flags...)
			cmd.Stdin = nil
			cmd.Stdout = f
			cmd.Stderr = f
			log.Printf("Starting: %s %s\n", cmd.Path, flagString)
			if err := cmd.Start(); err != nil {
				return nil, err
			}
			t.cmd = cmd
			t.pid = cmd.Process.Pid
			log.Printf("Started: %s [PID: %d]", cmd.Path, t.pid)
			processPID = t.pid
			go func() {
				if err := cmd.Wait(); err != nil {
					log.Printf("Start(%s) cmd.Wait returned %v\n", cmd.Path, err)
					return
				}
				log.Printf("Exiting %s\n", cmd.Path)
			}()

		case Request_ZooKeeper:
			_, err := os.Stat("/usr/bin/java")
			if err != nil {
				return nil, err
			}

			log.Printf("os.Chdir: %s\n", zkWorkingDir)
			if err := os.Chdir(zkWorkingDir); err != nil {
				return nil, err
			}

			log.Printf("os.MkdirAll: %s\n", zkDataDir)
			if err := os.MkdirAll(zkDataDir, 0777); err != nil {
				return nil, err
			}

			idFilePath := filepath.Join(zkDataDir, "myid")
			log.Printf("Writing %d to %s\n", r.ZookeeperMyID, idFilePath)
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
			log.Printf("Writing %q to %s\n", zc, configFilePath)
			if err := toFile(zc, configFilePath); err != nil {
				return nil, err
			}

			f, err := openToAppend(dataLogPath)
			if err != nil {
				return nil, err
			}
			t.logfile = f

			// this changes for different releases
			flagString := `-cp zookeeper-3.4.8.jar:lib/slf4j-api-1.6.1.jar:lib/slf4j-log4j12-1.6.1.jar:lib/log4j-1.2.16.jar:conf org.apache.zookeeper.server.quorum.QuorumPeerMain`
			args := []string{shell, "-c", "/usr/bin/java " + flagString + " " + configFilePath}

			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdin = nil
			cmd.Stdout = f
			cmd.Stderr = f
			log.Printf("Starting: %s %s\n", cmd.Path, strings.Join(args[1:], " "))
			if err := cmd.Start(); err != nil {
				return nil, err
			}
			t.cmd = cmd
			t.pid = cmd.Process.Pid
			log.Printf("Started: %s [PID: %d]", cmd.Path, t.pid)
			processPID = t.pid
			go func() {
				if err := cmd.Wait(); err != nil {
					log.Printf("Start(%s) cmd.Wait returned %v\n", cmd.Path, err)
					return
				}
				log.Printf("Exiting %s\n", cmd.Path)
			}()

		default:
			return nil, fmt.Errorf("unknown database (%q)", r.Database)
		}

	case Request_Restart:
		if t.cmd == nil {
			return nil, fmt.Errorf("nil command")
		}

		log.Printf("Restarting %s\n", t.req.Database)
		if r.Database == Request_ZooKeeper {
			if err := os.Chdir(zkWorkingDir); err != nil {
				return nil, err
			}
		}

		f, err := openToAppend(dataLogPath)
		if err != nil {
			return nil, err
		}
		t.logfile = f

		cmd := exec.Command(t.cmd.Path, t.cmd.Args[1:]...)
		cmd.Stdin = nil
		cmd.Stdout = f
		cmd.Stderr = f
		log.Printf("Restarting: %s\n", strings.Join(t.cmd.Args, " "))
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		t.cmd = cmd
		t.pid = cmd.Process.Pid
		log.Printf("Restarted: %s [PID: %d]", cmd.Path, t.pid)
		processPID = t.pid
		go func() {
			if err := cmd.Wait(); err != nil {
				log.Printf("Restart(%s) cmd.Wait returned %v\n", cmd.Path, err)
				return
			}
			log.Printf("Exiting %s\n", cmd.Path)
		}()

	case Request_Stop:
		if t.cmd == nil {
			return nil, fmt.Errorf("nil command")
		}
		log.Printf("Stopping %s [PID: %d]\n", t.req.Database, t.pid)
		if err := syscall.Kill(t.pid, syscall.SIGTERM); err != nil {
			log.Printf("%v\n", err)
		}
		time.Sleep(2 * time.Second)
		if err := syscall.Kill(t.pid, syscall.SIGKILL); err != nil {
			log.Printf("%v\n", err)
		}
		if t.logfile != nil {
			t.logfile.Close()
		}
		log.Printf("Stopped: %s [PID: %d\n]", t.cmd.Path, t.pid)
		processPID = t.pid

	default:
		return nil, fmt.Errorf("Not implemented %v", r.Operation)
	}

	firstCSV := true
	if globalFlags.Monitor && r.Operation == Request_Start {
		go func(processPID int) {
			notifier := make(chan os.Signal, 1)
			signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

			rFunc := func() error {
				pss, err := ps.ListStatus(&ps.Status{Pid: processPID})
				if err != nil {
					return err
				}

				if !strings.HasPrefix(globalFlags.MonitorFilePath, globalFlags.WorkingDir) {
					globalFlags.MonitorFilePath = filepath.Join(globalFlags.WorkingDir, globalFlags.MonitorFilePath)
				}
				f, err := openToAppend(globalFlags.MonitorFilePath)
				if err != nil {
					return err
				}
				defer f.Close()

				if err := ps.WriteToCSV(firstCSV, f, pss...); err != nil {
					return err
				}
				firstCSV = false
				return nil
			}

			log.Printf("Monitoring %v (file saved at %v)\n", r.Database, globalFlags.WorkingDir)
			var err error
			if err = rFunc(); err != nil {
				log.Println("error:", err)
			}

		escape:
			for {
				log.Printf("Monitoring %v at %v\n", r.Database, time.Now())
				select {
				case <-time.After(globalFlags.MonitorInterval):
					if err = rFunc(); err != nil {
						break escape
					}
				case sig := <-notifier:
					log.Printf("Received %v\n", sig)
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
