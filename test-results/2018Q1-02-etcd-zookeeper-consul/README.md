

<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, Best Throughput (etcd 1K clients with 100 conns, Zookeeper 700, Consul 500 clients)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 17.10 (GNU/Linux kernel 4.13.0-25-generic)
- `ulimit -n` is 120000
- etcd v3.3.0 (Go 1.9.2)
- Zookeeper r3.5.3-beta
  - Java 8
  - javac 1.8.0_151
  - Java(TM) SE Runtime Environment (build 1.8.0_151-b12)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.151-b12, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v1.0.2 (Go 1.9.2)


```
+---------------------------------------+---------------------+-----------------------------+-----------------------+
|                                       | etcd-v3.3.0-go1.9.2 | zookeeper-r3.5.3-beta-java8 | consul-v1.0.2-go1.9.2 |
+---------------------------------------+---------------------+-----------------------------+-----------------------+
|                         TOTAL-SECONDS |         28.3623 sec |                 59.2167 sec |          178.9443 sec |
|                  TOTAL-REQUEST-NUMBER |           1,000,000 |                   1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |      37,330 req/sec |              25,124 req/sec |        15,865 req/sec |
|                        AVG-THROUGHPUT |      35,258 req/sec |              16,842 req/sec |         5,588 req/sec |
|                        MIN-THROUGHPUT |      13,505 req/sec |                  20 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |           4.6073 ms |                   2.9094 ms |            11.6604 ms |
|                           AVG-LATENCY |          28.2625 ms |                  30.9499 ms |            89.4351 ms |
|                       SLOWEST-LATENCY |         117.4918 ms |                4564.6788 ms |          4616.2947 ms |
|                           Latency p10 |        13.508626 ms |                 9.068163 ms |          30.408863 ms |
|                           Latency p25 |        16.869586 ms |                 9.351597 ms |          34.224021 ms |
|                           Latency p50 |        22.167478 ms |                10.093377 ms |          39.881181 ms |
|                           Latency p75 |        34.855941 ms |                14.951189 ms |          52.644787 ms |
|                           Latency p90 |        54.613394 ms |                28.497256 ms |         118.340402 ms |
|                           Latency p95 |        59.785127 ms |                72.671788 ms |         229.129526 ms |
|                           Latency p99 |        74.139638 ms |               273.218523 ms |        1495.660763 ms |
|                         Latency p99.9 |        97.385495 ms |              2526.873285 ms |        3499.225138 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |              5.1 GB |                      4.6 GB |                5.6 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |              3.8 GB |                      3.6 GB |                4.4 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |              252 MB |                      357 MB |                206 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |              1.5 GB |                      1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |            446.83 % |                   1122.00 % |              426.33 % |
|               SERVER-MAX-MEMORY-USAGE |              1.1 GB |                       15 GB |                4.6 GB |
|                  CLIENT-MAX-CPU-USAGE |            606.00 % |                    314.00 % |              215.00 % |
|               CLIENT-MAX-MEMORY-USAGE |               96 MB |                      2.4 GB |                 86 MB |
|                    CLIENT-ERROR-COUNT |                   0 |                       2,652 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                   0 |                         237 |                     2 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                   0 |                           0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |             108,067 |                     157,034 |               675,072 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |          20,449,360 |                  16,480,488 |           106,836,768 |
|           SERVER-AVG-DISK-SPACE-USAGE |              2.6 GB |                      6.9 GB |                2.9 GB |
+---------------------------------------+---------------------+-----------------------------+-----------------------+


zookeeper__r3_5_3_beta errors:
"zk: connection closed" (count 2,264)
"zk: could not connect to a server" (count 388)
```


<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-CPU.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/MAX-CPU.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2018Q1-02-etcd-zookeeper-consul/write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">



