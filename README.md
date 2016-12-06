This is an experimental project. Project/code is subject to change anytime.
This is mainly for comparing etcd with other databases. For etcd, we recommend
to just run [etcd benchmark tool](https://github.com/coreos/etcd/tree/master/tools/benchmark).

# dbtester

[![Build Status](https://img.shields.io/travis/coreos/dbtester.svg?style=flat-square)][cistat] [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][dbtester-godoc]

Distributed database tester.

Please let us know or file an issue if:

- Need help with running this testing suite
- Questions about test results

We keep all logs at cloud storage:

- https://console.cloud.google.com/storage/browser/dbtester-results


[cistat]: https://travis-ci.org/coreos/dbtester
[dbtester-godoc]: https://godoc.org/github.com/coreos/dbtester




<br><br><hr>
### Latest Results

THIS IS WORKING IN PROGRESS (more accurate results coming soon...)

- Google Cloud Compute Engine
- 3 machines of 8 vCPUs + 1 6GB Memory + 50 GB SSD
- 1 machine(client) of 16 vCPUs + 30 GB Memory + 50 GB SSD
- Ubuntu 16.10
- Zookeepr v3.4.9
  - Java 8
  - Java(TM) SE Runtime Environment (build 1.8.0_111-b14)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.111-b14, mixed mode)
  - javac 1.8.0_111
- etcd v3.1 (Go 1.7.4)
- Consul v0.7.1 (Go 1.7.3)


<br><br>
##### Write 2M keys, 1000-client (etcd v3.1 100-conn), 8-byte key, 256-byte value

<img src="https://storage.googleapis.com/dbtester-results/2016Q402-etcd-zk-consul/write-2M-keys/avg-latency-ms.svg" alt="2016Q402-etcd-zk-consul/write-2M-keys/avg-latency-ms">

<img src="https://storage.googleapis.com/dbtester-results/2016Q402-etcd-zk-consul/write-2M-keys/throughput.svg" alt="2016Q402-etcd-zk-consul/write-2M-keys/throughput">

<img src="https://storage.googleapis.com/dbtester-results/2016Q402-etcd-zk-consul/write-2M-keys/avg-cpu.svg" alt="2016Q402-etcd-zk-consul/write-2M-keys/avg-cpu">

<img src="https://storage.googleapis.com/dbtester-results/2016Q402-etcd-zk-consul/write-2M-keys/avg-memory-mb.svg" alt="2016Q402-etcd-zk-consul/write-2M-keys/avg-memory-mb">
