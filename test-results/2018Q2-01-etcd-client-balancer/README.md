

<br><br><hr>
##### Read 3M same keys, 256-byte key, 1KB value, Best Throughput (etcd 1K clients with 100 conns)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 17.10 (GNU/Linux kernel 4.13.0-41-generic)
- `ulimit -n` is 120000
- etcd v3.2.20 (Go 1.8.7)
- etcd v3.3.5 (Go 1.9.6)
- etcd v3.4 67b1ff672 (Go 1.10.2)
- etcd v3.4 new balancer (Go 1.10.2)


```
+---------------------------------------+----------------------+---------------------+------------------------------+---------------------------------+
|                                       | etcd-v3.2.20-go1.8.7 | etcd-v3.3.5-go1.9.6 | etcd-v3.4-67b1ff672-go1.10.2 | etcd-v3.4-balancer0511-go1.10.2 |
+---------------------------------------+----------------------+---------------------+------------------------------+---------------------------------+
|                         TOTAL-SECONDS |          18.1562 sec |         18.3415 sec |                  18.0028 sec |                     17.8095 sec |
|                  TOTAL-REQUEST-NUMBER |            3,000,000 |           3,000,000 |                    3,000,000 |                       3,000,000 |
|                        MAX-THROUGHPUT |      176,854 req/sec |     176,783 req/sec |              179,961 req/sec |                 187,593 req/sec |
|                        AVG-THROUGHPUT |      165,233 req/sec |     163,563 req/sec |              166,640 req/sec |                 168,449 req/sec |
|                        MIN-THROUGHPUT |       62,416 req/sec |      15,304 req/sec |               64,844 req/sec |                  42,215 req/sec |
|                       FASTEST-LATENCY |            0.4388 ms |           0.4826 ms |                    0.4670 ms |                       0.4871 ms |
|                           AVG-LATENCY |            4.2135 ms |           4.3606 ms |                    4.4632 ms |                       4.5614 ms |
|                       SLOWEST-LATENCY |           61.7597 ms |         215.6992 ms |                   38.3503 ms |                      39.4441 ms |
|                           Latency p10 |          1.661360 ms |         1.752339 ms |                  1.763714 ms |                     1.892650 ms |
|                           Latency p25 |          2.307997 ms |         2.379004 ms |                  2.416912 ms |                     2.545477 ms |
|                           Latency p50 |          3.593472 ms |         3.651454 ms |                  3.811588 ms |                     3.966304 ms |
|                           Latency p75 |          5.587801 ms |         5.735338 ms |                  5.994316 ms |                     6.080458 ms |
|                           Latency p90 |          7.469498 ms |         7.795259 ms |                  7.904383 ms |                     7.884734 ms |
|                           Latency p95 |          8.925745 ms |         9.357837 ms |                  9.325560 ms |                     9.214828 ms |
|                           Latency p99 |         12.104650 ms |        13.055879 ms |                 12.643542 ms |                    12.341674 ms |
|                         Latency p99.9 |         18.393202 ms |        19.234115 ms |                 18.204708 ms |                    20.965550 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |               1.2 GB |              1.2 GB |                       1.2 GB |                          1.2 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |               4.5 GB |              4.5 GB |                       4.5 GB |                          4.5 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |               4.4 GB |              4.7 GB |                       4.4 GB |                          4.4 GB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |               1.1 GB |              1.2 GB |                       1.2 GB |                          1.2 GB |
|                  SERVER-MAX-CPU-USAGE |             941.00 % |            914.67 % |                     913.67 % |                        860.00 % |
|               SERVER-MAX-MEMORY-USAGE |                52 MB |               54 MB |                        59 MB |                           57 MB |
|                  CLIENT-MAX-CPU-USAGE |            1456.00 % |           1477.00 % |                    1461.00 % |                       1438.00 % |
|               CLIENT-MAX-MEMORY-USAGE |               166 MB |              185 MB |                       171 MB |                          160 MB |
|                    CLIENT-ERROR-COUNT |                    0 |                   0 |                            0 |                               0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                   27 |                  16 |                            0 |                               0 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                    0 |                   0 |                            0 |                               0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |                   30 |                 120 |                           91 |                             115 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |                  400 |               1,280 |                        1,016 |                           1,288 |
|           SERVER-AVG-DISK-SPACE-USAGE |                81 MB |               64 MB |                        64 MB |                           64 MB |
+---------------------------------------+----------------------+---------------------+------------------------------+---------------------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-CPU.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/MAX-CPU.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, Best Throughput (etcd 1K clients with 100 conns)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 17.10 (GNU/Linux kernel 4.13.0-41-generic)
- `ulimit -n` is 120000
- etcd v3.2.20 (Go 1.8.7)
- etcd v3.3.5 (Go 1.9.6)
- etcd v3.4 67b1ff672 (Go 1.10.2)
- etcd v3.4 new balancer (Go 1.10.2)


```
+---------------------------------------+----------------------+---------------------+------------------------------+---------------------------------+
|                                       | etcd-v3.2.20-go1.8.7 | etcd-v3.3.5-go1.9.6 | etcd-v3.4-67b1ff672-go1.10.2 | etcd-v3.4-balancer0511-go1.10.2 |
+---------------------------------------+----------------------+---------------------+------------------------------+---------------------------------+
|                         TOTAL-SECONDS |          28.3191 sec |         29.8842 sec |                  28.1384 sec |                     28.7298 sec |
|                  TOTAL-REQUEST-NUMBER |            1,000,000 |           1,000,000 |                    1,000,000 |                       1,000,000 |
|                        MAX-THROUGHPUT |       37,372 req/sec |      36,660 req/sec |               37,100 req/sec |                  37,574 req/sec |
|                        AVG-THROUGHPUT |       35,311 req/sec |      33,462 req/sec |               35,538 req/sec |                  34,807 req/sec |
|                        MIN-THROUGHPUT |       12,599 req/sec |      11,135 req/sec |               18,686 req/sec |                   1,353 req/sec |
|                       FASTEST-LATENCY |            4.4390 ms |           4.5028 ms |                    4.4703 ms |                       1.5290 ms |
|                           AVG-LATENCY |           28.2084 ms |          29.7919 ms |                   28.0337 ms |                      28.6681 ms |
|                       SLOWEST-LATENCY |          145.3477 ms |         121.5527 ms |                  120.5386 ms |                     122.6241 ms |
|                           Latency p10 |         11.936935 ms |        12.359409 ms |                 11.845129 ms |                    11.773927 ms |
|                           Latency p25 |         14.582232 ms |        15.363413 ms |                 14.208263 ms |                    16.379649 ms |
|                           Latency p50 |         21.711737 ms |        23.324211 ms |                 21.912266 ms |                    23.166898 ms |
|                           Latency p75 |         35.140638 ms |        38.881637 ms |                 34.065757 ms |                    39.052418 ms |
|                           Latency p90 |         57.640700 ms |        61.135803 ms |                 59.341919 ms |                    53.821186 ms |
|                           Latency p95 |         65.274153 ms |        67.717201 ms |                 65.998481 ms |                    60.028769 ms |
|                           Latency p99 |         83.692351 ms |        85.296525 ms |                 81.303646 ms |                    74.864572 ms |
|                         Latency p99.9 |        115.552287 ms |       109.913571 ms |                104.323381 ms |                   102.394398 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |               4.9 GB |              5.1 GB |                       5.0 GB |                          5.3 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |               3.7 GB |              3.8 GB |                       3.7 GB |                          4.0 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |               252 MB |              256 MB |                       254 MB |                          329 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |               1.5 GB |              1.5 GB |                       1.5 GB |                          1.6 GB |
|                  SERVER-MAX-CPU-USAGE |             491.83 % |            477.67 % |                     504.67 % |                        540.00 % |
|               SERVER-MAX-MEMORY-USAGE |               1.1 GB |              1.2 GB |                       1.2 GB |                          1.2 GB |
|                  CLIENT-MAX-CPU-USAGE |             622.00 % |            625.70 % |                     663.00 % |                        673.00 % |
|               CLIENT-MAX-MEMORY-USAGE |                97 MB |              101 MB |                        93 MB |                          104 MB |
|                    CLIENT-ERROR-COUNT |                    0 |                   0 |                            0 |                               0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                    0 |                   2 |                            1 |                               0 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                    0 |                   0 |                            0 |                               0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |              104,308 |             105,480 |                      103,923 |                         113,568 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |           20,456,512 |          20,718,432 |                   20,575,208 |                      20,831,256 |
|           SERVER-AVG-DISK-SPACE-USAGE |               2.7 GB |              2.7 GB |                       2.6 GB |                          2.6 GB |
+---------------------------------------+----------------------+---------------------+------------------------------+---------------------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-CPU.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/MAX-CPU.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2018Q2-01-etcd-client-balancer/write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">



