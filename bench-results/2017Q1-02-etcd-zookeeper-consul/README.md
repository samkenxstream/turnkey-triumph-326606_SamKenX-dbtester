

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
|      READS-COMPLETED-DELTA |                 3 |                    214 |                    67 |
|    SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| WRITES-COMPLETED-DELTA-SUM |           1224742 |                 957554 |               2253556 |
|  SECTORS-WRITTEN-DELTA-SUM |            730236 |               11224704 |              21641916 |
|              AVG-DATA-SIZE |            2.4 GB |                 7.4 GB |                3.1 GB |
|          RECEIVE-BYTES-SUM |            5.1 GB |                 6.5 GB |                9.5 GB |
|         TRANSMIT-BYTES-SUM |            3.9 GB |                 5.5 GB |                8.2 GB |
|              MAX-CPU-USAGE |          440.43 % |               567.33 % |              505.67 % |
|           MAX-MEMORY-USAGE |        1314.30 MB |             4165.09 MB |            6207.47 MB |
|              TOTAL-SECONDS |      344.1735 sec |           319.8246 sec |          719.7771 sec |
|             MIN-THROUGHPUT |       192 req/sec |              0 req/sec |             0 req/sec |
|             AVG-THROUGHPUT |     2,905 req/sec |          3,108 req/sec |         1,389 req/sec |
|             MAX-THROUGHPUT |    35,777 req/sec |         41,982 req/sec |        16,543 req/sec |
|            SLOWEST-LATENCY |       109.1354 ms |           3802.0451 ms |         24873.9576 ms |
|                AVG-LATENCY |        13.8141 ms |             26.6582 ms |           108.3479 ms |
|            FASTEST-LATENCY |         1.1083 ms |              1.0881 ms |             2.9683 ms |
|                Latency p10 |       2.312361 ms |            2.398301 ms |           3.915505 ms |
|                Latency p25 |       5.875704 ms |            3.627383 ms |           7.812279 ms |
|                Latency p50 |       9.941706 ms |            6.432763 ms |          18.977204 ms |
|                Latency p75 |      16.972703 ms |           11.814841 ms |          52.003318 ms |
|                Latency p90 |      28.773894 ms |           16.627078 ms |          84.974572 ms |
|                Latency p95 |      44.682598 ms |           21.045116 ms |         183.677871 ms |
|                Latency p99 |      58.665081 ms |          572.640038 ms |        1190.464210 ms |
|              Latency p99.9 |      90.857600 ms |         2619.729132 ms |       16556.840511 ms |
+----------------------------+-------------------+------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-02-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA">



