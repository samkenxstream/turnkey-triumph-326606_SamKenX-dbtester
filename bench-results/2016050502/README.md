

- Google Cloud Compute Engine
- 8 vCPUs + 16GB Memory + 50GB SSD
- 1 machine(client) of 16 vCPUs + 30GB Memory + 50GB SSD
- Ubuntu 15.10
- Zookeeper v3.4.8
- Consul v0.6.4



<br><br><hr>
##### Write 600K keys, 1000-client, 1000-conn, 8-byte same key, 256-byte value (Zookeeper)

<img src="https://storage.googleapis.com/dbtester-results/2016050502/01-avg-latency-ms.svg" alt="2016050502/01-avg-latency-ms">

<img src="https://storage.googleapis.com/dbtester-results/2016050502/01-throughput.svg" alt="2016050502/01-throughput">

<img src="https://storage.googleapis.com/dbtester-results/2016050502/01-avg-cpu.svg" alt="2016050502/01-avg-cpu">

<img src="https://storage.googleapis.com/dbtester-results/2016050502/01-avg-memory.svg" alt="2016050502/01-avg-memory">



<br><br><hr>
##### Write 600K keys, 1000-client, 1000-conn, 8-byte same key, 256-byte value (Consul)

<img src="https://storage.googleapis.com/dbtester-results/2016050502/02-avg-latency-ms.svg" alt="2016050502/02-avg-latency-ms">

<img src="https://storage.googleapis.com/dbtester-results/2016050502/02-throughput.svg" alt="2016050502/02-throughput">

<img src="https://storage.googleapis.com/dbtester-results/2016050502/02-avg-cpu.svg" alt="2016050502/02-avg-cpu">

<img src="https://storage.googleapis.com/dbtester-results/2016050502/02-avg-memory.svg" alt="2016050502/02-avg-memory">



