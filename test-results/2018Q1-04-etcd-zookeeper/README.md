

<br><br><hr>
##### Write 100K keys, 256-byte key, 1KB value, 1 client

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 17.10 (GNU/Linux kernel 4.13.0-25-generic)
- `ulimit -n` is 120000
- etcd v3.3.0 (Go 1.9.3)
- Zookeeper r3.5.3-beta
  - Java 8
  - javac 1.8.0_151
  - Java(TM) SE Runtime Environment (build 1.8.0_151-b12)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.151-b12, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`


```
+---------------------------------------+---------------------+-----------------------------+
|                                       | etcd-v3.3.0-go1.9.3 | zookeeper-r3.5.3-beta-java8 |
+---------------------------------------+---------------------+-----------------------------+
|                         TOTAL-SECONDS |        175.1067 sec |                193.2046 sec |
|                  TOTAL-REQUEST-NUMBER |             100,000 |                     100,000 |
|                        MAX-THROUGHPUT |         632 req/sec |                 594 req/sec |
|                        AVG-THROUGHPUT |         571 req/sec |                 517 req/sec |
|                        MIN-THROUGHPUT |         166 req/sec |                 213 req/sec |
|                       FASTEST-LATENCY |           0.9618 ms |                   1.1486 ms |
|                           AVG-LATENCY |           1.7462 ms |                   1.9308 ms |
|                       SLOWEST-LATENCY |          17.2672 ms |                 110.7389 ms |
|                           Latency p10 |         1.218092 ms |                 1.402825 ms |
|                           Latency p25 |         1.305288 ms |                 1.478765 ms |
|                           Latency p50 |         1.513458 ms |                 1.600954 ms |
|                           Latency p75 |         2.123979 ms |                 2.601295 ms |
|                           Latency p90 |         2.366836 ms |                 2.793483 ms |
|                           Latency p95 |         2.499598 ms |                 2.890841 ms |
|                           Latency p99 |         3.991311 ms |                 3.223381 ms |
|                         Latency p99.9 |         5.451962 ms |                 3.994881 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |              557 MB |                      617 MB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |              427 MB |                      511 MB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |               31 MB |                       35 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |              161 MB |                      140 MB |
|                  SERVER-MAX-CPU-USAGE |             53.20 % |                     91.33 % |
|               SERVER-MAX-MEMORY-USAGE |              456 MB |                      3.4 GB |
|                  CLIENT-MAX-CPU-USAGE |             24.00 % |                     10.00 % |
|               CLIENT-MAX-MEMORY-USAGE |               19 MB |                       22 MB |
|                    CLIENT-ERROR-COUNT |                   0 |                           0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                  18 |                          70 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                   0 |                           0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |             729,968 |                     699,827 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |           7,266,272 |                   7,540,328 |
|           SERVER-AVG-DISK-SPACE-USAGE |              399 MB |                      317 MB |
+---------------------------------------+---------------------+-----------------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-LATENCY-MS.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-LATENCY-MS-BY-KEY.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-THROUGHPUT.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-CPU.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/MAX-CPU.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-VMRSS-MB.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-VMRSS-MB-BY-KEY.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-READS-COMPLETED-DELTA.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-SECTORS-READ-DELTA.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-WRITES-COMPLETED-DELTA.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-READ-BYTES-NUM-DELTA.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2018Q1-04-etcd-zookeeper/write-100K-keys-1-client/AVG-TRANSMIT-BYTES-NUM-DELTA">



