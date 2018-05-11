

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
|                         TOTAL-SECONDS |          27.9457 sec |         27.8376 sec |                  29.8214 sec |                     31.7905 sec |
|                  TOTAL-REQUEST-NUMBER |            1,000,000 |           1,000,000 |                    1,000,000 |                       1,000,000 |
|                        MAX-THROUGHPUT |       37,870 req/sec |      38,399 req/sec |               35,302 req/sec |                  34,293 req/sec |
|                        AVG-THROUGHPUT |       35,783 req/sec |      35,922 req/sec |               33,532 req/sec |                  31,455 req/sec |
|                        MIN-THROUGHPUT |       33,308 req/sec |      30,366 req/sec |                4,518 req/sec |                  25,404 req/sec |
|                       FASTEST-LATENCY |            5.0600 ms |           5.0021 ms |                    6.6801 ms |                       6.5805 ms |
|                           AVG-LATENCY |           27.8413 ms |          27.7460 ms |                   29.7391 ms |                      31.7049 ms |
|                       SLOWEST-LATENCY |          137.2510 ms |         123.4449 ms |                  137.3843 ms |                     122.6154 ms |
|                           Latency p10 |         12.288933 ms |        11.079009 ms |                 13.443536 ms |                    14.645595 ms |
|                           Latency p25 |         15.050926 ms |        13.415975 ms |                 16.699900 ms |                    17.839535 ms |
|                           Latency p50 |         21.019031 ms |        20.939641 ms |                 23.205842 ms |                    25.704132 ms |
|                           Latency p75 |         35.016191 ms |        35.073160 ms |                 36.362920 ms |                    41.522521 ms |
|                           Latency p90 |         55.658405 ms |        59.034597 ms |                 58.360949 ms |                    60.267800 ms |
|                           Latency p95 |         63.626096 ms |        66.795356 ms |                 67.065177 ms |                    67.282624 ms |
|                           Latency p99 |         82.812214 ms |        86.663292 ms |                 81.502701 ms |                    81.240145 ms |
|                         Latency p99.9 |        110.040242 ms |       109.186148 ms |                108.321180 ms |                   101.836643 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |               4.9 GB |              4.8 GB |                       5.1 GB |                          4.9 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |               3.7 GB |              3.6 GB |                       3.8 GB |                          3.7 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |               4.7 GB |              4.9 GB |                       258 MB |                          254 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |               2.7 GB |              2.7 GB |                       1.5 GB |                          1.5 GB |
|                  SERVER-MAX-CPU-USAGE |             487.33 % |            475.67 % |                     490.00 % |                        496.67 % |
|               SERVER-MAX-MEMORY-USAGE |               1.1 GB |              1.1 GB |                       1.1 GB |                          1.2 GB |
|                  CLIENT-MAX-CPU-USAGE |            1456.00 % |           1477.00 % |                     630.00 % |                        562.00 % |
|               CLIENT-MAX-MEMORY-USAGE |               166 MB |              185 MB |                        88 MB |                           90 MB |
|                    CLIENT-ERROR-COUNT |                    0 |                   0 |                            0 |                               0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                   44 |                  26 |                            3 |                               1 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                    0 |                   0 |                            0 |                               0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |              103,627 |             101,928 |                      104,389 |                          99,307 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |           20,009,512 |          20,008,880 |                   20,716,200 |                      20,157,016 |
|           SERVER-AVG-DISK-SPACE-USAGE |               2.6 GB |              2.6 GB |                       2.7 GB |                          2.8 GB |
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



