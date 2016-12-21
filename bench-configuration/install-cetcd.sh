#!/usr/bin/env bash
set -e

go get -v github.com/coreos/cetcd/cmd/cetcd
cetcd -h
