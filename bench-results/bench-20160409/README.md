

- Google Cloud Compute Engine
- 8 vCPUs + 16GB Memory + 50GB SSD
- 1 machine(client) of 16 vCPUs + 30GB Memory + 50GB SSD
- Ubuntu 15.10
- Go 1.6
- Java 8
  - Java(TM) SE Runtime Environment (build 1.8.0_74-b02)
  - Java HotSpot(TM) 64-Bit Server VM (build 25.74-b02, mixed mode)
- etcd v3 (master branch)
- Zookeeper v3.4.8
- Consul v0.6.4



<br><br><hr>
##### Write 300K keys, 1K clients, key 64 bytes, value 256 bytes

<img src="https://storage.googleapis.com/bench-20160409/bench-avg-latency-ms.svg" alt="bench-20160409/bench-avg-latency-ms">

<img src="https://storage.googleapis.com/bench-20160409/bench-throughput.svg" alt="bench-20160409/bench-throughput">

<img src="https://storage.googleapis.com/bench-20160409/bench-avg-cpu.svg" alt="bench-20160409/bench-avg-cpu">

<img src="https://storage.googleapis.com/bench-20160409/bench-avg-memory.svg" alt="bench-20160409/bench-avg-memory">



