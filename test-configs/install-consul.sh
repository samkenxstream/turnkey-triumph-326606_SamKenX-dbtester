#!/usr/bin/env bash
set -e

rm -f /tmp/consul.zip
curl -sf -o /tmp/consul.zip https://releases.hashicorp.com/consul/0.8.3/consul_0.8.3_linux_amd64.zip

rm -f ${GOPATH}/bin/consul
unzip /tmp/consul.zip -d ${GOPATH}/bin
rm -f /tmp/consul.zip

consul version

<<COMMENT
https://github.com/hashicorp/consul/releases
COMMENT
