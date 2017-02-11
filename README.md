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

All logs and results can be found at https://github.com/coreos/dbtester/tree/master/test-results or https://console.cloud.google.com/storage/browser/dbtester-results/?authuser=0&project=etcd-development.



<br><br><hr>
##### Noticeable Warnings: Zookeeper

Snapshot, when writing 1-million entries (256-byte key, 1KB value value), with 500 concurrent clients

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

When writing more than 2-million entries (256-byte key, 1KB value value) with 500 concurrent clients

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

Logs do not tell much but average latency spikes (e.g. from 70.27517 ms to 10407.900082 ms)



<br><br><hr>
##### Write 1M keys, 256-byte key, 1KB value value, clients 1 to 1,000

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 30 GB Memory + 150 GB SSD (1 for client)
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.5)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.4 (Go 1.7.5)


```
+---------------------------------------+-------------------+------------------------+-----------------------+
|                                       | etcd-v3.1-go1.7.5 | zookeeper-r3.4.9-java8 | consul-v0.7.4-go1.7.5 |
+---------------------------------------+-------------------+------------------------+-----------------------+
|                         TOTAL-SECONDS |      342.2984 sec |           320.9968 sec |          888.9235 sec |
|                  TOTAL-REQUEST-NUMBER |         1,000,000 |              1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |    34,747 req/sec |         43,558 req/sec |        16,486 req/sec |
|                        AVG-THROUGHPUT |     2,921 req/sec |          3,115 req/sec |         1,124 req/sec |
|                        MIN-THROUGHPUT |        29 req/sec |              0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |         1.1239 ms |              1.1194 ms |             3.1083 ms |
|                           AVG-LATENCY |        13.9400 ms |             36.2980 ms |           152.1034 ms |
|                       SLOWEST-LATENCY |       115.2305 ms |           2626.1766 ms |         20615.4531 ms |
|                           Latency p10 |       2.309089 ms |            2.512687 ms |           3.992750 ms |
|                           Latency p25 |       5.936953 ms |            3.853528 ms |           8.522133 ms |
|                           Latency p50 |      10.201295 ms |            6.619094 ms |          36.083945 ms |
|                           Latency p75 |      16.895621 ms |           13.507464 ms |         107.288320 ms |
|                           Latency p90 |      28.576840 ms |           20.348043 ms |         313.922534 ms |
|                           Latency p95 |      44.602367 ms |           28.054291 ms |         478.639678 ms |
|                           Latency p99 |      60.177421 ms |         1108.928408 ms |         623.624914 ms |
|                         Latency p99.9 |      92.142990 ms |         2619.595630 ms |       14534.009968 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |            5.0 GB |                 5.4 GB |                6.7 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |            3.9 GB |                 4.3 GB |                5.4 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |            270 MB |                 356 MB |                210 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |            1.5 GB |                 1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |          431.33 % |               601.67 % |              466.00 % |
|               SERVER-MAX-MEMORY-USAGE |            1.3 GB |                 3.9 GB |                5.0 GB |
|                  CLIENT-MAX-CPU-USAGE |          559.00 % |               685.00 % |              231.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            327 MB |                 4.4 GB |                201 MB |
|                    CLIENT-ERROR-COUNT |                 0 |                     13 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                72 |                    389 |                   255 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |         1,525,739 |              1,199,135 |             4,434,402 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |        32,134,464 |             40,860,544 |           183,627,904 |
|           SERVER-AVG-DISK-SPACE-USAGE |            3.1 GB |                 7.6 GB |                2.8 GB |
+---------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: could not connect to a server" (count 13)
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
##### Write 1M keys, 256-byte key, 1KB value, Best Throughput (etcd 1,000, Zookeeper 500, Consul 500 clients)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 30 GB Memory + 150 GB SSD (1 for client)
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.5)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.4 (Go 1.7.5)


```
+---------------------------------------+-------------------+------------------------+-----------------------+
|                                       | etcd-v3.1-go1.7.5 | zookeeper-r3.4.9-java8 | consul-v0.7.4-go1.7.5 |
+---------------------------------------+-------------------+------------------------+-----------------------+
|                         TOTAL-SECONDS |       36.5583 sec |            63.6203 sec |          261.0246 sec |
|                  TOTAL-REQUEST-NUMBER |         1,000,000 |              1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |    35,187 req/sec |         44,883 req/sec |        15,009 req/sec |
|                        AVG-THROUGHPUT |    27,353 req/sec |         15,684 req/sec |         3,831 req/sec |
|                        MIN-THROUGHPUT |    13,891 req/sec |              0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |         4.7388 ms |              2.4777 ms |            14.0583 ms |
|                           AVG-LATENCY |        36.3305 ms |             22.5664 ms |           130.4658 ms |
|                       SLOWEST-LATENCY |       346.3847 ms |           3517.8313 ms |         19316.7564 ms |
|                           Latency p10 |      15.483941 ms |            7.766447 ms |          32.229589 ms |
|                           Latency p25 |      18.596901 ms |            9.156112 ms |          37.179339 ms |
|                           Latency p50 |      23.920164 ms |           10.532408 ms |          47.314383 ms |
|                           Latency p75 |      44.685986 ms |           12.228320 ms |          70.740623 ms |
|                           Latency p90 |      64.373404 ms |           14.214414 ms |         104.953131 ms |
|                           Latency p95 |     113.441501 ms |           16.216272 ms |         282.116427 ms |
|                           Latency p99 |     167.148590 ms |          340.373624 ms |        1092.781013 ms |
|                         Latency p99.9 |     194.350987 ms |         2151.870409 ms |       11587.696623 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |            5.0 GB |                 5.4 GB |                 10 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |            3.8 GB |                 4.3 GB |                9.1 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |            282 MB |                 352 MB |                218 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |            1.4 GB |                 1.4 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |          446.67 % |               731.27 % |              379.33 % |
|               SERVER-MAX-MEMORY-USAGE |            1.2 GB |                 4.6 GB |                5.8 GB |
|                  CLIENT-MAX-CPU-USAGE |          568.00 % |               404.00 % |              223.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            248 MB |                 3.1 GB |                 81 MB |
|                    CLIENT-ERROR-COUNT |                 0 |                  2,150 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                 2 |                    213 |                   147 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           102,486 |                 86,871 |               641,255 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |        20,504,912 |             26,252,736 |           110,850,712 |
|           SERVER-AVG-DISK-SPACE-USAGE |            2.7 GB |                 6.9 GB |                3.0 GB |
+---------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: could not connect to a server" (count 765)
"zk: connection closed" (count 1,385)
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
- etcd v3.1 (Go 1.7.5)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.4 (Go 1.7.5)


```
+---------------------------------------+-------------------+------------------------+-----------------------+
|                                       | etcd-v3.1-go1.7.5 | zookeeper-r3.4.9-java8 | consul-v0.7.4-go1.7.5 |
+---------------------------------------+-------------------+------------------------+-----------------------+
|                         TOTAL-SECONDS |       36.5091 sec |            59.0429 sec |          116.4349 sec |
|                  TOTAL-REQUEST-NUMBER |         1,000,000 |              1,000,000 |             1,000,000 |
|                        MAX-THROUGHPUT |    35,162 req/sec |         50,196 req/sec |        14,615 req/sec |
|                        AVG-THROUGHPUT |    27,390 req/sec |         16,854 req/sec |         8,588 req/sec |
|                        MIN-THROUGHPUT |    14,654 req/sec |              0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |         4.9384 ms |              1.4337 ms |            15.3309 ms |
|                           AVG-LATENCY |        36.3098 ms |             44.2298 ms |           115.8624 ms |
|                       SLOWEST-LATENCY |       353.6275 ms |           7585.4064 ms |         12902.7031 ms |
|                           Latency p10 |      14.627933 ms |           12.299326 ms |          67.088002 ms |
|                           Latency p25 |      17.957404 ms |           15.830300 ms |          73.441500 ms |
|                           Latency p50 |      23.290321 ms |           18.660221 ms |          82.581044 ms |
|                           Latency p75 |      43.118899 ms |           21.623243 ms |          96.395872 ms |
|                           Latency p90 |      64.488242 ms |           26.143493 ms |         166.706106 ms |
|                           Latency p95 |     126.257124 ms |           77.780210 ms |         218.288165 ms |
|                           Latency p99 |     173.468949 ms |          843.928907 ms |         323.392356 ms |
|                         Latency p99.9 |     194.822900 ms |         2909.696060 ms |        7179.211229 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |            5.0 GB |                 5.3 GB |                5.6 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |            3.8 GB |                 4.3 GB |                4.4 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |            282 MB |                 371 MB |                243 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |            1.4 GB |                 1.5 GB |                1.5 GB |
|                  SERVER-MAX-CPU-USAGE |          433.67 % |               634.40 % |              478.67 % |
|               SERVER-MAX-MEMORY-USAGE |            1.1 GB |                 4.9 GB |                5.1 GB |
|                  CLIENT-MAX-CPU-USAGE |          582.00 % |               422.00 % |              219.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            266 MB |                 4.9 GB |                115 MB |
|                    CLIENT-ERROR-COUNT |                 0 |                  4,886 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |                70 |                    217 |                    46 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           102,259 |                 80,718 |               377,108 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |        20,520,584 |             27,706,344 |            48,882,704 |
|           SERVER-AVG-DISK-SPACE-USAGE |            2.7 GB |                 5.9 GB |                3.1 GB |
+---------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: connection closed" (count 2,381)
"zk: could not connect to a server" (count 2,505)
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
##### Write 3-million keys, 256-byte key, 1KB value, Best Throughput (etcd 1,000, Zookeeper 500, Consul 500 clients)

- Google Cloud Compute Engine
- 4 machines of 16 vCPUs + 30 GB Memory + 150 GB SSD (1 for client)
- Ubuntu 16.10
- etcd v3.1 (Go 1.7.5)
- Zookeeper r3.4.9
  - Java 8
  - javac 1.8.0_121
  - Java(TM) SE Runtime Environment (build 1.8.0_121-b13)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.121-b13, mixed mode)
- Consul v0.7.4 (Go 1.7.5)


```
+---------------------------------------+-------------------+------------------------+-----------------------+
|                                       | etcd-v3.1-go1.7.5 | zookeeper-r3.4.9-java8 | consul-v0.7.4-go1.7.5 |
+---------------------------------------+-------------------+------------------------+-----------------------+
|                         TOTAL-SECONDS |      138.0049 sec |          2242.5513 sec |         2361.4227 sec |
|                  TOTAL-REQUEST-NUMBER |         3,000,000 |              3,000,000 |             3,000,000 |
|                        MAX-THROUGHPUT |    33,963 req/sec |         42,804 req/sec |        16,033 req/sec |
|                        AVG-THROUGHPUT |    21,738 req/sec |            906 req/sec |         1,270 req/sec |
|                        MIN-THROUGHPUT |     7,006 req/sec |              0 req/sec |             0 req/sec |
|                       FASTEST-LATENCY |         3.2812 ms |              1.2984 ms |            16.9263 ms |
|                           AVG-LATENCY |        45.9259 ms |             42.7162 ms |           393.5551 ms |
|                       SLOWEST-LATENCY |       259.3153 ms |           6921.5532 ms |         30425.8070 ms |
|                           Latency p10 |      15.658060 ms |            8.200137 ms |          34.534783 ms |
|                           Latency p25 |      18.852826 ms |            9.417761 ms |          42.917753 ms |
|                           Latency p50 |      23.711419 ms |           10.902322 ms |          64.087200 ms |
|                           Latency p75 |      53.911600 ms |           12.742504 ms |          91.427820 ms |
|                           Latency p90 |     130.485283 ms |           15.178863 ms |         167.938406 ms |
|                           Latency p95 |     151.376530 ms |           38.140465 ms |         951.339409 ms |
|                           Latency p99 |     171.722975 ms |         1540.586327 ms |       10968.875166 ms |
|                         Latency p99.9 |     188.102916 ms |         2276.156658 ms |       18546.023298 ms |
|      SERVER-TOTAL-NETWORK-RX-DATA-SUM |             15 GB |                  12 GB |                 89 GB |
|      SERVER-TOTAL-NETWORK-TX-DATA-SUM |             12 GB |                  10 GB |                 85 GB |
|           CLIENT-TOTAL-NETWORK-RX-SUM |            852 MB |                 997 MB |                685 MB |
|           CLIENT-TOTAL-NETWORK-TX-SUM |            4.3 GB |                 3.4 GB |                4.4 GB |
|                  SERVER-MAX-CPU-USAGE |          442.00 % |               785.67 % |              439.33 % |
|               SERVER-MAX-MEMORY-USAGE |            2.8 GB |                 7.4 GB |                 15 GB |
|                  CLIENT-MAX-CPU-USAGE |          554.00 % |               368.00 % |              222.00 % |
|               CLIENT-MAX-MEMORY-USAGE |            348 MB |                 3.1 GB |                157 MB |
|                    CLIENT-ERROR-COUNT |                 0 |                967,222 |                     0 |
|  SERVER-AVG-READS-COMPLETED-DELTA-SUM |               168 |                    443 |                23,290 |
|    SERVER-AVG-SECTORS-READS-DELTA-SUM |                 0 |                      0 |                     0 |
| SERVER-AVG-WRITES-COMPLETED-DELTA-SUM |           303,974 |                231,817 |             2,976,074 |
|  SERVER-AVG-SECTORS-WRITTEN-DELTA-SUM |        62,193,472 |             84,808,712 |           990,569,940 |
|           SERVER-AVG-DISK-SPACE-USAGE |            6.5 GB |                  21 GB |                8.3 GB |
+---------------------------------------+-------------------+------------------------+-----------------------+


zookeeper errors:
"zk: could not connect to a server" (count 963,503)
"zk: connection closed" (count 3,719)
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



