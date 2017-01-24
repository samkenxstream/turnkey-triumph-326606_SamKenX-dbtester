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
- Zookeeper r3.4.9
  - Java 8
  - Java(TM) SE Runtime Environment (build 1.8.0_111-b14)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.111-b14, mixed mode)
  - javac 1.8.0_111
- Consul v0.7.2 (Go 1.7.4)
- zetcd v3.1 (Go 1.7.4)
- cetcd v3.1 (Go 1.7.4)

<br><br>


Below is latency distribution.

| Write 2M | etcd | Zookeeper | Consul |
|:-:|:-:|:-:|:-:|:-:|
| Total | 49.9517 sec | 57.648 sec | 196.5391 sec |
| Slowest latency | 219.4978 ms | 3673.511 ms | 3456.8632 ms |
| Fastest latency | 2.4053 ms | 1.3334 ms | 11.6683 ms |
| Average latency | 24.9106 ms | 25.6085 ms | 98.0761 ms |
| 10th percentile | 15.252641 ms | 8.694673 ms | 58.773931 ms |
| 90th percentile | 34.603153 ms | 21.260465 ms | 155.155478 ms |
| 95th percentile | 58.790464 ms | 38.915503 ms | 228.378322 ms |
| 99th percentile | 109.932998 ms | 249.609641 ms | 397.145186 ms |
| 99.9th percentile | 163.174532 ms | 1791.313948 ms | 2063.012564 ms |

| Write 2M, 1000QPS | etcd | Zookeeper | Consul |
|:-:|:-:|:-:|:-:|:-:|
| Total | 1999.0071 sec | 2001.5838 sec | 2136.9951 sec |
| Slowest latency | 266.4827 ms | 2391.7506 ms | 16507.9402 ms |
| Fastest latency | 1.1254 ms | 0.8987 ms | 3.2398 ms |
| Average latency | 2.9778 ms | 6.6063 ms | 209.3409 ms |
| 10th percentile | 1.673424 ms | 1.168258 ms | 8.883748 ms |
| 90th percentile | 4.249615 ms | 2.169022 ms | 192.123086 ms |
| 95th percentile | 4.869974 ms | 2.547155 ms | 203.922319 ms |
| 99th percentile | 6.493646 ms | 39.196783 ms | 4457.758814 ms |
| 99.9th percentile | 26.053565 ms | 1205.014692 ms | 12980.91952 ms |

##### Write 2M keys, 1000-client (etcd 100 TCP conns), 8-byte key, 256-byte value

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
##### Write 2M keys 1000QPS, 1000-client (etcd 100 TCP conns), 8-byte key, 256-byte value

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-LATENCY-MS.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-THROUGHPUT.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-CPU.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-VMRSS-MB.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/02-write-2M-keys-rate-limited/AVG-TRANSMIT-BYTES-NUM-DELTA">


<br><br>
##### Write 2M keys, 8-byte key, 256-byte value

clients increase from 1 to 1000

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
