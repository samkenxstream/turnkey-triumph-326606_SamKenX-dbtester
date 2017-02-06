

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
+----------------------------------------+-------------------+------------------------+-----------------------+
|                                        | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.3-go1.7.4 |
+----------------------------------------+-------------------+------------------------+-----------------------+
|                          TOTAL-SECONDS |      345.7743 sec |           320.9959 sec |          676.6576 sec |
|                   TOTAL-REQUEST-NUMBER |         1,000,000 |              1,000,000 |             1,000,000 |
|                         MAX-THROUGHPUT |    35,245 req/sec |         37,580 req/sec |        15,488 req/sec |
|                         AVG-THROUGHPUT |     2,892 req/sec |          3,108 req/sec |         1,477 req/sec |
|                         MIN-THROUGHPUT |        23 req/sec |              0 req/sec |             0 req/sec |
|                        FASTEST-LATENCY |         1.0942 ms |              0.9952 ms |             3.0193 ms |
|                            AVG-LATENCY |        14.0207 ms |             34.9720 ms |            44.3717 ms |
|                        SLOWEST-LATENCY |       101.5395 ms |           2334.6007 ms |          3475.2713 ms |
|                            Latency p10 |       2.308513 ms |            2.378813 ms |           3.971904 ms |
|                            Latency p25 |       5.879109 ms |            3.595997 ms |           7.811056 ms |
|                            Latency p50 |      10.100823 ms |            6.320842 ms |          20.224796 ms |
|                            Latency p75 |      17.160120 ms |           12.677824 ms |          56.790033 ms |
|                            Latency p90 |      28.890553 ms |           21.506040 ms |          89.507900 ms |
|                            Latency p95 |      44.861387 ms |           30.713207 ms |         106.807611 ms |
|                            Latency p99 |      62.570152 ms |         1186.951737 ms |         301.241777 ms |
|                          Latency p99.9 |      87.467729 ms |         2241.127019 ms |        2634.872172 ms |
|  SERVER-TOTAL-NETWORK-RECEIVE-DATA-SUM |            5.1 GB |                 5.4 GB |                5.5 GB |
| SERVER-TOTAL-NETWORK-TRANSMIT-DATA-SUM |            3.9 GB |                 4.4 GB |                4.3 GB |
|       CLIENT-TOTAL-NETWORK-RECEIVE-SUM |            271 MB |                 356 MB |                202 MB |
|      CLIENT-TOTAL-NETWORK-TRANSMIT-SUM |            1.5 GB |                 1.4 GB |                1.5 GB |
|                   SERVER-MAX-CPU-USAGE |          449.90 % |               558.97 % |              421.60 % |
|                SERVER-MAX-MEMORY-USAGE |            1.4 GB |                 3.9 GB |                4.8 GB |
|                   CLIENT-MAX-CPU-USAGE |          552.00 % |               617.80 % |              224.00 % |
|                CLIENT-MAX-MEMORY-USAGE |            361 MB |                 4.4 GB |                225 MB |
|                     CLIENT-ERROR-COUNT |                 0 |                  2,093 |                     0 |
|   SERVER-AVG-READS-COMPLETED-DELTA-SUM |                 3 |                    221 |                    11 |
|     SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
|  SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |         1,222,851 |                960,275 |             2,333,199 |
|   SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |           739,936 |             12,769,072 |             4,013,148 |
|           SERVER-AVG-DATA-SIZE-ON-DISK |            2.4 GB |                 8.8 GB |                3.0 GB |
+----------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: connection closed" (count 1,747)
"zk: could not connect to a server" (count 346)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA">



