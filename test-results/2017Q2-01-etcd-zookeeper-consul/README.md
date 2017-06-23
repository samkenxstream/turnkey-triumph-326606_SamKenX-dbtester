
<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/2017Q2-01-write-1M-cpu-client-scaling.png" alt="2017Q2-01-write-1M-cpu-client-scaling">

<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/2017Q2-01-write-1M-latency-best-throughput.png" alt="2017Q2-01-write-1M-latency-best-throughput">

<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/2017Q2-01-write-1M-latency-by-key-best-throughput.png" alt="2017Q2-01-write-1M-latency-by-key-best-throughput">

<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/2017Q2-01-write-1M-memory-by-key-best-throughput.png" alt="2017Q2-01-write-1M-memory-by-key-best-throughput">

<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/2017Q2-01-write-1M-network-traffic-best-throughput.png" alt="2017Q2-01-write-1M-network-traffic-throughput">

<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/2017Q2-01-write-1M-sector-writes-client-scaling.png" alt="2017Q2-01-write-1M-sector-writes-client-scaling">

<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/2017Q2-01-write-1M-throughput-client-scaling.png" alt="2017Q2-01-write-1M-throughput-client-scaling">

<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/2017Q2-01-read-3M-latency-1000-clients.png" alt="2017Q2-01-read-3M-latency-1000-clients">


<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value value, clients 1 to 1,000

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10 (GNU/Linux kernel 4.8.0-49-generic)
- `ulimit -n` is 120000
- etcd tip (Go 1.8.1, git SHA f4641accc34be80c255ff87673a1cb92342abedd)
- Zookeeper r3.5.3-beta
  - Java 8
  - javac 1.8.0_131
  - Java(TM) SE Runtime Environment (build 1.8.0_131-b11)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.131-b11, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v0.8.3 (Go 1.8.1)


```
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                                       | etcd-tip-go1.8.1 | zookeeper-r3.5.3-beta-java8 | consul-v0.8.3-go1.8.1 |
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                         TOTAL-SECONDS |     351.6020 sec |                312.1674 sec |          693.4096 sec |
|                  TOTAL-REQUEST-NUMBER |        1,000,000 |                   1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |   38,664 req/sec |              26,396 req/sec |        16,112 req/sec |
|                        AVG-THROUGHPUT |    2,844 req/sec |               3,203 req/sec |         1,442 req/sec |
|                        MIN-THROUGHPUT |        0 req/sec |                 322 req/sec |             7 req/sec |
|                       FASTEST-LATENCY |        1.0256 ms |                   1.0387 ms |             3.1768 ms |
|                           AVG-LATENCY |       13.5717 ms |                  19.3475 ms |            43.1156 ms |
|                       SLOWEST-LATENCY |     1691.4180 ms |                2227.7036 ms |           606.1826 ms |
|                           Latency p10 |      2.402266 ms |                 2.591794 ms |           4.077258 ms |
|                           Latency p25 |      6.143995 ms |                 4.015241 ms |           8.037402 ms |
|                           Latency p50 |     10.108521 ms |                 6.376892 ms |          21.082171 ms |
|                           Latency p75 |     16.263905 ms |                14.289960 ms |          57.802632 ms |
|                           Latency p90 |     26.587891 ms |                35.872999 ms |          92.286646 ms |
|                           Latency p95 |     43.256411 ms |                54.634488 ms |         146.055426 ms |
|                           Latency p99 |     59.931281 ms |               207.573536 ms |         345.999896 ms |
|                         Latency p99.9 |     85.660538 ms |              1645.054495 ms |         434.422801 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |           4.9 GB |                      5.6 GB |                5.6 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |           3.8 GB |                      4.6 GB |                4.3 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |           270 MB |                      356 MB |                202 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |           1.5 GB |                      1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |         378.67 % |                    486.67 % |              478.47 % |
|               SERVER-MAX-MEMORY-USAGE |           1.4 GB |                       17 GB |                4.7 GB |
|                  CLIENT-MAX-CPU-USAGE |         430.00 % |                    553.50 % |              226.00 % |
|               CLIENT-MAX-MEMORY-USAGE |           302 MB |                      3.5 GB |                192 MB |
|                    CLIENT-ERROR-COUNT |                0 |                           0 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |               80 |                         325 |                   245 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                0 |                           0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |        1,477,712 |                   1,207,585 |             3,470,004 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |       31,450,244 |                  29,392,728 |           100,575,608 |
|           SERVER-AVG-DISK-SPACE-USAGE |           2.4 GB |                      7.2 GB |                3.0 GB |
+---------------------------------------+------------------+-----------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/MAX-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READ-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/01-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, Best Throughput (etcd 1,000, Zookeeper 700, Consul 500 clients)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10 (GNU/Linux kernel 4.8.0-49-generic)
- `ulimit -n` is 120000
- etcd tip (Go 1.8.1, git SHA f4641accc34be80c255ff87673a1cb92342abedd)
- Zookeeper r3.5.3-beta
  - Java 8
  - javac 1.8.0_131
  - Java(TM) SE Runtime Environment (build 1.8.0_131-b11)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.131-b11, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v0.8.3 (Go 1.8.1)


```
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                                       | etcd-tip-go1.8.1 | zookeeper-r3.5.3-beta-java8 | consul-v0.8.3-go1.8.1 |
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                         TOTAL-SECONDS |      27.1129 sec |                 73.7509 sec |          104.7749 sec |
|                  TOTAL-REQUEST-NUMBER |        1,000,000 |                   1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |   39,310 req/sec |              25,609 req/sec |        16,187 req/sec |
|                        AVG-THROUGHPUT |   36,882 req/sec |              13,527 req/sec |         9,544 req/sec |
|                        MIN-THROUGHPUT |   11,832 req/sec |                   0 req/sec |         1,153 req/sec |
|                       FASTEST-LATENCY |        5.4015 ms |                   3.6710 ms |            17.5701 ms |
|                           AVG-LATENCY |       27.0628 ms |                  40.5812 ms |            52.3596 ms |
|                       SLOWEST-LATENCY |      259.4610 ms |                6730.8732 ms |           726.8671 ms |
|                           Latency p10 |     12.479063 ms |                10.913450 ms |          30.501150 ms |
|                           Latency p25 |     15.698429 ms |                13.308653 ms |          34.493203 ms |
|                           Latency p50 |     21.313025 ms |                22.341612 ms |          40.393695 ms |
|                           Latency p75 |     32.893052 ms |                38.831905 ms |          49.501634 ms |
|                           Latency p90 |     53.129533 ms |                55.978192 ms |          78.371086 ms |
|                           Latency p95 |     58.230130 ms |                86.862387 ms |         113.602111 ms |
|                           Latency p99 |     69.868009 ms |               396.844374 ms |         302.290504 ms |
|                         Latency p99.9 |     86.566285 ms |              1509.821322 ms |         428.773213 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |           5.0 GB |                      5.3 GB |                5.6 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |           3.8 GB |                      4.3 GB |                4.3 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |           278 MB |                      365 MB |                206 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |           1.4 GB |                      1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |         405.77 % |                    509.00 % |              467.00 % |
|               SERVER-MAX-MEMORY-USAGE |           1.2 GB |                       17 GB |                5.0 GB |
|                  CLIENT-MAX-CPU-USAGE |         446.00 % |                    199.00 % |              212.00 % |
|               CLIENT-MAX-MEMORY-USAGE |           229 MB |                      4.0 GB |                 94 MB |
|                    CLIENT-ERROR-COUNT |                0 |                       2,358 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                8 |                         270 |                    12 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                0 |                           0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           96,781 |                     201,371 |               455,890 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |       20,479,856 |                  31,624,496 |            46,815,008 |
|           SERVER-AVG-DISK-SPACE-USAGE |           2.6 GB |                      8.6 GB |                3.0 GB |
+---------------------------------------+------------------+-----------------------------+-----------------------+


zookeeper__r3_5_3_beta errors:
"zk: could not connect to a server" (count 808)
"zk: connection closed" (count 1,550)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/MAX-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/02-write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, 100 clients, 1000 QPS Limit

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10 (GNU/Linux kernel 4.8.0-49-generic)
- `ulimit -n` is 120000
- etcd tip (Go 1.8.1, git SHA f4641accc34be80c255ff87673a1cb92342abedd)
- Zookeeper r3.5.3-beta
  - Java 8
  - javac 1.8.0_131
  - Java(TM) SE Runtime Environment (build 1.8.0_131-b11)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.131-b11, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v0.8.3 (Go 1.8.1)


```
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                                       | etcd-tip-go1.8.1 | zookeeper-r3.5.3-beta-java8 | consul-v0.8.3-go1.8.1 |
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                         TOTAL-SECONDS |     999.0049 sec |               1000.5826 sec |         1013.0493 sec |
|                  TOTAL-REQUEST-NUMBER |        1,000,000 |                   1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |    1,889 req/sec |               1,758 req/sec |         2,095 req/sec |
|                        AVG-THROUGHPUT |    1,000 req/sec |                 999 req/sec |           987 req/sec |
|                        MIN-THROUGHPUT |      111 req/sec |                   0 req/sec |            66 req/sec |
|                       FASTEST-LATENCY |        1.0461 ms |                   1.0048 ms |             3.7279 ms |
|                           AVG-LATENCY |        4.5826 ms |                   2.8069 ms |            49.6017 ms |
|                       SLOWEST-LATENCY |       49.1574 ms |                2309.1645 ms |          1467.9840 ms |
|                           Latency p10 |      2.660425 ms |                 1.673224 ms |           6.966911 ms |
|                           Latency p25 |      3.375709 ms |                 1.929275 ms |           9.864446 ms |
|                           Latency p50 |      4.379026 ms |                 2.243025 ms |          38.349274 ms |
|                           Latency p75 |      5.552898 ms |                 2.502263 ms |          84.282425 ms |
|                           Latency p90 |      6.651107 ms |                 2.702043 ms |         100.327132 ms |
|                           Latency p95 |      7.509832 ms |                 2.846511 ms |         108.347549 ms |
|                           Latency p99 |      9.331735 ms |                 3.849668 ms |         135.330942 ms |
|                         Latency p99.9 |     14.513159 ms |               115.777978 ms |         899.092677 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |           5.7 GB |                      5.7 GB |                6.0 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |           4.5 GB |                      4.6 GB |                4.7 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |           264 MB |                      355 MB |                209 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |           1.5 GB |                      1.5 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |          77.23 % |                    269.33 % |              271.60 % |
|               SERVER-MAX-MEMORY-USAGE |           1.6 GB |                       17 GB |                4.3 GB |
|                  CLIENT-MAX-CPU-USAGE |          48.00 % |                     38.00 % |               51.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            90 MB |                      647 MB |                 75 MB |
|                    CLIENT-ERROR-COUNT |                0 |                           0 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |              115 |                         539 |                   213 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                0 |                           0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |        5,648,559 |                   6,255,113 |            10,048,916 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |       64,808,440 |                 116,857,296 |           256,791,472 |
|           SERVER-AVG-DISK-SPACE-USAGE |           2.5 GB |                       11 GB |                2.8 GB |
+---------------------------------------+------------------+-----------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-LATENCY-MS.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-THROUGHPUT.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/MAX-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-VMRSS-MB.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-SECTORS-READ-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-READ-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/03-write-1M-keys-1000QPS/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Read 3M same keys, 256-byte key, 1KB value, Best Throughput (etcd 1,000, Zookeeper 700, Consul 500 clients)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10 (GNU/Linux kernel 4.8.0-49-generic)
- `ulimit -n` is 120000
- etcd tip (Go 1.8.1, git SHA f4641accc34be80c255ff87673a1cb92342abedd)
- Zookeeper r3.5.3-beta
  - Java 8
  - javac 1.8.0_131
  - Java(TM) SE Runtime Environment (build 1.8.0_131-b11)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.131-b11, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v0.8.3 (Go 1.8.1)


```
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                                       | etcd-tip-go1.8.1 | zookeeper-r3.5.3-beta-java8 | consul-v0.8.3-go1.8.1 |
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                         TOTAL-SECONDS |      25.5649 sec |                 33.3964 sec |           42.7809 sec |
|                  TOTAL-REQUEST-NUMBER |        3,000,000 |                   3,000,000 |             3,000,000 |
|                        MAX-THROUGHPUT |  129,450 req/sec |             105,215 req/sec |        75,899 req/sec |
|                        AVG-THROUGHPUT |  117,348 req/sec |              89,830 req/sec |        70,124 req/sec |
|                        MIN-THROUGHPUT |   24,852 req/sec |              20,212 req/sec |        55,371 req/sec |
|                       FASTEST-LATENCY |        0.4735 ms |                   0.8851 ms |             0.3436 ms |
|                           AVG-LATENCY |        8.1089 ms |                   7.4485 ms |             6.1135 ms |
|                       SLOWEST-LATENCY |      878.4009 ms |                 226.4025 ms |            79.0516 ms |
|                           Latency p10 |      3.528392 ms |                 2.801653 ms |           2.347963 ms |
|                           Latency p25 |      4.954795 ms |                 3.205985 ms |           3.477169 ms |
|                           Latency p50 |      6.985102 ms |                 4.003191 ms |           5.524614 ms |
|                           Latency p75 |      8.392994 ms |                12.555739 ms |           7.989689 ms |
|                           Latency p90 |     12.082250 ms |                15.000391 ms |          10.686574 ms |
|                           Latency p95 |     15.594535 ms |                17.438278 ms |          12.588298 ms |
|                           Latency p99 |     35.948728 ms |                26.880671 ms |          16.373871 ms |
|                         Latency p99.9 |    223.739181 ms |                44.972673 ms |          22.810455 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |           1.3 GB |                      2.4 GB |                5.2 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |           4.8 GB |                      4.9 GB |                9.5 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |           4.7 GB |                      4.6 GB |                5.9 GB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |           1.3 GB |                      2.0 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |         669.27 % |                    622.67 % |              883.00 % |
|               SERVER-MAX-MEMORY-USAGE |            84 MB |                       11 GB |                 38 MB |
|                  CLIENT-MAX-CPU-USAGE |        1112.00 % |                   1386.00 % |             1358.00 % |
|               CLIENT-MAX-MEMORY-USAGE |           287 MB |                      4.5 GB |                158 MB |
|                    CLIENT-ERROR-COUNT |                0 |                           0 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                0 |                          15 |                     1 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                0 |                           0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |               81 |                       5,151 |                   182 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |              912 |                      26,496 |                 2,096 |
|           SERVER-AVG-DISK-SPACE-USAGE |            81 MB |                       67 MB |                 69 kB |
+---------------------------------------+------------------+-----------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/MAX-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/04-read-3M-same-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Read 3M same keys, 256-byte key, 1KB value, 1,000 clients

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 16.10 (GNU/Linux kernel 4.8.0-49-generic)
- `ulimit -n` is 120000
- etcd tip (Go 1.8.1, git SHA f4641accc34be80c255ff87673a1cb92342abedd)
- Zookeeper r3.5.3-beta
  - Java 8
  - javac 1.8.0_131
  - Java(TM) SE Runtime Environment (build 1.8.0_131-b11)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.131-b11, mixed mode)
  - `/usr/bin/java -Djute.maxbuffer=33554432 -Xms50G -Xmx50G`
- Consul v0.8.3 (Go 1.8.1)


```
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                                       | etcd-tip-go1.8.1 | zookeeper-r3.5.3-beta-java8 | consul-v0.8.3-go1.8.1 |
+---------------------------------------+------------------+-----------------------------+-----------------------+
|                         TOTAL-SECONDS |      25.5690 sec |                 32.5781 sec |           45.3777 sec |
|                  TOTAL-REQUEST-NUMBER |        3,000,000 |                   3,000,000 |             3,000,000 |
|                        MAX-THROUGHPUT |  131,842 req/sec |             108,889 req/sec |        70,944 req/sec |
|                        AVG-THROUGHPUT |  117,329 req/sec |              92,078 req/sec |        66,111 req/sec |
|                        MIN-THROUGHPUT |   25,889 req/sec |               4,444 req/sec |           902 req/sec |
|                       FASTEST-LATENCY |        0.4944 ms |                   0.5857 ms |             0.3899 ms |
|                           AVG-LATENCY |        8.1291 ms |                  10.0708 ms |            13.2565 ms |
|                       SLOWEST-LATENCY |      888.7094 ms |                 499.0320 ms |          1073.5244 ms |
|                           Latency p10 |      3.644467 ms |                 3.666637 ms |           4.953445 ms |
|                           Latency p25 |      4.999142 ms |                 4.307426 ms |           7.085859 ms |
|                           Latency p50 |      6.981474 ms |                 5.416655 ms |          11.658595 ms |
|                           Latency p75 |      8.258696 ms |                17.228526 ms |          17.624440 ms |
|                           Latency p90 |     11.807836 ms |                19.635761 ms |          24.435628 ms |
|                           Latency p95 |     15.451811 ms |                20.759451 ms |          28.323236 ms |
|                           Latency p99 |     36.775950 ms |                23.300235 ms |          34.916833 ms |
|                         Latency p99.9 |    223.901331 ms |               223.894715 ms |          42.851037 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |           1.3 GB |                      2.5 GB |                5.2 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |           4.8 GB |                      5.1 GB |                9.6 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |           4.7 GB |                      4.6 GB |                5.9 GB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |           1.3 GB |                      2.0 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |         665.00 % |                    693.67 % |              848.63 % |
|               SERVER-MAX-MEMORY-USAGE |            84 MB |                       12 GB |                 52 MB |
|                  CLIENT-MAX-CPU-USAGE |        1122.00 % |                   1306.00 % |             1261.00 % |
|               CLIENT-MAX-MEMORY-USAGE |           311 MB |                      6.3 GB |                181 MB |
|                    CLIENT-ERROR-COUNT |                0 |                         272 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                0 |                           9 |                     0 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                0 |                           0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |               47 |                       3,908 |                   342 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |              504 |                      23,192 |                 3,328 |
|           SERVER-AVG-DISK-SPACE-USAGE |            81 MB |                       67 MB |                 69 kB |
+---------------------------------------+------------------+-----------------------------+-----------------------+


zookeeper__r3_5_3_beta errors:
"zk: could not connect to a server" (count 272)
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-LATENCY-MS.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-LATENCY-MS-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-THROUGHPUT.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/MAX-CPU.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-VMRSS-MB.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-VMRSS-MB-BY-KEY.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-SECTORS-READ-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-READ-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q2-01-etcd-zookeeper-consul/05-read-3M-same-keys-1K-client/AVG-TRANSMIT-BYTES-NUM-DELTA">



