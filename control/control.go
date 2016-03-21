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

package control

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/dbtester/agent"
	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type (
	Flags struct {
		Database                string
		AgentEndpoints          []string
		ZookeeperPreAllocSize   int64
		ZookeeperMaxClientCnxns int64
	}
)

var (
	StartCommand = &cobra.Command{
		Use:   "start",
		Short: "Starts database through RPC calls.",
		Run:   CommandFunc,
	}
	StopCommand = &cobra.Command{
		Use:   "stop",
		Short: "Stops database through RPC calls.",
		Run:   CommandFunc,
	}
	RestartCommand = &cobra.Command{
		Use:   "restart",
		Short: "Restarts database through RPC calls.",
		Run:   CommandFunc,
	}
	globalFlags = Flags{}
)

func init() {
	StartCommand.PersistentFlags().StringVar(&globalFlags.Database, "database", "", "etcd, zookeeper, zk, etcd2, consul.")
	StartCommand.PersistentFlags().StringSliceVar(&globalFlags.AgentEndpoints, "agent-endpoints", []string{""}, "Endpoints to send client requests to, then it automatically configures.")
	StartCommand.PersistentFlags().Int64Var(&globalFlags.ZookeeperPreAllocSize, "zk-pre-alloc-size", 65536*1024, "Disk pre-allocation size in bytes.")
	StartCommand.PersistentFlags().Int64Var(&globalFlags.ZookeeperMaxClientCnxns, "zk-max-client-conns", 5000, "Maximum number of concurrent Zookeeper connection.")

	StopCommand.PersistentFlags().StringSliceVar(&globalFlags.AgentEndpoints, "agent-endpoints", []string{""}, "Endpoints to send client requests to, then it automatically configures.")

	RestartCommand.PersistentFlags().StringSliceVar(&globalFlags.AgentEndpoints, "agent-endpoints", []string{""}, "Endpoints to send client requests to, then it automatically configures.")
}

func CommandFunc(cmd *cobra.Command, args []string) {
	if globalFlags.Database == "zk" {
		globalFlags.Database = "zookeeper"
	}
	req := agent.Request{}

	switch cmd.Use {
	case "start":
		req.Operation = agent.Request_Start
	case "stop":
		req.Operation = agent.Request_Stop
	case "restart":
		req.Operation = agent.Request_Restart
	default:
		log.Printf("Operation '%s' is not supported!\n", cmd.Use)
		os.Exit(-1)
	}

	switch globalFlags.Database {
	case "etcd":
		req.Database = agent.Request_etcd
	case "zookeeper":
		req.Database = agent.Request_ZooKeeper
	default:
		log.Printf("'%s' is not supported!\n", globalFlags.Database)
		os.Exit(-1)
	}
	peerIPs := extractIPs(globalFlags.AgentEndpoints)
	req.PeerIPs = strings.Join(peerIPs, "___") // because protoc mixes the order of 'repeated' type data

	if cmd.Use == "start" {
		req.ZookeeperPreAllocSize = globalFlags.ZookeeperPreAllocSize
		req.ZookeeperMaxClientCnxns = globalFlags.ZookeeperMaxClientCnxns
	}

	for i := range peerIPs {
		nreq := req

		nreq.EtcdServerIndex = uint32(i)
		nreq.ZookeeperMyID = uint32(i + 1)
		ep := globalFlags.AgentEndpoints[nreq.EtcdServerIndex]

		log.Printf("%s to %s\n", req.Database, ep)
		conn, err := grpc.Dial(ep, grpc.WithInsecure())
		if err != nil {
			log.Printf("error %v when connecting to %s\n", err, ep)
			os.Exit(-1)
		}
		defer conn.Close()

		cli := agent.NewTransporterClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := cli.Transfer(ctx, &nreq)
		cancel()
		if err != nil {
			log.Printf("error %v when transferring to %s\n", err, ep)
			return
		}
		log.Printf("Response from %s (%+v)\n", ep, resp)
	}
}

func extractIPs(es []string) []string {
	var rs []string
	for _, v := range es {
		sl := strings.Split(v, ":")
		rs = append(rs, sl[0])
	}
	return rs
}
