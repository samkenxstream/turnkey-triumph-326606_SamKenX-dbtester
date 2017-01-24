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

All logs and results can be found at https://console.cloud.google.com/storage/browser/dbtester-results

<br><br>
##### Write 1M keys, 1000-client, 256-byte key, 1KB value

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 30 GB Memory + 150 GB SSD (1 for client)
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.4)
- Zookeeper r3.4.9
  - Java 8
  - Java(TM) SE Runtime Environment (build 1.8.0_111-b14)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.111-b14, mixed mode)
  - javac 1.8.0_111
- Consul v0.7.2 (Go 1.7.4)


```
+----------------------------+--------------------+------------------------+-----------------------+
|                            | etcd-v3.1-go1.7.4  | zookeeper-r3.4.9-java8 | consul-v0.7.2-go1.7.4 |
+----------------------------+--------------------+------------------------+-----------------------+
|   READS-COMPLETED-DELTA    |         6          |           6            |          15           |
|  SECTORS-READS-DELTA-SUM   |         0          |           0            |           0           |
| WRITES-COMPLETED-DELTA-SUM |       96474        |         77628          |        940695         |
| SECTORS-WRITTEN-DELTA-SUM  |       542512       |        9387436         |       41272068        |
|     RECEIVE-BYTES-SUM      |       4.9 GB       |         5.1 GB         |        7.7 GB         |
|     TRANSMIT-BYTES-SUM     |       3.7 GB       |         4.1 GB         |        6.5 GB         |
|       TOTAL-SECONDS        |    36.2024 sec     |      62.0373 sec       |     467.9311 sec      |
|       AVG-THROUGHPUT       | 27622.4453 req/sec |   15951.5555 req/sec   |   2137.0667 req/sec   |
|      SLOWEST-LATENCY       |    246.4560 ms     |      6650.0930 ms      |     30388.9318 ms     |
|      FASTEST-LATENCY       |     5.3413 ms      |       1.7698 ms        |      21.5605 ms       |
|        AVG-LATENCY         |     36.1057 ms     |       37.7865 ms       |      467.4253 ms      |
|        Latency p10         |    13.712090 ms    |      11.923543 ms      |     65.910086 ms      |
|        Latency p25         |    16.625779 ms    |      14.581663 ms      |     77.221971 ms      |
|        Latency p50         |    22.306160 ms    |      19.217649 ms      |     120.663354 ms     |
|        Latency p75         |    40.376905 ms    |      23.642903 ms      |     716.373543 ms     |
|        Latency p90         |    65.849751 ms    |      28.756700 ms      |    1068.038406 ms     |
|        Latency p95         |   137.545464 ms    |      59.868096 ms      |    1080.751412 ms     |
|        Latency p99         |   177.127309 ms    |     544.858078 ms      |    2686.919571 ms     |
|       Latency p99.9        |   198.540415 ms    |     2457.827147 ms     |    19041.188919 ms    |
+----------------------------+--------------------+------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-TRANSMIT-BYTES-NUM-DELTA">
