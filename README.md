This is an experimental project. Project/code is subject to change anytime.
This is mainly for comparing etcd with other databases. For etcd, we recommend
to just run [etcd benchmark tool](https://github.com/coreos/etcd/tree/master/tools/benchmark).

# dbtester

[![Build Status](https://img.shields.io/travis/coreos/dbtester.svg?style=flat-square)][cistat] [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][dbtester-godoc]

Distributed database tester.

Please let us know or file an issue if:

- Need help with running this testing suite
- Questions about test results

We keep full logs here and cloud storage(when it's over 1MB):

- https://console.cloud.google.com/storage/browser/dbtester-results

Test results:

- Compression experiment: https://github.com/coreos/dbtester/tree/master/bench-results/2016051001/README.md
- Compression experiment: https://github.com/coreos/dbtester/tree/master/bench-results/2016050901/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016050504/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016050503/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016050502/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016050501/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016050301/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016050101/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016043002/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016043001/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016042502/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016042501/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016041801/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016041601/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016041502/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016041501/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016041401/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016041203/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016041202/README.md
- https://github.com/coreos/dbtester/tree/master/bench-results/2016041201/README.md

[cistat]: https://travis-ci.org/coreos/dbtester
[dbtester-godoc]: https://godoc.org/github.com/coreos/dbtester
