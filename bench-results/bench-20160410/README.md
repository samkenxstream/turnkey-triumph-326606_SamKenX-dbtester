

- Google Cloud Compute Engine
- 8 vCPUs + 16GB Memory + 375GB local SSD (SCSI)
- 1 machine(client) of 16 vCPUs + 30GB Memory + 50GB SSD
- Ubuntu 15.10
- Go 1.6 with etcd master branch as of testing date
- Java 8 with etcd v3 (patch(current)
  - Java(TM) SE Runtime Environment (build 1.8.0_74-b02)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.74-b02, mixed mode)



<br><br><hr>
##### Write 300K keys, 1K clients, key 64 bytes, value 256 bytes

<img src="bench-20160410/bench-avg-latency-ms.svg" alt="bench-20160410/bench-avg-latency-ms"><img src="bench-20160410/bench-throughput.svg" alt="bench-20160410/bench-throughput"><img src="bench-20160410/bench-avg-cpu.svg" alt="bench-20160410/bench-avg-cpu"><img src="bench-20160410/bench-avg-memory.svg" alt="bench-20160410/bench-avg-memory">

