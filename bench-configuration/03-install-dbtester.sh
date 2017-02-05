#!/usr/bin/env bash
set -e

USER_NAME=coreos
BRANCH_NAME=master
cd $HOME
rm -rf $HOME/go/src/github.com/coreos/dbtester
git clone https://github.com/$USER_NAME/dbtester --branch $BRANCH_NAME $HOME/go/src/github.com/coreos/dbtester

cd $HOME
go install -v ./go/src/github.com/coreos/dbtester

dbtester agent -h
dbtester control -h
