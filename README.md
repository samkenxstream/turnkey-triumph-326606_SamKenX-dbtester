# dbtester

[![Build Status](https://img.shields.io/travis/coreos/dbtester.svg?style=flat-square)](https://travis-ci.org/coreos/dbtester) [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/coreos/dbtester)

Distributed database tester.

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

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/01-write-2M-keys/avg-latency-ms.svg" alt="2016Q4-01-etcd-zk-consul/01-write-2M-keys/avg-latency-ms">

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/01-write-2M-keys/avg-throughput.svg" alt="2016Q4-01-etcd-zk-consul/01-write-2M-keys/avg-throughput">

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/01-write-2M-keys/avg-cpu.svg" alt="2016Q4-01-etcd-zk-consul/01-write-2M-keys/avg-cpu">

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/01-write-2M-keys/avg-memory-mb.svg" alt="2016Q4-01-etcd-zk-consul/01-write-2M-keys/avg-memory-mb">


<br><br>
##### Write 2M keys (1000QPS), 1000-client (etcd v3.1 100-conn), 8-byte key, 256-byte value

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/02-write-2M-keys-1000QPS-no-zetcd-cetcd/avg-latency-ms.svg" alt="2016Q4-01-etcd-zk-consul/02-write-2M-keys-1000QPS-no-zetcd-cetcd/avg-latency-ms">

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/02-write-2M-keys-1000QPS-no-zetcd-cetcd/avg-throughput.svg" alt="2016Q4-01-etcd-zk-consul/02-write-2M-keys-1000QPS-no-zetcd-cetcd/avg-throughput">

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/02-write-2M-keys-1000QPS-no-zetcd-cetcd/avg-cpu.svg" alt="2016Q4-01-etcd-zk-consul/02-write-2M-keys-1000QPS-no-zetcd-cetcd/avg-cpu">

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/02-write-2M-keys-1000QPS-no-zetcd-cetcd/avg-memory-mb.svg" alt="2016Q4-01-etcd-zk-consul/02-write-2M-keys-1000QPS-no-zetcd-cetcd/avg-memory-mb">


<br><br>
##### Write 500K keys, 1-client, 8-byte key, 256-byte value

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/03-write-500K-keys-1CONN-no-zetcd-cetcd/avg-latency-ms.svg" alt="2016Q4-01-etcd-zk-consul/03-write-500K-keys-1CONN-no-zetcd-cetcd/avg-latency-ms">

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/03-write-500K-keys-1CONN-no-zetcd-cetcd/avg-throughput.svg" alt="2016Q4-01-etcd-zk-consul/03-write-500K-keys-1CONN-no-zetcd-cetcd/avg-throughput">

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/03-write-500K-keys-1CONN-no-zetcd-cetcd/avg-cpu.svg" alt="2016Q4-01-etcd-zk-consul/03-write-500K-keys-1CONN-no-zetcd-cetcd/avg-cpu">

<img src="https://storage.googleapis.com/dbtester-results/2016Q4-01-etcd-zk-consul/03-write-500K-keys-1CONN-no-zetcd-cetcd/avg-memory-mb.svg" alt="2016Q4-01-etcd-zk-consul/03-write-500K-keys-1CONN-no-zetcd-cetcd/avg-memory-mb">

