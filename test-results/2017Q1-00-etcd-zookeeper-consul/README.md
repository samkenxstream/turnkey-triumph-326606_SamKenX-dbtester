

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





<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, Best Throughput (etcd 700, Zookeeper 300, Consul 500 clients)

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
|                          TOTAL-SECONDS |       37.4731 sec |            63.8379 sec |          309.0797 sec |
|                   TOTAL-REQUEST-NUMBER |         1,000,000 |              1,000,000 |             1,000,000 |
|                         MAX-THROUGHPUT |    32,555 req/sec |         38,537 req/sec |        16,919 req/sec |
|                         AVG-THROUGHPUT |    26,685 req/sec |         15,656 req/sec |         3,235 req/sec |
|                         MIN-THROUGHPUT |    15,350 req/sec |              0 req/sec |             0 req/sec |
|                        FASTEST-LATENCY |         4.5030 ms |              2.3895 ms |            14.5025 ms |
|                            AVG-LATENCY |        26.1688 ms |             17.9204 ms |           154.4720 ms |
|                        SLOWEST-LATENCY |       205.8298 ms |           3471.7581 ms |         21463.6662 ms |
|                            Latency p10 |      10.580892 ms |            5.948893 ms |          29.430795 ms |
|                            Latency p25 |      12.626734 ms |            6.817741 ms |          34.096243 ms |
|                            Latency p50 |      18.452075 ms |            7.718642 ms |          42.522161 ms |
|                            Latency p75 |      26.930691 ms |            8.659562 ms |          63.027422 ms |
|                            Latency p90 |      49.843748 ms |            9.750484 ms |         263.017095 ms |
|                            Latency p95 |      67.991988 ms |           10.757740 ms |         601.571541 ms |
|                            Latency p99 |     135.153213 ms |          114.903303 ms |        1237.545701 ms |
|                          Latency p99.9 |     160.393852 ms |         2248.915050 ms |       11496.380241 ms |
|  SERVER-TOTAL-NETWORK-RECEIVE-DATA-SUM |            4.9 GB |                 5.4 GB |                 14 GB |
| SERVER-TOTAL-NETWORK-TRANSMIT-DATA-SUM |            3.8 GB |                 4.3 GB |                 12 GB |
|       CLIENT-TOTAL-NETWORK-RECEIVE-SUM |            274 MB |                 351 MB |                211 MB |
|      CLIENT-TOTAL-NETWORK-TRANSMIT-SUM |            1.4 GB |                 1.4 GB |                1.5 GB |
|                   SERVER-MAX-CPU-USAGE |          416.67 % |               710.17 % |              410.67 % |
|                SERVER-MAX-MEMORY-USAGE |            1.2 GB |                 4.6 GB |                5.9 GB |
|                   CLIENT-MAX-CPU-USAGE |          448.00 % |               365.00 % |              223.00 % |
|                CLIENT-MAX-MEMORY-USAGE |            219 MB |                 1.9 GB |                 90 MB |
|                     CLIENT-ERROR-COUNT |                 0 |                    509 |                     0 |
|   SERVER-AVG-READS-COMPLETED-DELTA-SUM |                 3 |                    253 |                   162 |
|     SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
|  SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           110,436 |                101,558 |               665,193 |
|   SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |           371,988 |              9,306,936 |            37,002,460 |
|           SERVER-AVG-DATA-SIZE-ON-DISK |            2.8 GB |                 7.1 GB |                2.9 GB |
+----------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: could not connect to a server" (count 208)
"zk: connection closed" (count 301)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-CPU.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, 1,000 client

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
|                          TOTAL-SECONDS |       36.5547 sec |            59.5935 sec |          244.8620 sec |
|                   TOTAL-REQUEST-NUMBER |         1,000,000 |              1,000,000 |             1,000,000 |
|                         MAX-THROUGHPUT |    35,025 req/sec |         47,693 req/sec |        15,146 req/sec |
|                         AVG-THROUGHPUT |    27,356 req/sec |         16,697 req/sec |         4,083 req/sec |
|                         MIN-THROUGHPUT |       228 req/sec |              0 req/sec |             0 req/sec |
|                        FASTEST-LATENCY |         4.1857 ms |              1.4036 ms |            27.1684 ms |
|                            AVG-LATENCY |        36.4506 ms |             46.1962 ms |           244.4663 ms |
|                        SLOWEST-LATENCY |       233.1590 ms |           6787.1829 ms |         30479.0204 ms |
|                            Latency p10 |      14.321495 ms |           12.447591 ms |          66.233418 ms |
|                            Latency p25 |      17.784730 ms |           16.822188 ms |          74.271576 ms |
|                            Latency p50 |      23.783217 ms |           20.026128 ms |          92.736772 ms |
|                            Latency p75 |      44.930445 ms |           23.615952 ms |         151.332505 ms |
|                            Latency p90 |      67.957086 ms |           31.726640 ms |         338.858136 ms |
|                            Latency p95 |     121.316986 ms |           78.842487 ms |         696.780894 ms |
|                            Latency p99 |     163.614045 ms |          821.197693 ms |        1887.708027 ms |
|                          Latency p99.9 |     187.122508 ms |         2709.006521 ms |       20942.858889 ms |
|  SERVER-TOTAL-NETWORK-RECEIVE-DATA-SUM |            5.0 GB |                 5.3 GB |                6.6 GB |
| SERVER-TOTAL-NETWORK-TRANSMIT-DATA-SUM |            3.9 GB |                 4.2 GB |                5.4 GB |
|       CLIENT-TOTAL-NETWORK-RECEIVE-SUM |            282 MB |                 357 MB |                243 MB |
|      CLIENT-TOTAL-NETWORK-TRANSMIT-SUM |            1.4 GB |                 1.4 GB |                1.5 GB |
|                   SERVER-MAX-CPU-USAGE |          451.00 % |               723.67 % |              436.67 % |
|                SERVER-MAX-MEMORY-USAGE |            1.2 GB |                 4.7 GB |                4.9 GB |
|                   CLIENT-MAX-CPU-USAGE |          554.00 % |              1227.00 % |              218.00 % |
|                CLIENT-MAX-MEMORY-USAGE |            264 MB |                 5.5 GB |                114 MB |
|                     CLIENT-ERROR-COUNT |                 0 |                  4,962 |                     0 |
|   SERVER-AVG-READS-COMPLETED-DELTA-SUM |                74 |                    315 |                   102 |
|     SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
|  SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           100,644 |                 82,186 |               481,109 |
|   SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |           442,860 |              8,179,180 |            20,434,060 |
|           SERVER-AVG-DATA-SIZE-ON-DISK |            2.7 GB |                 7.3 GB |                3.0 GB |
+----------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: could not connect to a server" (count 3,301)
"zk: connection closed" (count 1,661)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-LATENCY-MS.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-THROUGHPUT.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-CPU.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-VMRSS-MB.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/03-write-1M-keys-1000-client/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 3-million keys, 256-byte key, 1KB value, Best Throughput (etcd 700, Zookeeper 300, Consul 500 clients)

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
|                          TOTAL-SECONDS |      141.1728 sec |          4134.1565 sec |         2979.6229 sec |
|                   TOTAL-REQUEST-NUMBER |         3,000,000 |              3,000,000 |             3,000,000 |
|                         MAX-THROUGHPUT |    32,215 req/sec |         38,259 req/sec |        15,784 req/sec |
|                         AVG-THROUGHPUT |    21,250 req/sec |            465 req/sec |         1,006 req/sec |
|                         MIN-THROUGHPUT |     9,081 req/sec |              0 req/sec |             0 req/sec |
|                        FASTEST-LATENCY |         3.6479 ms |              1.8912 ms |            10.3322 ms |
|                            AVG-LATENCY |        32.9036 ms |             28.5018 ms |           496.5874 ms |
|                        SLOWEST-LATENCY |      1030.3392 ms |           4729.8104 ms |         34484.1252 ms |
|                            Latency p10 |      10.882650 ms |            5.694308 ms |          35.045139 ms |
|                            Latency p25 |      13.005974 ms |            6.357600 ms |          44.711879 ms |
|                            Latency p50 |      18.743999 ms |            7.379354 ms |          66.739638 ms |
|                            Latency p75 |      31.054943 ms |            8.886502 ms |         116.150396 ms |
|                            Latency p90 |      96.027241 ms |           11.474394 ms |         540.305996 ms |
|                            Latency p95 |     123.844111 ms |           14.391699 ms |        1209.603980 ms |
|                            Latency p99 |     148.837786 ms |          722.847010 ms |       11957.052613 ms |
|                          Latency p99.9 |     162.994644 ms |         2457.522611 ms |       26068.610794 ms |
|  SERVER-TOTAL-NETWORK-RECEIVE-DATA-SUM |             15 GB |                  26 GB |                102 GB |
| SERVER-TOTAL-NETWORK-TRANSMIT-DATA-SUM |             12 GB |                  24 GB |                 98 GB |
|       CLIENT-TOTAL-NETWORK-RECEIVE-SUM |            832 MB |                 978 MB |                688 MB |
|      CLIENT-TOTAL-NETWORK-TRANSMIT-SUM |            4.3 GB |                 3.3 GB |                4.4 GB |
|                   SERVER-MAX-CPU-USAGE |          421.00 % |               752.33 % |              410.67 % |
|                SERVER-MAX-MEMORY-USAGE |            2.8 GB |                 7.1 GB |                 17 GB |
|                   CLIENT-MAX-CPU-USAGE |          442.00 % |               376.00 % |              217.00 % |
|                CLIENT-MAX-MEMORY-USAGE |            265 MB |                 1.9 GB |                174 MB |
|                     CLIENT-ERROR-COUNT |                 0 |              1,076,632 |                     0 |
|   SERVER-AVG-READS-COMPLETED-DELTA-SUM |                14 |                    309 |                94,864 |
|     SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
|  SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           333,196 |                360,314 |             3,689,221 |
|   SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |         1,604,172 |             67,295,610 |           794,919,246 |
|           SERVER-AVG-DATA-SIZE-ON-DISK |            6.5 GB |                  27 GB |                8.2 GB |
+----------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: connection closed" (count 6,678)
"zk: could not connect to a server" (count 1,069,954)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-LATENCY-MS.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-THROUGHPUT.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-CPU.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-VMRSS-MB.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys/AVG-TRANSMIT-BYTES-NUM-DELTA">



