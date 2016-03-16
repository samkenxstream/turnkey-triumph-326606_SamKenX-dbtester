# dbtester

[![Build Status](https://img.shields.io/travis/coreos/dbtester.svg?style=flat-square)][cistat] [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][dbtester-godoc]

Distributed database tester.

Run [`install.sh`](install.sh):

```
go get github.com/coreos/dbtester

# For each machine
dbtester agent --working-dir $WORKING_DIR --monitor
dbtester agent --working-dir $WORKING_DIR --monitor
dbtester agent --working-dir $WORKING_DIR --monitor
dbtester agent --working-dir $WORKING_DIR --monitor
dbtester agent --working-dir $WORKING_DIR --monitor

# Client machine
dbtester start --database="etcd" --agent-endpoints="$(echo $ETCD_RPC_ENDPOINTS)"
dbtester start --database="zk" --zk-max-client-conns=5000 --agent-endpoints="$(echo $ZK_RPC_ENDPOINTS)"
```

[cistat]: https://travis-ci.org/coreos/dbtester
[dbtester-godoc]: https://godoc.org/github.com/coreos/dbtester
