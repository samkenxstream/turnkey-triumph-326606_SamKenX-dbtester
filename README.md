# dbtester

[![Build Status](https://img.shields.io/travis/coreos/dbtester.svg?style=flat-square)](https://travis-ci.org/coreos/dbtester) [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/coreos/dbtester)

Distributed database benchmark tester: etcd, Zookeeper, Consul, zetcd, cetcd

- Database Agent
  - https://github.com/coreos/dbtester/tree/master/agent
- Database Client
  - https://github.com/coreos/dbtester/tree/master/control
- System Metrics
  - https://github.com/gyuho/psn
- Test Data Analysis
  - https://github.com/coreos/dbtester/tree/master/analyze
  - https://github.com/gyuho/dataframe
  - https://github.com/gonum/plot

For etcd, we recommend [etcd benchmark tool](https://github.com/coreos/etcd/tree/master/tools/benchmark).

All logs and results can be found at https://github.com/coreos/dbtester/tree/master/test-results



<br><br><hr>
##### Noticeable Warnings: Zookeeper

Snapshot, when writing 1-million entries (256-byte key, 1KB value value), with 300 concurrent clients

```
# snapshot warnings
cd 2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput
grep -r -i fsync-ing\ the zookeeper-r3.4.9-java8-* | less

2017-02-10 18:55:38,997 [myid:3] - WARN  [SyncThread:3:SyncRequestProcessor@148] - Too busy to snap, skipping
2017-02-10 18:55:38,998 [myid:3] - INFO  [SyncThread:3:FileTxnLog@203] - Creating new log file: log.1000c0c51
2017-02-10 18:55:40,855 [myid:3] - INFO  [SyncThread:3:FileTxnLog@203] - Creating new log file: log.1000cd2e6
2017-02-10 18:55:40,855 [myid:3] - INFO  [Snapshot Thread:FileTxnSnapLog@240] - Snapshotting: 0x1000cd1ca to /home/gyuho/zookeeper/zookeeper.data/version-2/snapshot.1000cd1ca
2017-02-10 18:55:46,382 [myid:3] - WARN  [SyncThread:3:FileTxnLog@338] - fsync-ing the write ahead log in SyncThread:3 took 1062ms which will adversely effect operation latency. See the ZooKeeper troubleshooting guide
2017-02-10 18:55:47,471 [myid:3] - WARN  [SyncThread:3:FileTxnLog@338] - fsync-ing the write ahead log in SyncThread:3 took 1084ms which will adversely effect operation latency. See the ZooKeeper troubleshooting guide
2017-02-10 18:55:49,425 [myid:3] - WARN  [SyncThread:3:FileTxnLog@338] - fsync-ing the write ahead log in SyncThread:3 took 1142ms which will adversely effect operation latency. See the ZooKeeper troubleshooting guide
2017-02-10 18:55:51,188 [myid:3] - WARN  [SyncThread:3:FileTxnLog@338] - fsync-ing the write ahead log in SyncThread:3 took 1201ms which will adversely effect operation latency. See the ZooKeeper troubleshooting guide
2017-02-10 18:55:52,292 [myid:3] - WARN  [SyncThread:3:FileTxnLog@338] - fsync-ing the write ahead log in SyncThread:3 took 1102ms which will adversely effect operation latency. See the ZooKeeper troubleshooting guide
```

When writing more than 2-million entries (256-byte key, 1KB value value) with 300 concurrent clients

```
# leader election
cd 2017Q1-00-etcd-zookeeper-consul/04-write-too-many-keys
grep -r -i election\ took  zookeeper-r3.4.9-java8-* | less

# leader election is taking more than 10 seconds...
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:22:16,549 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Follower@61] - FOLLOWING - LEADER ELECTION TOOK - 22978
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:23:02,279 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 10210
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:23:14,498 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 203
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:23:36,303 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 9791
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:23:52,151 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 3836
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:24:13,849 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 9686
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:24:29,694 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 3573
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:24:51,392 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 8686
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:25:07,231 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 3827
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:25:28,940 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 9697
zookeeper-r3.4.9-java8-2-database.log:2017-02-10 19:25:44,772 [myid:2] - INFO  [QuorumPeer[myid=2]/0:0:0:0:0:0:0:0:2181:Leader@361] - LEADING - LEADER ELECTION TOOK - 3820
```


<br><br><hr>
##### Noticeable Warnings: Consul

Snapshot, when writing 1-million entries (256-byte key, 1KB value value), with 500 concurrent clients

```
# snapshot warnings
cd 2017Q1-00-etcd-zookeeper-consul/02-write-1M-keys-best-throughput
grep -r -i installed\ remote consul-v0.7.4-go1.7.5-* | less

    2017/02/10 18:58:43 [INFO] snapshot: Creating new snapshot at /home/gyuho/consul.data/raft/snapshots/2-900345-1486753123478.tmp
    2017/02/10 18:58:45 [INFO] snapshot: reaping snapshot /home/gyuho/consul.data/raft/snapshots/2-849399-1486753096972
    2017/02/10 18:58:46 [INFO] raft: Copied 1223270573 bytes to local snapshot
    2017/02/10 18:58:55 [INFO] raft: Compacting logs from 868354 to 868801
    2017/02/10 18:58:56 [INFO] raft: Installed remote snapshot
    2017/02/10 18:58:57 [INFO] snapshot: Creating new snapshot at /home/gyuho/consul.data/raft/snapshots/2-911546-1486753137827.tmp
    2017/02/10 18:58:59 [INFO] consul.fsm: snapshot created in 32.255µs
    2017/02/10 18:59:01 [INFO] snapshot: reaping snapshot /home/gyuho/consul.data/raft/snapshots/2-873921-1486753116619
    2017/02/10 18:59:02 [INFO] raft: Copied 1238491373 bytes to local snapshot
    2017/02/10 18:59:11 [INFO] raft: Compacting logs from 868802 to 868801
    2017/02/10 18:59:11 [INFO] raft: Installed remote snapshot
```

Logs do not tell much but average latency spikes (e.g. from 258.28656 ms to 6265.185836 ms)



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
+---------------------------------------+-------------------+------------------------+-----------------------+
|                                       | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.3-go1.7.4 |
+---------------------------------------+-------------------+------------------------+-----------------------+
|                         TOTAL-SECONDS |      338.7661 sec |           344.3563 sec |          703.7060 sec |
|                  TOTAL-REQUEST-NUMBER |         1,000,000 |              1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |    35,147 req/sec |         31,726 req/sec |        15,328 req/sec |
|                        AVG-THROUGHPUT |     2,951 req/sec |          2,903 req/sec |         1,421 req/sec |
|                        MIN-THROUGHPUT |        90 req/sec |              0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |         1.1001 ms |              1.1093 ms |             2.9964 ms |
|                           AVG-LATENCY |        13.8862 ms |             34.9948 ms |            72.5791 ms |
|                       SLOWEST-LATENCY |       109.4800 ms |           2618.2703 ms |         20860.6692 ms |
|                           Latency p10 |       2.295037 ms |            2.620473 ms |           3.982040 ms |
|                           Latency p25 |       5.788546 ms |            3.932461 ms |           7.888984 ms |
|                           Latency p50 |       9.935599 ms |            7.747493 ms |          21.950488 ms |
|                           Latency p75 |      17.040088 ms |           16.891088 ms |          58.936521 ms |
|                           Latency p90 |      28.513968 ms |           28.114578 ms |         126.568085 ms |
|                           Latency p95 |      44.023164 ms |           44.599685 ms |         165.331967 ms |
|                           Latency p99 |      60.351324 ms |         1063.554863 ms |         430.071868 ms |
|                         Latency p99.9 |      94.058105 ms |         2336.144865 ms |       12648.734251 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |            5.1 GB |                 5.4 GB |                7.9 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |            3.9 GB |                 4.4 GB |                6.6 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |            270 MB |                 357 MB |                202 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |            1.5 GB |                 1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |          434.00 % |               600.67 % |              416.00 % |
|               SERVER-MAX-MEMORY-USAGE |            1.3 GB |                 4.0 GB |                5.4 GB |
|                  CLIENT-MAX-CPU-USAGE |          540.00 % |               322.00 % |              204.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            330 MB |                 3.6 GB |                199 MB |
|                    CLIENT-ERROR-COUNT |                 0 |                     24 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                76 |                    334 |                    66 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |         1,217,294 |                953,784 |             2,381,092 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |           714,152 |              9,304,072 |            10,711,132 |
|           SERVER-AVG-DISK-SPACE-USAGE |            3.0 GB |                 7.9 GB |                3.0 GB |
+---------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: could not connect to a server" (count 24)
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
+---------------------------------------+-------------------+------------------------+-----------------------+
|                                       | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.3-go1.7.4 |
+---------------------------------------+-------------------+------------------------+-----------------------+
|                         TOTAL-SECONDS |       37.3284 sec |            75.0700 sec |          304.4858 sec |
|                  TOTAL-REQUEST-NUMBER |         1,000,000 |              1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |    33,567 req/sec |         35,499 req/sec |        15,141 req/sec |
|                        AVG-THROUGHPUT |    26,789 req/sec |         13,274 req/sec |         3,284 req/sec |
|                        MIN-THROUGHPUT |    10,018 req/sec |              0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |         4.2842 ms |              2.7405 ms |            11.4297 ms |
|                           AVG-LATENCY |        26.0603 ms |             18.2231 ms |           152.1359 ms |
|                       SLOWEST-LATENCY |       520.8716 ms |           4264.2996 ms |         28029.3953 ms |
|                           Latency p10 |      10.171289 ms |            6.401553 ms |          30.579107 ms |
|                           Latency p25 |      12.254908 ms |            7.300705 ms |          35.763003 ms |
|                           Latency p50 |      17.138243 ms |            8.302805 ms |          48.467608 ms |
|                           Latency p75 |      23.925669 ms |            9.453586 ms |          80.519456 ms |
|                           Latency p90 |      48.690057 ms |           10.764813 ms |         248.959013 ms |
|                           Latency p95 |      76.533161 ms |           11.992104 ms |         349.281928 ms |
|                           Latency p99 |     146.318242 ms |          153.580393 ms |        1324.508306 ms |
|                         Latency p99.9 |     183.924901 ms |         1935.929712 ms |       10622.316021 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |            5.0 GB |                 6.2 GB |                 11 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |            3.8 GB |                 5.1 GB |                 10 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |            274 MB |                 350 MB |                216 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |            1.4 GB |                 1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |          407.67 % |               704.97 % |              380.00 % |
|               SERVER-MAX-MEMORY-USAGE |            1.1 GB |                 5.1 GB |                6.2 GB |
|                  CLIENT-MAX-CPU-USAGE |          454.00 % |               292.00 % |              202.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            210 MB |                 1.7 GB |                 88 MB |
|                    CLIENT-ERROR-COUNT |                 0 |                  3,452 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                 5 |                    212 |                   270 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           112,190 |                109,945 |               681,774 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |           492,444 |             10,249,020 |            32,988,480 |
|           SERVER-AVG-DISK-SPACE-USAGE |            2.8 GB |                 7.3 GB |                2.9 GB |
+---------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: could not connect to a server" (count 3,152)
"zk: connection closed" (count 300)
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
+---------------------------------------+-------------------+------------------------+-----------------------+
|                                       | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.3-go1.7.4 |
+---------------------------------------+-------------------+------------------------+-----------------------+
|                         TOTAL-SECONDS |       36.3917 sec |            72.1698 sec |          140.9480 sec |
|                  TOTAL-REQUEST-NUMBER |         1,000,000 |              1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |    36,861 req/sec |         39,923 req/sec |        14,455 req/sec |
|                        AVG-THROUGHPUT |    27,478 req/sec |         13,704 req/sec |         7,094 req/sec |
|                        MIN-THROUGHPUT |     7,088 req/sec |              0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |         3.7509 ms |              4.3350 ms |            12.9159 ms |
|                           AVG-LATENCY |        36.2639 ms |             49.9165 ms |           140.4538 ms |
|                       SLOWEST-LATENCY |       244.3595 ms |           6056.0204 ms |         21808.2916 ms |
|                           Latency p10 |      13.700258 ms |           14.719617 ms |          65.494475 ms |
|                           Latency p25 |      16.855903 ms |           20.289440 ms |          71.570399 ms |
|                           Latency p50 |      21.895662 ms |           24.861967 ms |          80.420560 ms |
|                           Latency p75 |      43.387011 ms |           28.165478 ms |          96.334874 ms |
|                           Latency p90 |      69.199925 ms |           32.086272 ms |         161.048409 ms |
|                           Latency p95 |     133.035887 ms |           58.675701 ms |         210.603403 ms |
|                           Latency p99 |     173.021441 ms |          880.253639 ms |         484.701759 ms |
|                         Latency p99.9 |     188.144407 ms |         3260.921550 ms |       15472.281038 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |            5.0 GB |                 5.5 GB |                5.6 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |            3.9 GB |                 4.4 GB |                4.4 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |            282 MB |                 357 MB |                243 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |            1.4 GB |                 1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |          443.00 % |               695.00 % |              414.00 % |
|               SERVER-MAX-MEMORY-USAGE |            1.1 GB |                 4.8 GB |                4.9 GB |
|                  CLIENT-MAX-CPU-USAGE |          566.00 % |               324.00 % |              205.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            278 MB |                 5.3 GB |                114 MB |
|                    CLIENT-ERROR-COUNT |                 0 |                 10,929 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                 2 |                    216 |                    64 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           100,985 |                 91,154 |               289,238 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |           552,592 |              9,923,988 |             9,537,168 |
|           SERVER-AVG-DISK-SPACE-USAGE |            2.7 GB |                 6.7 GB |                3.1 GB |
+---------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: could not connect to a server" (count 8,818)
"zk: connection closed" (count 2,111)
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
+---------------------------------------+-------------------+------------------------+-----------------------+
|                                       | etcd-v3.1-go1.7.4 | zookeeper-r3.4.9-java8 | consul-v0.7.3-go1.7.4 |
+---------------------------------------+-------------------+------------------------+-----------------------+
|                         TOTAL-SECONDS |      141.1728 sec |          4134.1565 sec |         2979.6229 sec |
|                  TOTAL-REQUEST-NUMBER |         3,000,000 |              3,000,000 |             3,000,000 |
|                        MAX-THROUGHPUT |    32,215 req/sec |         38,259 req/sec |        15,784 req/sec |
|                        AVG-THROUGHPUT |    21,250 req/sec |            465 req/sec |         1,006 req/sec |
|                        MIN-THROUGHPUT |     9,081 req/sec |              0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |         3.6479 ms |              1.8912 ms |            10.3322 ms |
|                           AVG-LATENCY |        32.9036 ms |             28.5018 ms |           496.5874 ms |
|                       SLOWEST-LATENCY |      1030.3392 ms |           4729.8104 ms |         34484.1252 ms |
|                           Latency p10 |      10.882650 ms |            5.694308 ms |          35.045139 ms |
|                           Latency p25 |      13.005974 ms |            6.357600 ms |          44.711879 ms |
|                           Latency p50 |      18.743999 ms |            7.379354 ms |          66.739638 ms |
|                           Latency p75 |      31.054943 ms |            8.886502 ms |         116.150396 ms |
|                           Latency p90 |      96.027241 ms |           11.474394 ms |         540.305996 ms |
|                           Latency p95 |     123.844111 ms |           14.391699 ms |        1209.603980 ms |
|                           Latency p99 |     148.837786 ms |          722.847010 ms |       11957.052613 ms |
|                         Latency p99.9 |     162.994644 ms |         2457.522611 ms |       26068.610794 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |             15 GB |                  26 GB |                102 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |             12 GB |                  24 GB |                 98 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |            832 MB |                 978 MB |                688 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |            4.3 GB |                 3.3 GB |                4.4 GB |
|                  SERVER-MAX-CPU-USAGE |          421.00 % |               752.33 % |              410.67 % |
|               SERVER-MAX-MEMORY-USAGE |            2.8 GB |                 7.1 GB |                 17 GB |
|                  CLIENT-MAX-CPU-USAGE |          442.00 % |               376.00 % |              217.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            265 MB |                 1.9 GB |                174 MB |
|                    CLIENT-ERROR-COUNT |                 0 |              1,076,632 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                14 |                    309 |                94,864 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           333,196 |                360,314 |             3,689,221 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |         1,604,172 |             67,295,610 |           794,919,246 |
|           SERVER-AVG-DISK-SPACE-USAGE |            6.5 GB |                  27 GB |                8.2 GB |
+---------------------------------------+-------------------+------------------------+-----------------------+


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



