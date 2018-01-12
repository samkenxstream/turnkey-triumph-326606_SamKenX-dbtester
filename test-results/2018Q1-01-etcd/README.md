
<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-3m-avg-latency.png" alt="write-3m-avg-latency">

<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-3m-avg-latency-by-key.png" alt="write-3m-avg-latency-by-key">

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





<br><br><hr>
##### Write 3-million keys, 256-byte key, 1KB value, Best Throughput (etcd 1,000)

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
|                         TOTAL-SECONDS |         88.5767 sec |         84.0750 sec |
|                  TOTAL-REQUEST-NUMBER |           3,000,000 |           3,000,000 |
|                        MAX-THROUGHPUT |      37,110 req/sec |      39,478 req/sec |
|                        AVG-THROUGHPUT |      32,976 req/sec |      35,682 req/sec |
|                        MIN-THROUGHPUT |      16,075 req/sec |       8,116 req/sec |
|                       FASTEST-LATENCY |           3.5655 ms |           3.1295 ms |
|                           AVG-LATENCY |          30.1489 ms |          27.9304 ms |
|                       SLOWEST-LATENCY |         346.7163 ms |         270.7468 ms |
|                           Latency p10 |        12.991756 ms |        13.338154 ms |
|                           Latency p25 |        15.870495 ms |        16.324596 ms |
|                           Latency p50 |        22.419681 ms |        21.073518 ms |
|                           Latency p75 |        37.436826 ms |        33.602292 ms |
|                           Latency p90 |        58.439201 ms |        53.845497 ms |
|                           Latency p95 |        68.347220 ms |        59.987225 ms |
|                           Latency p99 |       113.337461 ms |        87.168222 ms |
|                         Latency p99.9 |       185.483190 ms |       134.969831 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |               15 GB |               16 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |               11 GB |               12 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |              1.2 GB |              1.0 GB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |              4.7 GB |              4.8 GB |
|                  SERVER-MAX-CPU-USAGE |            602.33 % |            550.00 % |
|               SERVER-MAX-MEMORY-USAGE |              2.7 GB |              2.9 GB |
|                  CLIENT-MAX-CPU-USAGE |            896.00 % |            713.00 % |
|               CLIENT-MAX-MEMORY-USAGE |              323 MB |              316 MB |
|                    CLIENT-ERROR-COUNT |              79,079 |                   0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                  62 |                   7 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                   0 |                   0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |             308,267 |             333,907 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |          60,593,932 |          62,305,992 |
|           SERVER-AVG-DISK-SPACE-USAGE |              6.6 GB |              6.6 GB |
+---------------------------------------+---------------------+---------------------+


etcd__v3_2 errors:
"etcdserver: too many requests" (count 79,079)
```


<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-LATENCY-MS.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-LATENCY-MS-BY-KEY.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-THROUGHPUT.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-CPU.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/MAX-CPU.svg" alt="2018Q1-01-etcd/write-too-many-keys/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-VMRSS-MB.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-VMRSS-MB-BY-KEY.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-READS-COMPLETED-DELTA.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-SECTORS-READ-DELTA.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-WRITES-COMPLETED-DELTA.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-READ-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd/write-too-many-keys/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/write-too-many-keys/AVG-TRANSMIT-BYTES-NUM-DELTA">




<br><br><hr>
##### Read 3M same keys, 256-byte key, 1KB value, Best Throughput (etcd 1K clients with 100 conns)

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
|                         TOTAL-SECONDS |         19.7093 sec |         19.1632 sec |
|                  TOTAL-REQUEST-NUMBER |           3,000,000 |           3,000,000 |
|                        MAX-THROUGHPUT |     167,224 req/sec |     171,500 req/sec |
|                        AVG-THROUGHPUT |     152,212 req/sec |     156,549 req/sec |
|                        MIN-THROUGHPUT |       7,792 req/sec |      87,028 req/sec |
|                       FASTEST-LATENCY |           0.5423 ms |           0.5970 ms |
|                           AVG-LATENCY |           5.3364 ms |           4.8534 ms |
|                       SLOWEST-LATENCY |         659.9911 ms |         215.0872 ms |
|                           Latency p10 |         2.132814 ms |         2.028502 ms |
|                           Latency p25 |         2.851045 ms |         2.705135 ms |
|                           Latency p50 |         4.488098 ms |         4.167273 ms |
|                           Latency p75 |         6.742433 ms |         6.396360 ms |
|                           Latency p90 |         8.985732 ms |         8.496297 ms |
|                           Latency p95 |        10.723795 ms |        10.039053 ms |
|                           Latency p99 |        15.161416 ms |        13.638065 ms |
|                         Latency p99.9 |        33.201151 ms |        20.599315 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |              1.2 GB |              1.2 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |              4.8 GB |              4.5 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |              4.7 GB |              4.4 GB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |              1.2 GB |              1.2 GB |
|                  SERVER-MAX-CPU-USAGE |            824.00 % |            908.67 % |
|               SERVER-MAX-MEMORY-USAGE |               50 MB |               54 MB |
|                  CLIENT-MAX-CPU-USAGE |           1438.00 % |           1461.00 % |
|               CLIENT-MAX-MEMORY-USAGE |              171 MB |              167 MB |
|                    CLIENT-ERROR-COUNT |                   0 |                   0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                   0 |                   0 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                   0 |                   0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |                  97 |                  30 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |               1,280 |                 376 |
|           SERVER-AVG-DISK-SPACE-USAGE |               81 MB |               64 MB |
+---------------------------------------+---------------------+---------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-LATENCY-MS.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-LATENCY-MS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-LATENCY-MS-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-THROUGHPUT.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-THROUGHPUT">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-NON-VOLUNTARY-CTXT-SWITCHES">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-CPU.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/MAX-CPU.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/MAX-CPU">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-VMRSS-MB.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-VMRSS-MB">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-VMRSS-MB-BY-KEY-ERROR-POINTS">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-READS-COMPLETED-DELTA.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-READS-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-SECTORS-READ-DELTA.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-SECTORS-READ-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-WRITES-COMPLETED-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-SECTORS-WRITTEN-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-READ-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-WRITE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-RECEIVE-BYTES-NUM-DELTA">



<img src="https://storage.googleapis.com/dbtester-results/2018Q1-01-etcd-read-3M-same-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2018Q1-01-etcd/read-3M-same-keys-best-throughput/AVG-TRANSMIT-BYTES-NUM-DELTA">



