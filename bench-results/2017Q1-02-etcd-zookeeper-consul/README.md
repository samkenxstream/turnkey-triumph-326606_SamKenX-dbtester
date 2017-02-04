

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
+----------------------------+-------------------+------------------------+-----------------------+
|                            | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.3-go1.7.4 |
+----------------------------+-------------------+------------------------+-----------------------+
|      READS-COMPLETED-DELTA |               335 |                    217 |                   257 |
|    SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| WRITES-COMPLETED-DELTA-SUM |         1,231,556 |                955,211 |             2,516,639 |
|  SECTORS-WRITTEN-DELTA-SUM |           729,836 |             10,554,668 |            44,590,452 |
|              AVG-DATA-SIZE |            2.4 GB |                 7.7 GB |                3.0 GB |
|          RECEIVE-BYTES-SUM |            5.1 GB |                 5.4 GB |                7.9 GB |
|         TRANSMIT-BYTES-SUM |            3.9 GB |                 4.4 GB |                6.7 GB |
|              MAX-CPU-USAGE |          445.33 % |               547.73 % |              411.23 % |
|           MAX-MEMORY-USAGE |        1348.55 MB |             3886.00 MB |            6300.68 MB |
|              TOTAL-SECONDS |      348.6809 sec |           325.9369 sec |          833.0403 sec |
|             MAX-THROUGHPUT |    35,607 req/sec |         36,742 req/sec |        14,389 req/sec |
|             AVG-THROUGHPUT |     2,867 req/sec |          3,060 req/sec |         1,200 req/sec |
|             MIN-THROUGHPUT |       154 req/sec |              0 req/sec |             0 req/sec |
|            FASTEST-LATENCY |         1.1331 ms |              1.1178 ms |             3.0202 ms |
|                AVG-LATENCY |        13.9505 ms |             30.7027 ms |           199.2154 ms |
|            SLOWEST-LATENCY |       268.1365 ms |           4791.4295 ms |         22578.4875 ms |
|                Latency p10 |       2.329807 ms |            2.425453 ms |           3.935582 ms |
|                Latency p25 |       5.932790 ms |            3.847339 ms |           7.653584 ms |
|                Latency p50 |      10.211023 ms |            6.907074 ms |          21.364391 ms |
|                Latency p75 |      17.023281 ms |           13.768245 ms |          59.155494 ms |
|                Latency p90 |      28.458718 ms |           23.767844 ms |         161.821492 ms |
|                Latency p95 |      43.808491 ms |           29.821734 ms |         317.361200 ms |
|                Latency p99 |      62.109107 ms |          699.922794 ms |        4416.572182 ms |
|              Latency p99.9 |      92.473754 ms |         2228.556340 ms |       20750.406953 ms |
+----------------------------+-------------------+------------------------+-----------------------+
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



