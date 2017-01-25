# dbtester

[![Build Status](https://img.shields.io/travis/coreos/dbtester.svg?style=flat-square)](https://travis-ci.org/coreos/dbtester) [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/coreos/dbtester)

Distributed database benchmark tester: etcd, Zookeeper, Consul

- Database agent and runner are implemented at https://github.com/coreos/dbtester/tree/master/agent
- Client is implemented at https://github.com/coreos/dbtester/tree/master/control
- System metrics are collected via https://github.com/gyuho/psn
- Data analysis is done via https://github.com/coreos/dbtester/tree/master/analyze
  - https://github.com/gyuho/dataframe
  - https://github.com/gonum/plot

For etcd, we also recommend [etcd benchmark tool](https://github.com/coreos/etcd/tree/master/tools/benchmark).

All logs and results can be found at https://console.cloud.google.com/storage/browser/dbtester-results


<br><br><hr>
##### Write 1M keys, 1000-client, 256-byte key, 1KB value

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 30 GB Memory + 150 GB SSD (1 for client)
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.4)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.2 (Go 1.7.4)


```
+----------------------------+--------------------+------------------------+-----------------------+
|                            | etcd-v3.1-go1.7.4  | zookeeper-r3.4.9-java8 | consul-v0.7.2-go1.7.4 |
+----------------------------+--------------------+------------------------+-----------------------+
|      READS-COMPLETED-DELTA |                  6 |                    311 |                    15 |
|    SECTORS-READS-DELTA-SUM |                  0 |                      0 |                     0 |
| WRITES-COMPLETED-DELTA-SUM |              96474 |                  84881 |                940695 |
|  SECTORS-WRITTEN-DELTA-SUM |             542512 |                9221640 |              41272068 |
|          RECEIVE-BYTES-SUM |             4.9 GB |                 5.3 GB |                7.7 GB |
|         TRANSMIT-BYTES-SUM |             3.7 GB |                 4.2 GB |                6.5 GB |
|              MAX-CPU-USAGE |           291.56 % |               363.65 % |              226.18 % |
|           MAX-MEMORY-USAGE |         1198.60 MB |             4688.05 MB |            4329.38 MB |
|              TOTAL-SECONDS |        36.2024 sec |            61.0944 sec |          467.9311 sec |
|             AVG-THROUGHPUT | 27622.4453 req/sec |     16298.0557 req/sec |     2137.0667 req/sec |
|            SLOWEST-LATENCY |        246.4560 ms |           5570.6375 ms |         30388.9318 ms |
|            FASTEST-LATENCY |          5.3413 ms |              2.4757 ms |            21.5605 ms |
|                AVG-LATENCY |         36.1057 ms |             50.4279 ms |           467.4253 ms |
|                Latency p10 |       13.712090 ms |           14.861507 ms |          65.910086 ms |
|                Latency p25 |       16.625779 ms |           18.884719 ms |          77.221971 ms |
|                Latency p50 |       22.306160 ms |           22.291879 ms |         120.663354 ms |
|                Latency p75 |       40.376905 ms |           25.751846 ms |         716.373543 ms |
|                Latency p90 |       65.849751 ms |           30.030446 ms |        1068.038406 ms |
|                Latency p95 |      137.545464 ms |           81.141780 ms |        1080.751412 ms |
|                Latency p99 |      177.127309 ms |          965.771377 ms |        2686.919571 ms |
|              Latency p99.9 |      198.540415 ms |         2911.408642 ms |       19041.188919 ms |
+----------------------------+--------------------+------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/01-write-1M-keys/AVG-TRANSMIT-BYTES-NUM-DELTA">



<br><br><hr>
##### Write 1M keys, 1000-client, 1000QPS, 256-byte key, 1KB value

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 30 GB Memory + 150 GB SSD (1 for client)
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.4)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.2 (Go 1.7.4)


```
+----------------------------+-------------------+------------------------+-----------------------+
|                            | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.2-go1.7.4 |
+----------------------------+-------------------+------------------------+-----------------------+
|      READS-COMPLETED-DELTA |                 0 |                    205 |                   141 |
|    SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| WRITES-COMPLETED-DELTA-SUM |           4429435 |                4666587 |               7230463 |
|  SECTORS-WRITTEN-DELTA-SUM |           2262752 |               14830288 |              81718092 |
|          RECEIVE-BYTES-SUM |            5.8 GB |                 5.8 GB |                6.0 GB |
|         TRANSMIT-BYTES-SUM |            4.6 GB |                 4.7 GB |                4.7 GB |
|              MAX-CPU-USAGE |           52.02 % |                46.67 % |               84.40 % |
|           MAX-MEMORY-USAGE |        1626.60 MB |             3568.49 MB |            3936.74 MB |
|              TOTAL-SECONDS |      999.0083 sec |          1000.5091 sec |         1195.0587 sec |
|             AVG-THROUGHPUT | 1000.9927 req/sec |       998.8554 req/sec |      836.7790 req/sec |
|            SLOWEST-LATENCY |       198.1010 ms |           2392.6933 ms |         23594.1354 ms |
|            FASTEST-LATENCY |         1.2003 ms |              0.9210 ms |             3.6024 ms |
|                AVG-LATENCY |         4.5224 ms |              8.1224 ms |           363.4513 ms |
|                Latency p10 |       2.710760 ms |            1.593791 ms |           8.339832 ms |
|                Latency p25 |       3.382362 ms |            1.815819 ms |          36.060302 ms |
|                Latency p50 |       4.375098 ms |            2.109050 ms |         134.805513 ms |
|                Latency p75 |       5.431229 ms |            2.403126 ms |         205.643246 ms |
|                Latency p90 |       6.422301 ms |            2.662319 ms |         227.073525 ms |
|                Latency p95 |       7.127944 ms |            3.112525 ms |         250.450359 ms |
|                Latency p99 |       8.805373 ms |          236.442165 ms |       10010.579309 ms |
|              Latency p99.9 |      17.743337 ms |         1006.184186 ms |       22707.315056 ms |
+----------------------------+-------------------+------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-NON-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/02-write-1M-keys-rate-limited/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value value (clients 1 to 1000)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 30 GB Memory + 150 GB SSD (1 for client)
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.4)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.2 (Go 1.7.4)


```
+----------------------------+-------------------+------------------------+-----------------------+
|                            | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.2-go1.7.4 |
+----------------------------+-------------------+------------------------+-----------------------+
|      READS-COMPLETED-DELTA |               186 |                    294 |                   383 |
|    SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| WRITES-COMPLETED-DELTA-SUM |           2091099 |                1668912 |               3973507 |
|  SECTORS-WRITTEN-DELTA-SUM |           1175004 |               14572148 |              32486700 |
|          RECEIVE-BYTES-SUM |            5.2 GB |                 5.4 GB |                8.9 GB |
|         TRANSMIT-BYTES-SUM |            4.1 GB |                 4.3 GB |                7.7 GB |
|              MAX-CPU-USAGE |           56.53 % |                56.39 % |               56.33 % |
|           MAX-MEMORY-USAGE |        1472.87 MB |             3652.64 MB |            4561.35 MB |
|              TOTAL-SECONDS |      570.7510 sec |           477.1301 sec |         1056.0748 sec |
|             AVG-THROUGHPUT | 1752.0775 req/sec |      2090.3063 req/sec |      946.9027 req/sec |
|            SLOWEST-LATENCY |       556.0183 ms |           3654.9846 ms |         17465.0180 ms |
|            FASTEST-LATENCY |         1.1633 ms |              0.9644 ms |             3.2258 ms |
|                AVG-LATENCY |        12.1934 ms |             30.7887 ms |            96.5730 ms |
|                Latency p10 |       2.312432 ms |            2.226581 ms |           4.048553 ms |
|                Latency p25 |       4.204692 ms |            2.847540 ms |           5.682306 ms |
|                Latency p50 |       7.689732 ms |            4.229944 ms |          10.091410 ms |
|                Latency p75 |      14.481892 ms |           10.914450 ms |          44.163918 ms |
|                Latency p90 |      24.692357 ms |           17.725627 ms |          94.140188 ms |
|                Latency p95 |      40.230709 ms |           24.428028 ms |         208.474321 ms |
|                Latency p99 |      60.787420 ms |         1188.220202 ms |         895.758025 ms |
|              Latency p99.9 |     527.752867 ms |         2326.553678 ms |       13622.807237 ms |
+----------------------------+-------------------+------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/03-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 1M keys, 1000-client, 256-byte key, 1KB value

- Packet Bare Metal
- 4 machines of 16 Physical Cores @ 2.6 GHz + 128 GB DDR4 RAM + 120 GB SSD + 20Gbps Bonded Network (1 for client)
- Ubuntu 16.04
- etcd v3.1 (Go 1.7.4)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.2 (Go 1.7.4)


```
+----------------------------+--------------------+------------------------+-----------------------+
|                            | etcd-v3.1-go1.7.4  | zookeeper-r3.4.9-java8 | consul-v0.7.2-go1.7.4 |
+----------------------------+--------------------+------------------------+-----------------------+
|      READS-COMPLETED-DELTA |                  6 |                      4 |                    13 |
|    SECTORS-READS-DELTA-SUM |                  0 |                      0 |                     0 |
| WRITES-COMPLETED-DELTA-SUM |             332647 |                 646300 |               1115043 |
|  SECTORS-WRITTEN-DELTA-SUM |                  0 |                      0 |                     0 |
|          RECEIVE-BYTES-SUM |             5.1 GB |                 5.3 GB |                5.6 GB |
|         TRANSMIT-BYTES-SUM |             3.9 GB |                 4.2 GB |                4.4 GB |
|              MAX-CPU-USAGE |           362.72 % |               388.31 % |              282.29 % |
|           MAX-MEMORY-USAGE |         1213.43 MB |             9291.73 MB |            4523.80 MB |
|              TOTAL-SECONDS |        40.8774 sec |            30.3340 sec |          164.0401 sec |
|             AVG-THROUGHPUT | 24463.3803 req/sec |     32964.2605 req/sec |     6096.0706 req/sec |
|            SLOWEST-LATENCY |        150.4056 ms |           1565.8847 ms |          7306.8885 ms |
|            FASTEST-LATENCY |          0.8704 ms |              0.3387 ms |             8.9377 ms |
|                AVG-LATENCY |         40.7504 ms |             29.4415 ms |           163.8167 ms |
|                Latency p10 |       19.337223 ms |           18.111482 ms |          58.038157 ms |
|                Latency p25 |       24.134461 ms |           20.482972 ms |          72.112776 ms |
|                Latency p50 |       33.954506 ms |           23.463435 ms |         106.720600 ms |
|                Latency p75 |       54.301812 ms |           27.142309 ms |         213.816984 ms |
|                Latency p90 |       72.115813 ms |           31.443281 ms |         333.154609 ms |
|                Latency p95 |       81.127899 ms |           35.534083 ms |         373.636304 ms |
|                Latency p99 |       98.866146 ms |          147.958507 ms |         461.374153 ms |
|              Latency p99.9 |      132.016073 ms |         1378.812699 ms |        5311.162416 ms |
+----------------------------+--------------------+------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-NON-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/04-write-1M-keys/AVG-TRANSMIT-BYTES-NUM-DELTA">




<br><br><hr>
##### Write 1M keys, 1000-client, 1000QPS, 256-byte key, 1KB value

- Packet Bare Metal
- 4 machines of 16 Physical Cores @ 2.6 GHz + 128 GB DDR4 RAM + 120 GB SSD + 20Gbps Bonded Network (1 for client)
- Ubuntu 16.04
- etcd v3.1 (Go 1.7.4)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.2 (Go 1.7.4)


```
+----------------------------+-------------------+------------------------+-----------------------+
|                            | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.2-go1.7.4 |
+----------------------------+-------------------+------------------------+-----------------------+
|      READS-COMPLETED-DELTA |               473 |                    706 |                   571 |
|    SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| WRITES-COMPLETED-DELTA-SUM |           8257619 |                9913357 |              13254939 |
|  SECTORS-WRITTEN-DELTA-SUM |                 0 |                      0 |                     0 |
|          RECEIVE-BYTES-SUM |            6.5 GB |                 5.8 GB |                6.4 GB |
|         TRANSMIT-BYTES-SUM |            5.2 GB |                 4.7 GB |                5.1 GB |
|              MAX-CPU-USAGE |           91.10 % |                64.88 % |              147.91 % |
|           MAX-MEMORY-USAGE |        1694.38 MB |             2842.88 MB |            4231.91 MB |
|              TOTAL-SECONDS |      999.0582 sec |           999.8117 sec |         1040.1450 sec |
|             AVG-THROUGHPUT | 1000.9427 req/sec |       999.5752 req/sec |      961.4044 req/sec |
|            SLOWEST-LATENCY |        56.0360 ms |           2635.4569 ms |          8546.8756 ms |
|            FASTEST-LATENCY |         0.7225 ms |              0.3447 ms |             1.3011 ms |
|                AVG-LATENCY |         1.3866 ms |              6.9357 ms |            93.2921 ms |
|                Latency p10 |       1.019233 ms |            0.667793 ms |           2.297153 ms |
|                Latency p25 |       1.140797 ms |            0.757403 ms |           2.458065 ms |
|                Latency p50 |       1.303873 ms |            0.876627 ms |           2.654298 ms |
|                Latency p75 |       1.429593 ms |            0.982231 ms |           3.600769 ms |
|                Latency p90 |       1.584000 ms |            1.121141 ms |          10.765193 ms |
|                Latency p95 |       1.733878 ms |            1.236501 ms |          65.069180 ms |
|                Latency p99 |       3.298035 ms |           28.890276 ms |        4628.091116 ms |
|              Latency p99.9 |      17.796168 ms |         1369.914634 ms |        7243.951610 ms |
+----------------------------+-------------------+------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-NON-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/05-write-1M-keys-rate-limited/AVG-TRANSMIT-BYTES-NUM-DELTA">





<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value value (clients 1 to 1000)

- Packet Bare Metal
- 4 machines of 16 Physical Cores @ 2.6 GHz + 128 GB DDR4 RAM + 120 GB SSD + 20Gbps Bonded Network (1 for client)
- Ubuntu 16.04
- etcd v3.1 (Go 1.7.4)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.2 (Go 1.7.4)


```
+----------------------------+-------------------+------------------------+-----------------------+
|                            | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.2-go1.7.4 |
+----------------------------+-------------------+------------------------+-----------------------+
|      READS-COMPLETED-DELTA |               308 |                    217 |                   541 |
|    SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| WRITES-COMPLETED-DELTA-SUM |           3480190 |                3479238 |               6605626 |
|  SECTORS-WRITTEN-DELTA-SUM |                 0 |                      0 |                     0 |
|          RECEIVE-BYTES-SUM |            5.4 GB |                 5.5 GB |                5.8 GB |
|         TRANSMIT-BYTES-SUM |            4.2 GB |                 4.4 GB |                4.5 GB |
|              MAX-CPU-USAGE |          135.68 % |               122.75 % |              124.26 % |
|           MAX-MEMORY-USAGE |        1333.86 MB |             6610.87 MB |            4790.88 MB |
|              TOTAL-SECONDS |      276.4103 sec |           185.1561 sec |          582.4233 sec |
|             AVG-THROUGHPUT | 3617.8096 req/sec |      5400.8493 req/sec |     1716.9642 req/sec |
|            SLOWEST-LATENCY |       155.0781 ms |           2140.8240 ms |          5637.2541 ms |
|            FASTEST-LATENCY |         0.6443 ms |              0.2915 ms |             1.3058 ms |
|                AVG-LATENCY |        12.7328 ms |              9.5968 ms |            47.3448 ms |
|                Latency p10 |       1.115963 ms |            0.706270 ms |           2.296140 ms |
|                Latency p25 |       1.484697 ms |            0.890714 ms |           3.300920 ms |
|                Latency p50 |       2.696072 ms |            1.594529 ms |           6.696194 ms |
|                Latency p75 |      14.830030 ms |           10.053333 ms |          42.282100 ms |
|                Latency p90 |      36.917286 ms |           18.814867 ms |          93.939972 ms |
|                Latency p95 |      67.270709 ms |           24.070593 ms |         131.658011 ms |
|                Latency p99 |      92.520682 ms |           38.044010 ms |         473.239177 ms |
|              Latency p99.9 |     122.274461 ms |         1725.413184 ms |        4949.483214 ms |
+----------------------------+-------------------+------------------------+-----------------------+
```


<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-LATENCY-MS.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-LATENCY-MS">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-THROUGHPUT.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-THROUGHPUT">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-NON-VOLUNTARY-CTXT-SWITCHES">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-CPU.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-CPU">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-VMRSS-MB.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-VMRSS-MB">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-READS-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-SECTORS-READ-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-WRITES-COMPLETED-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-SECTORS-WRITTEN-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-RECEIVE-BYTES-NUM-DELTA">

<img src="https://storage.googleapis.com/dbtester-results/2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA.svg" alt="2017Q1-01-etcd-zookeeper-consul/06-write-1M-keys-client-variable/AVG-TRANSMIT-BYTES-NUM-DELTA">
