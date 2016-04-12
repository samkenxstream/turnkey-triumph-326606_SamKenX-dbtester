

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
##### Write 300K keys, 1 client, key 64 bytes, value 256 bytes

<img src="https://storage.googleapis.com/bench-2016041201/bench-01-avg-latency-ms.svg" alt="bench-2016041201/bench-01-avg-latency-ms">

<img src="https://storage.googleapis.com/bench-2016041201/bench-01-throughput.svg" alt="bench-2016041201/bench-01-throughput">

<img src="https://storage.googleapis.com/bench-2016041201/bench-01-avg-cpu.svg" alt="bench-2016041201/bench-01-avg-cpu">

<img src="https://storage.googleapis.com/bench-2016041201/bench-01-avg-memory.svg" alt="bench-2016041201/bench-01-avg-memory">



<br><br><hr>
##### Write 1M keys, 300 clients(50 etcd conns), key 64 bytes, value 1 kb

<img src="https://storage.googleapis.com/bench-2016041201/bench-02-avg-latency-ms.svg" alt="bench-2016041201/bench-02-avg-latency-ms">

<img src="https://storage.googleapis.com/bench-2016041201/bench-02-throughput.svg" alt="bench-2016041201/bench-02-throughput">

<img src="https://storage.googleapis.com/bench-2016041201/bench-02-avg-cpu.svg" alt="bench-2016041201/bench-02-avg-cpu">

<img src="https://storage.googleapis.com/bench-2016041201/bench-02-avg-memory.svg" alt="bench-2016041201/bench-02-avg-memory">



<br><br><hr>
##### Write 1M keys, 1K clients(50 etcd conns), key 8 bytes, value 256 bytes

<img src="https://storage.googleapis.com/bench-2016041201/bench-03-avg-latency-ms.svg" alt="bench-2016041201/bench-03-avg-latency-ms">

<img src="https://storage.googleapis.com/bench-2016041201/bench-03-throughput.svg" alt="bench-2016041201/bench-03-throughput">

<img src="https://storage.googleapis.com/bench-2016041201/bench-03-avg-cpu.svg" alt="bench-2016041201/bench-03-avg-cpu">

<img src="https://storage.googleapis.com/bench-2016041201/bench-03-avg-memory.svg" alt="bench-2016041201/bench-03-avg-memory">



<br><br><hr>
##### Write 2M keys, 1K clients(50 etcd conns), key 64 bytes, value 1 kb

<img src="https://storage.googleapis.com/bench-2016041201/bench-04-avg-latency-ms.svg" alt="bench-2016041201/bench-04-avg-latency-ms">

<img src="https://storage.googleapis.com/bench-2016041201/bench-04-throughput.svg" alt="bench-2016041201/bench-04-throughput">

<img src="https://storage.googleapis.com/bench-2016041201/bench-04-avg-cpu.svg" alt="bench-2016041201/bench-04-avg-cpu">

<img src="https://storage.googleapis.com/bench-2016041201/bench-04-avg-memory.svg" alt="bench-2016041201/bench-04-avg-memory">



