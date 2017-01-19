# dbtester

[![Build Status](https://img.shields.io/travis/coreos/dbtester.svg?style=flat-square)](https://travis-ci.org/coreos/dbtester) [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/coreos/dbtester)

Distributed database tester

- Database agent and runner are implemented at https://github.com/coreos/dbtester/tree/master/agent
- Client is implemented at https://github.com/coreos/dbtester/tree/master/control
- System metrics are collected via https://github.com/gyuho/psn
- Data analysis is done via https://github.com/coreos/dbtester/tree/master/analyze
  - https://github.com/gyuho/dataframe
  - https://github.com/gonum/plot

For etcd, we also recommend [etcd benchmark tool](https://github.com/coreos/etcd/tree/master/tools/benchmark).

<br><br><hr>
### Latest Results

All logs and results can be found at https://console.cloud.google.com/storage/browser/dbtester-results

- Google Cloud Compute Engine
- 3 machines of 8 vCPUs + 1 6GB Memory + 50 GB SSD
- 1 machine(client) of 16 vCPUs + 30 GB Memory + 50 GB SSD
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.4)
- Zookeepr r3.4.9
  - Java 8
  - Java(TM) SE Runtime Environment (build 1.8.0_111-b14)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.111-b14, mixed mode)
  - javac 1.8.0_111
- Consul v0.7.2 (Go 1.7.4)
- zetcd v3.1 (Go 1.7.4)
- cetcd v3.1 (Go 1.7.4)

<br><br>
##### Write 2M keys, 1000-client (etcd v3.1 100-conn), 8-byte key, 256-byte value

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-LATENCY-MS.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-THROUGHPUT.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-CPU.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-VMRSS-MB.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/00-write-2M-keys/AVG-TRANSMIT-BYTES-NUM-DELTA">

<br><br>
##### Write 2M keys, 8-byte key, 256-byte value

clients increase from 100 to 2000

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-LATENCY-MS.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-THROUGHPUT.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-CPU.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-VMRSS-MB.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write/AVG-TRANSMIT-BYTES-NUM-DELTA">
