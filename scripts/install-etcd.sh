#!/usr/bin/env bash
set -e

GIT_PATH=github.com/coreos/etcd

USER_NAME=coreos
BRANCH_NAME=master

rm -rf ${GOPATH}/src/${GIT_PATH}
mkdir -p ${GOPATH}/src/github.com/coreos

git clone https://github.com/${USER_NAME}/etcd \
  --branch ${BRANCH_NAME} \
  ${GOPATH}/src/${GIT_PATH}

cd ${GOPATH}/src/${GIT_PATH}

git reset --hard HEAD

./build

${GOPATH}/src/${GIT_PATH}/bin/etcd --version
${GOPATH}/src/${GIT_PATH}/bin/etcdctl --version

cp ${GOPATH}/src/${GIT_PATH}/bin/etcd ${GOPATH}/bin/etcd
sudo cp ${GOPATH}/src/${GIT_PATH}/bin/etcd /etcd

cp ${GOPATH}/src/${GIT_PATH}/bin/etcdctl ${GOPATH}/bin/etcdctl
sudo cp ${GOPATH}/src/${GIT_PATH}/bin/etcdctl /etcdctl

${GOPATH}/bin/etcd --version
ETCDCTL_API=3 ${GOPATH}/bin/etcdctl version
etc
ETCDCTL_API=3 etcdctl ersion
