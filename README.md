# dbtester

[![Build Status](https://img.shields.io/travis/coreos/dbtester.svg?style=flat-square)][cistat] [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][dbtester-godoc]

Distributed database tester.

Please databases as in [`install.sh`](install.sh). And:

```
go get github.com/coreos/dbtester

# For each machine
dbtester agent
dbtester agent
dbtester agent
dbtester agent
dbtester agent

# Client machine
dbtester start --agent-endpoints=$(echo $RPC_ENDPOINTS) --database=etcd 
dbtester start --agent-endpoints=$(echo $RPC_ENDPOINTS) --database=etcd2 
dbtester start --agent-endpoints=$(echo $RPC_ENDPOINTS) --database=zk
dbtester start --agent-endpoints=$(echo $RPC_ENDPOINTS) --database=consul
```

[cistat]: https://travis-ci.org/coreos/dbtester
[dbtester-godoc]: https://godoc.org/github.com/coreos/dbtester
