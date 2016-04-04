

- Google Cloud Compute Engine
- 3 machines(server) of 8 vCPUs + 16GB Memory + 375GB local SSD (SCSI)
- 1 machine(client) of 16 vCPUs + 30GB Memory + 50GB SSD
- Ubuntu 15.10
- Go 1.6 with etcd master branch as of testing date
- Java 8 with Zookeeper 3.4.8(current)
  - Java(TM) SE Runtime Environment (build 1.8.0_74-b02)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.74-b02, mixed mode)



<br><br><hr>
##### Write 300K keys, 1 client, key 64 bytes, value 256 bytes

![bench-01-avg-latency-ms.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-01-avg-latency-ms.svg)

![bench-01-throughput.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-01-throughput.svg)

![bench-01-avg-cpu.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-01-avg-cpu.svg)

![bench-01-avg-memory.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-01-avg-memory.svg)



<br><br><hr>
##### Write 1M keys, 10 clients, key 64 bytes, value 256 bytes

![bench-02-avg-latency-ms.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-02-avg-latency-ms.svg)

![bench-02-throughput.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-02-throughput.svg)

![bench-02-avg-cpu.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-02-avg-cpu.svg)

![bench-02-avg-memory.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-02-avg-memory.svg)



<br><br><hr>
##### Write 3M keys, 500 clients, key 64 bytes, value 256 bytes

![bench-03-avg-latency-ms.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-03-avg-latency-ms.svg)

![bench-03-throughput.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-03-throughput.svg)

![bench-03-avg-cpu.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-03-avg-cpu.svg)

![bench-03-avg-memory.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-03-avg-memory.svg)



<br><br><hr>
##### Write 3M keys, 1K clients, key 64 bytes, value 256 bytes

![bench-04-avg-latency-ms.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-04-avg-latency-ms.svg)

![bench-04-throughput.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-04-throughput.svg)

![bench-04-avg-cpu.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-04-avg-cpu.svg)

![bench-04-avg-memory.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-04-avg-memory.svg)



<br><br><hr>
##### Read single key 1M times, 1 client, key 64 bytes, value 1 kb

![bench-05-avg-latency-ms.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-05-avg-latency-ms.svg)

![bench-05-throughput.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-05-throughput.svg)

![bench-05-avg-cpu.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-05-avg-cpu.svg)

![bench-05-avg-memory.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-05-avg-memory.svg)



<br><br><hr>
##### Read single key 1M times, 100 clients, key 64 bytes, value 1 kb

![bench-06-avg-latency-ms.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-06-avg-latency-ms.svg)

![bench-06-throughput.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-06-throughput.svg)

![bench-06-avg-cpu.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-06-avg-cpu.svg)

![bench-06-avg-memory.svg](https://cdn.rawgit.com/coreos/dbtester/master/bench-results/bench-06-avg-memory.svg)



