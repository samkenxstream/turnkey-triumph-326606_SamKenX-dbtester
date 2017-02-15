#!/usr/bin/env bash
set -e

if ! [[ "$0" =~ "scripts/tests.sh" ]]; then
    echo "must be run from repository root"
    exit 255
fi
gofmt -l -s -d *.go
TESTS="./analyze ./pkg/fileinspect ./pkg/ntp"

echo "Checking gofmt..."
fmtRes=$(gofmt -l -s -d $TESTS)
if [ -n "${fmtRes}" ]; then
	echo -e "gofmt checking failed:\n${fmtRes}"
	exit 255
fi

echo "Checking govet..."
vetRes=$(go vet $TESTS 2>&1 >/dev/null)
if [ -n "${vetRes}" ]; then
	echo -e "govet checking failed:\n${vetRes}"
	exit 255
fi

echo "Running tests...";
go test -v $TESTS;
go test -v -race $TESTS;

echo "Success";
