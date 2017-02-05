# dbtester

[![Build Status](https://img.shields.io/travis/coreos/dbtester.svg?style=flat-square)](https://travis-ci.org/coreos/dbtester) [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/coreos/dbtester)

Distributed database benchmark tester: etcd, Zookeeper, Consul

- Database agent and runner are implemented at https://github.com/coreos/dbtester/tree/master/agent
- Client is implemented at https://github.com/coreos/dbtester/tree/master/control
- System metrics are collected via https://github.com/gyuho/psn
- Data analysis is done via https://github.com/coreos/dbtester/tree/master/analyze
  - https://github.com/gyuho/dataframe
  - https://github.com/gonum/plot

For etcd, we also recommend [etcd benchmark tool](https://github.com/coreos/etcd/tree/master/tools/benchmark).

All logs and results can be found at https://console.cloud.google.com/storage/browser/dbtester-results


<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value value, clients 1 to 1,000

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 30 GB Memory + 150 GB SSD (1 for client)
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.4)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.3 (Go 1.7.4)


```
+-----------------------------+-------------------+------------------------+-----------------------+
|                             | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.3-go1.7.4 |
+-----------------------------+-------------------+------------------------+-----------------------+
|   READS-COMPLETED-DELTA-SUM |                 2 |                    218 |                   126 |
|     SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
|  WRITES-COMPLETED-DELTA-SUM |         1,217,752 |                955,419 |             2,183,202 |
|   SECTORS-WRITTEN-DELTA-SUM |           702,780 |             11,846,404 |             3,991,232 |
|       AVG-DATA-SIZE-ON-DISK |            2.5 GB |                 7.4 GB |                3.1 GB |
|    NETWORK-RECEIVE-DATA-SUM |            5.1 GB |                 5.0 GB |                5.5 GB |
|   NETWORK-TRANSMIT-DATA-SUM |            3.9 GB |                 4.0 GB |                4.2 GB |
|               MAX-CPU-USAGE |          451.00 % |               800.00 % |              409.33 % |
|            MAX-MEMORY-USAGE |        1316.36 MB |             3596.23 MB |            4691.18 MB |
|               TOTAL-SECONDS |      324.5626 sec |           335.0999 sec |          667.9389 sec |
|              MAX-THROUGHPUT |    36,479 req/sec |         41,799 req/sec |        15,969 req/sec |
|              AVG-THROUGHPUT |     3,081 req/sec |          2,961 req/sec |         1,497 req/sec |
|              MIN-THROUGHPUT |        79 req/sec |              0 req/sec |            49 req/sec |
|             FASTEST-LATENCY |         1.0313 ms |              1.1256 ms |             2.9151 ms |
|                 AVG-LATENCY |        13.4807 ms |             25.5780 ms |            47.0926 ms |
|             SLOWEST-LATENCY |       261.8235 ms |           4260.5699 ms |         22260.1112 ms |
|                 Latency p10 |       2.201212 ms |            2.528874 ms |           3.926152 ms |
|                 Latency p25 |       5.636563 ms |            3.821957 ms |           7.640734 ms |
|                 Latency p50 |       9.696159 ms |            6.656013 ms |          19.126381 ms |
|                 Latency p75 |      16.202583 ms |           11.576279 ms |          54.750630 ms |
|                 Latency p90 |      28.434395 ms |           14.472618 ms |          77.993718 ms |
|                 Latency p95 |      44.336815 ms |           16.786180 ms |          91.026490 ms |
|                 Latency p99 |      60.008762 ms |          479.797108 ms |         201.844359 ms |
|               Latency p99.9 |      84.831886 ms |         2725.947720 ms |        1502.535463 ms |
|  CLIENT-NETWORK-RECEIVE-SUM |            270 MB |                 353 MB |                200 MB |
| CLIENT-NETWORK-TRANSMIT-SUM |            1.5 GB |                 1.4 GB |                1.5 GB |
|        CLIENT-MAX-CPU-USAGE |          577.00 % |               496.00 % |              210.00 % |
|     CLIENT-MAX-MEMORY-USAGE |            355 MB |                 3.3 GB |                227 MB |
|          CLIENT-ERROR-COUNT |                 0 |                  7,495 |                     0 |
+-----------------------------+-------------------+------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA">


