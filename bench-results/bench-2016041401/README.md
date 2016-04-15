

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



<br><br><hr>
##### Write 500K keys, 1 client, key 32 bytes, value 500 bytes

<img src="https://storage.googleapis.com/bench-2016041401/bench-01-avg-latency-ms.svg" alt="bench-2016041401/bench-01-avg-latency-ms">

<img src="https://storage.googleapis.com/bench-2016041401/bench-01-throughput.svg" alt="bench-2016041401/bench-01-throughput">

<img src="https://storage.googleapis.com/bench-2016041401/bench-01-avg-cpu.svg" alt="bench-2016041401/bench-01-avg-cpu">

<img src="https://storage.googleapis.com/bench-2016041401/bench-01-avg-memory.svg" alt="bench-2016041401/bench-01-avg-memory">



<br><br><hr>
##### Write 2M keys, 1000 clients (etcd 100 conns), key 32 bytes, value 500 bytes

<img src="https://storage.googleapis.com/bench-2016041401/bench-02-avg-latency-ms.svg" alt="bench-2016041401/bench-02-avg-latency-ms">

<img src="https://storage.googleapis.com/bench-2016041401/bench-02-throughput.svg" alt="bench-2016041401/bench-02-throughput">

<img src="https://storage.googleapis.com/bench-2016041401/bench-02-avg-cpu.svg" alt="bench-2016041401/bench-02-avg-cpu">

<img src="https://storage.googleapis.com/bench-2016041401/bench-02-avg-memory.svg" alt="bench-2016041401/bench-02-avg-memory">



