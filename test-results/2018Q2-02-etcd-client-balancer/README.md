

<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, Best Throughput (etcd 1K clients with 100 conns)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 17.10 (GNU/Linux kernel 4.13.0-45-generic)
- `ulimit -n` is 120000
- etcd v3.4 b241e383 (Go 1.10.3)
- etcd v3.4 new balancer (Go 1.10.3)


```
+---------------------------------------+-----------------------------+---------------------------------+
|                                       | etcd-v3.4-b241e383-go1.10.3 | etcd-v3.4-balancer0615-go1.10.3 |
+---------------------------------------+-----------------------------+---------------------------------+
|                         TOTAL-SECONDS |                 31.1256 sec |                     31.2477 sec |
|                  TOTAL-REQUEST-NUMBER |                   1,000,000 |                       1,000,000 |
|                        MAX-THROUGHPUT |              33,760 req/sec |                  34,587 req/sec |
|                        AVG-THROUGHPUT |              32,127 req/sec |                  32,002 req/sec |
|                        MIN-THROUGHPUT |               4,965 req/sec |                  10,454 req/sec |
|                       FASTEST-LATENCY |                   4.6587 ms |                       2.4888 ms |
|                           AVG-LATENCY |                  31.0604 ms |                      31.2033 ms |
|                       SLOWEST-LATENCY |                 117.5620 ms |                     114.9492 ms |
|                           Latency p10 |                13.431526 ms |                    14.691959 ms |
|                           Latency p25 |                17.993337 ms |                    19.586467 ms |
|                           Latency p50 |                24.734914 ms |                    25.571253 ms |
|                           Latency p75 |                42.801499 ms |                    42.138762 ms |
|                           Latency p90 |                57.777309 ms |                    55.289961 ms |
|                           Latency p95 |                65.311487 ms |                    60.855029 ms |
|                           Latency p99 |                78.819013 ms |                    75.192049 ms |
|                         Latency p99.9 |                97.808156 ms |                    92.254135 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |                      5.2 GB |                          5.2 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |                      3.9 GB |                          4.0 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |                      258 MB |                          324 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |                      1.5 GB |                          1.6 GB |
|                  SERVER-MAX-CPU-USAGE |                    440.30 % |                        537.67 % |
|               SERVER-MAX-MEMORY-USAGE |                      1.2 GB |                          1.2 GB |
|                  CLIENT-MAX-CPU-USAGE |                    570.00 % |                        593.00 % |
|               CLIENT-MAX-MEMORY-USAGE |                       95 MB |                          171 MB |
|                    CLIENT-ERROR-COUNT |                           0 |                               0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                           0 |                              73 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                           0 |                               0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |                     103,846 |                         109,864 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |                  23,873,928 |                      20,586,688 |
|           SERVER-AVG-DISK-SPACE-USAGE |                      2.7 GB |                          2.7 GB |
+---------------------------------------+-----------------------------+---------------------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-CPU.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/MAX-CPU.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Read 3M same keys, 256-byte key, 1KB value, Best Throughput (etcd 1K clients with 100 conns)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 17.10 (GNU/Linux kernel 4.13.0-45-generic)
- `ulimit -n` is 120000
- etcd v3.4 b241e383 (Go 1.10.3)
- etcd v3.4 new balancer (Go 1.10.3)


```
+---------------------------------------+-----------------------------+---------------------------------+
|                                       | etcd-v3.4-b241e383-go1.10.3 | etcd-v3.4-balancer0615-go1.10.3 |
+---------------------------------------+-----------------------------+---------------------------------+
|                         TOTAL-SECONDS |                 17.8744 sec |                     17.8226 sec |
|                  TOTAL-REQUEST-NUMBER |                   3,000,000 |                       3,000,000 |
|                        MAX-THROUGHPUT |             176,763 req/sec |                 172,164 req/sec |
|                        AVG-THROUGHPUT |             167,837 req/sec |                 168,325 req/sec |
|                        MIN-THROUGHPUT |              38,290 req/sec |                   7,453 req/sec |
|                       FASTEST-LATENCY |                   0.5131 ms |                       0.5025 ms |
|                           AVG-LATENCY |                   4.6043 ms |                       4.6358 ms |
|                       SLOWEST-LATENCY |                  37.8623 ms |                      29.7872 ms |
|                           Latency p10 |                 1.729814 ms |                     2.372096 ms |
|                           Latency p25 |                 2.383698 ms |                     3.036887 ms |
|                           Latency p50 |                 3.961112 ms |                     4.055946 ms |
|                           Latency p75 |                 6.137971 ms |                     5.684766 ms |
|                           Latency p90 |                 8.458589 ms |                     7.767217 ms |
|                           Latency p95 |                10.006860 ms |                     9.068512 ms |
|                           Latency p99 |                13.232563 ms |                    12.085174 ms |
|                         Latency p99.9 |                18.042299 ms |                    16.128133 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |                      1.2 GB |                          1.3 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |                      4.5 GB |                          4.6 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |                      4.4 GB |                          4.8 GB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |                      1.2 GB |                          1.3 GB |
|                  SERVER-MAX-CPU-USAGE |                    891.33 % |                        867.33 % |
|               SERVER-MAX-MEMORY-USAGE |                       58 MB |                           68 MB |
|                  CLIENT-MAX-CPU-USAGE |                   1453.00 % |                       1510.00 % |
|               CLIENT-MAX-MEMORY-USAGE |                      158 MB |                          255 MB |
|                    CLIENT-ERROR-COUNT |                           0 |                               0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                           0 |                               0 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                           0 |                               0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |                          51 |                              96 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |                         448 |                           1,112 |
|           SERVER-AVG-DISK-SPACE-USAGE |                       64 MB |                           64 MB |
+---------------------------------------+-----------------------------+---------------------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-CPU.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/MAX-CPU.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2018Q2-02-etcd-client-balancer/read-3M-same-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">



