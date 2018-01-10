

<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value, Best Throughput (etcd 1K clients with 100 conns)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 60 GB Memory + 300 GB SSD (1 for client)
- Ubuntu 17.10 (GNU/Linux kernel 4.13.0-25-generic)
- `ulimit -n` is 120000
- etcd v3.2.0 (Go 1.8.3)
- etcd v3.3.0 (Go 1.9.2)


```
+---------------------------------------+---------------------+---------------------+
|                                       | etcd-v3.2.0-go1.8.3 | etcd-v3.3.0-go1.9.2 |
+---------------------------------------+---------------------+---------------------+
|                         TOTAL-SECONDS |         28.2484 sec |         27.6024 sec |
|                  TOTAL-REQUEST-NUMBER |           1,000,000 |           1,000,000 |
|                        MAX-THROUGHPUT |      36,917 req/sec |      38,830 req/sec |
|                        AVG-THROUGHPUT |      35,400 req/sec |      36,228 req/sec |
|                        MIN-THROUGHPUT |       2,826 req/sec |      23,191 req/sec |
|                       FASTEST-LATENCY |           5.4766 ms |           5.0665 ms |
|                           AVG-LATENCY |          28.1631 ms |          27.5046 ms |
|                       SLOWEST-LATENCY |         169.7880 ms |         137.9128 ms |
|                           Latency p10 |        13.171913 ms |        12.482529 ms |
|                           Latency p25 |        16.340837 ms |        15.299216 ms |
|                           Latency p50 |        21.839763 ms |        21.250243 ms |
|                           Latency p75 |        34.903072 ms |        32.857025 ms |
|                           Latency p90 |        54.765537 ms |        55.143954 ms |
|                           Latency p95 |        60.558651 ms |        61.219371 ms |
|                           Latency p99 |        78.386660 ms |        75.881935 ms |
|                         Latency p99.9 |       100.343567 ms |       100.020188 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |              5.1 GB |              5.0 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |              3.9 GB |              3.7 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |              333 MB |              251 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |              1.5 GB |              1.5 GB |
|                  SERVER-MAX-CPU-USAGE |            443.00 % |            485.00 % |
|               SERVER-MAX-MEMORY-USAGE |              1.1 GB |              1.2 GB |
|                  CLIENT-MAX-CPU-USAGE |            597.00 % |            639.00 % |
|               CLIENT-MAX-MEMORY-USAGE |               82 MB |               93 MB |
|                    CLIENT-ERROR-COUNT |                   0 |                   0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                   2 |                   1 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                   0 |                   0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |             107,342 |             104,710 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |          20,746,024 |          20,265,400 |
|           SERVER-AVG-DISK-SPACE-USAGE |              2.6 GB |              2.6 GB |
+---------------------------------------+---------------------+---------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-CPU.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/MAX-CPU.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/write-1M-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">



