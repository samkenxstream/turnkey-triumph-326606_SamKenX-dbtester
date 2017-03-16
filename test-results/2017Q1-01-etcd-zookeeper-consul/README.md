

<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value value, clients 1 to 1,000

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10 (GNU/Linux kernel 4.8.0-42-generic)
- `ulimit -n` is 120000
- etcd tip (Go 1.8.0, git SHA d78b03fb27374c370d82973a94dda9f59647e750)
- Zookeeper r3.5.2-alpha
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v0.7.5 (Go 1.8.0)


```
+---------------------------------------+------------------+------------------------------+-----------------------+
|                                       | etcd-tip-go1.8.0 | zookeeper-r3.5.2-alpha-java8 | consul-v0.7.5-go1.8.0 |
+---------------------------------------+------------------+------------------------------+-----------------------+
|                         TOTAL-SECONDS |     377.2927 sec |                 327.6324 sec |          698.7478 sec |
|                  TOTAL-REQUEST-NUMBER |        1,000,000 |                    1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |   37,106 req/sec |               25,009 req/sec |        16,003 req/sec |
|                        AVG-THROUGHPUT |    2,650 req/sec |                3,048 req/sec |         1,431 req/sec |
|                        MIN-THROUGHPUT |      201 req/sec |                    0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |        1.1514 ms |                    1.0570 ms |             3.1691 ms |
|                           AVG-LATENCY |       13.9042 ms |                   22.0313 ms |            47.5751 ms |
|                       SLOWEST-LATENCY |      112.1736 ms |                 2662.2537 ms |          3416.4836 ms |
|                           Latency p10 |      2.590518 ms |                  2.554608 ms |           4.120475 ms |
|                           Latency p25 |      6.222186 ms |                  3.882311 ms |           7.906411 ms |
|                           Latency p50 |     10.331986 ms |                  6.389982 ms |          19.743291 ms |
|                           Latency p75 |     16.628607 ms |                 22.490612 ms |          61.329955 ms |
|                           Latency p90 |     28.355336 ms |                 34.082148 ms |          89.313236 ms |
|                           Latency p95 |     43.345061 ms |                 51.386743 ms |         122.111399 ms |
|                           Latency p99 |     60.595924 ms |                224.327398 ms |         309.823246 ms |
|                         Latency p99.9 |     95.336990 ms |               1549.797635 ms |        2714.083344 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |           5.1 GB |                       5.4 GB |                5.6 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |           3.9 GB |                       4.3 GB |                4.3 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |           270 MB |                       356 MB |                202 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |           1.5 GB |                       1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |         425.80 % |                     257.00 % |              424.53 % |
|               SERVER-MAX-MEMORY-USAGE |           1.4 GB |                        16 GB |                4.9 GB |
|                  CLIENT-MAX-CPU-USAGE |         462.00 % |                     354.50 % |              425.70 % |
|               CLIENT-MAX-MEMORY-USAGE |           308 MB |                       4.8 GB |                201 MB |
|                    CLIENT-ERROR-COUNT |                0 |                        1,194 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |               30 |                          207 |                    31 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                0 |                            0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |        1,525,546 |                    1,234,154 |             3,352,596 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |       32,521,080 |                   35,044,520 |           106,551,704 |
|           SERVER-AVG-DISK-SPACE-USAGE |           2.4 GB |                          0 B |                2.9 GB |
+---------------------------------------+------------------+------------------------------+-----------------------+


zookeeper__r3_5_2_alpha errors:
"zk: could not connect to a server" (count 195)
"zk: connection closed" (count 999)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/MAX-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READ-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, Best Throughput (etcd 1,000, Zookeeper 500, Consul 500 clients)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10 (GNU/Linux kernel 4.8.0-42-generic)
- `ulimit -n` is 120000
- etcd tip (Go 1.8.0, git SHA d78b03fb27374c370d82973a94dda9f59647e750)
- Zookeeper r3.5.2-alpha
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v0.7.5 (Go 1.8.0)


```
+---------------------------------------+------------------+------------------------------+-----------------------+
|                                       | etcd-tip-go1.8.0 | zookeeper-r3.5.2-alpha-java8 | consul-v0.7.5-go1.8.0 |
+---------------------------------------+------------------+------------------------------+-----------------------+
|                         TOTAL-SECONDS |      27.9354 sec |                  51.7589 sec |          482.4360 sec |
|                  TOTAL-REQUEST-NUMBER |        1,000,000 |                    1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |   38,791 req/sec |               24,142 req/sec |        15,858 req/sec |
|                        AVG-THROUGHPUT |   35,796 req/sec |               19,319 req/sec |         2,072 req/sec |
|                        MIN-THROUGHPUT |   12,028 req/sec |                   47 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |        4.0583 ms |                    3.6546 ms |            13.9067 ms |
|                           AVG-LATENCY |       27.7560 ms |                   25.2746 ms |           241.1081 ms |
|                       SLOWEST-LATENCY |     1084.3727 ms |                 1690.0324 ms |          2829.0558 ms |
|                           Latency p10 |     12.653617 ms |                  7.775213 ms |          32.308527 ms |
|                           Latency p25 |     15.548106 ms |                  8.989187 ms |          45.157053 ms |
|                           Latency p50 |     21.642950 ms |                 11.920882 ms |         306.531596 ms |
|                           Latency p75 |     34.058936 ms |                 18.416282 ms |         383.361003 ms |
|                           Latency p90 |     54.289308 ms |                 30.998200 ms |         399.483213 ms |
|                           Latency p95 |     59.761141 ms |                 42.943748 ms |         407.839444 ms |
|                           Latency p99 |     76.450636 ms |                340.584561 ms |         426.220478 ms |
|                         Latency p99.9 |     94.377491 ms |               1198.611571 ms |        1792.641180 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |           5.0 GB |                       4.9 GB |                5.7 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |           3.8 GB |                       3.9 GB |                4.4 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |           279 MB |                       353 MB |                228 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |           1.4 GB |                       1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |         448.33 % |                     580.70 % |              406.00 % |
|               SERVER-MAX-MEMORY-USAGE |           1.1 GB |                        16 GB |                4.7 GB |
|                  CLIENT-MAX-CPU-USAGE |         464.40 % |                     220.00 % |              255.00 % |
|               CLIENT-MAX-MEMORY-USAGE |           244 MB |                       2.5 GB |                 90 MB |
|                    CLIENT-ERROR-COUNT |                0 |                           20 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |               11 |                          212 |                    16 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                0 |                            0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |          100,970 |                      151,117 |             3,254,568 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |       20,736,456 |                   16,501,248 |           285,723,776 |
|           SERVER-AVG-DISK-SPACE-USAGE |           2.6 GB |                          0 B |                2.8 GB |
+---------------------------------------+------------------+------------------------------+-----------------------+


zookeeper__r3_5_2_alpha errors:
"zk: could not connect to a server" (count 20)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/MAX-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 3-million keys, 256-byte key, 1KB value, Best Throughput (etcd 1,000, Zookeeper 500, Consul 500 clients)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10 (GNU/Linux kernel 4.8.0-42-generic)
- `ulimit -n` is 120000
- etcd tip (Go 1.8.0, git SHA d78b03fb27374c370d82973a94dda9f59647e750)
- Zookeeper r3.5.2-alpha
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v0.7.5 (Go 1.8.0)


```
+---------------------------------------+------------------+------------------------------+-----------------------+
|                                       | etcd-tip-go1.8.0 | zookeeper-r3.5.2-alpha-java8 | consul-v0.7.5-go1.8.0 |
+---------------------------------------+------------------+------------------------------+-----------------------+
|                         TOTAL-SECONDS |      84.2543 sec |                 305.1659 sec |         1206.0352 sec |
|                  TOTAL-REQUEST-NUMBER |        3,000,000 |                    3,000,000 |             3,000,000 |
|                        MAX-THROUGHPUT |   39,973 req/sec |               23,847 req/sec |        15,016 req/sec |
|                        AVG-THROUGHPUT |   35,606 req/sec |                9,723 req/sec |         2,487 req/sec |
|                        MIN-THROUGHPUT |    8,413 req/sec |                    0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |        3.6522 ms |                    1.7530 ms |             9.3109 ms |
|                           AVG-LATENCY |       28.0355 ms |                   35.8848 ms |           200.9806 ms |
|                       SLOWEST-LATENCY |     1045.7943 ms |                 4200.7088 ms |         24326.9072 ms |
|                           Latency p10 |     12.163786 ms |                  7.437302 ms |          33.257408 ms |
|                           Latency p25 |     14.799691 ms |                  8.277894 ms |          37.723426 ms |
|                           Latency p50 |     21.589291 ms |                 10.689390 ms |          44.522580 ms |
|                           Latency p75 |     34.128977 ms |                 15.809325 ms |          57.370064 ms |
|                           Latency p90 |     55.153568 ms |                 30.973644 ms |         153.991054 ms |
|                           Latency p95 |     62.265352 ms |                 68.815970 ms |         459.915452 ms |
|                           Latency p99 |     95.655432 ms |                656.062715 ms |        4503.134832 ms |
|                         Latency p99.9 |    151.005460 ms |               2498.386467 ms |       12409.886568 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |            15 GB |                        14 GB |                100 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |            11 GB |                        11 GB |                 97 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |           834 MB |                       1.1 GB |                647 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |           4.3 GB |                       4.2 GB |                4.4 GB |
|                  SERVER-MAX-CPU-USAGE |         489.00 % |                     495.67 % |              462.00 % |
|               SERVER-MAX-MEMORY-USAGE |           2.8 GB |                        26 GB |                 19 GB |
|                  CLIENT-MAX-CPU-USAGE |         490.00 % |                     252.00 % |              225.00 % |
|               CLIENT-MAX-MEMORY-USAGE |           301 MB |                       3.3 GB |                156 MB |
|                    CLIENT-ERROR-COUNT |                0 |                       32,647 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |              147 |                          792 |                 2,141 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                0 |                            0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |          301,127 |                      672,380 |             2,944,882 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |       61,959,344 |                  196,435,376 |           958,732,432 |
|           SERVER-AVG-DISK-SPACE-USAGE |           6.6 GB |                          0 B |                8.5 GB |
+---------------------------------------+------------------+------------------------------+-----------------------+


zookeeper__r3_5_2_alpha errors:
"zk: could not connect to a server" (count 27,419)
"zk: connection closed" (count 5,228)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/MAX-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-READ-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-too-many-keys/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, 100 clients, 1000 QPS Limit

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10 (GNU/Linux kernel 4.8.0-42-generic)
- `ulimit -n` is 120000
- etcd tip (Go 1.8.0, git SHA d78b03fb27374c370d82973a94dda9f59647e750)
- Zookeeper r3.5.2-alpha
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v0.7.5 (Go 1.8.0)


```
+---------------------------------------+------------------+------------------------------+-----------------------+
|                                       | etcd-tip-go1.8.0 | zookeeper-r3.5.2-alpha-java8 | consul-v0.7.5-go1.8.0 |
+---------------------------------------+------------------+------------------------------+-----------------------+
|                         TOTAL-SECONDS |     999.0088 sec |                1002.8312 sec |         1010.3055 sec |
|                  TOTAL-REQUEST-NUMBER |        1,000,000 |                    1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |    1,181 req/sec |                1,975 req/sec |         2,064 req/sec |
|                        AVG-THROUGHPUT |    1,000 req/sec |                  997 req/sec |           989 req/sec |
|                        MIN-THROUGHPUT |      819 req/sec |                    0 req/sec |           100 req/sec |
|                       FASTEST-LATENCY |        1.1741 ms |                    1.0934 ms |             3.7535 ms |
|                           AVG-LATENCY |        4.6191 ms |                    3.2957 ms |            63.9923 ms |
|                       SLOWEST-LATENCY |      213.1974 ms |                 2801.7577 ms |          1837.6108 ms |
|                           Latency p10 |      2.690404 ms |                  1.722341 ms |           8.356415 ms |
|                           Latency p25 |      3.402413 ms |                  1.947233 ms |          19.390184 ms |
|                           Latency p50 |      4.436716 ms |                  2.237011 ms |          71.852670 ms |
|                           Latency p75 |      5.603319 ms |                  2.531433 ms |          94.524931 ms |
|                           Latency p90 |      6.684710 ms |                  2.764637 ms |         105.113648 ms |
|                           Latency p95 |      7.463689 ms |                  2.944441 ms |         112.451436 ms |
|                           Latency p99 |      9.220852 ms |                  4.915394 ms |         143.747948 ms |
|                         Latency p99.9 |     15.102828 ms |                299.410120 ms |         864.272976 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |           5.7 GB |                       5.8 GB |                5.9 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |           4.5 GB |                       4.6 GB |                4.7 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |           264 MB |                       354 MB |                221 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |           1.5 GB |                       1.5 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |          92.17 % |                     191.33 % |              278.60 % |
|               SERVER-MAX-MEMORY-USAGE |           1.6 GB |                        17 GB |                4.3 GB |
|                  CLIENT-MAX-CPU-USAGE |          51.00 % |                      59.40 % |               53.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            92 MB |                       662 MB |                 73 MB |
|                    CLIENT-ERROR-COUNT |                0 |                            4 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |              183 |                          312 |                   339 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                0 |                            0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |        5,703,724 |                    6,488,292 |            10,385,551 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |       64,969,512 |                  123,190,708 |           266,572,140 |
|           SERVER-AVG-DISK-SPACE-USAGE |           2.5 GB |                          0 B |                2.8 GB |
+---------------------------------------+------------------+------------------------------+-----------------------+


zookeeper__r3_5_2_alpha errors:
"zk: connection closed" (count 4)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/MAX-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-READ-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys-1000QPS/AVG-TRANSMIT-BYTES-NUM-DELTA">



