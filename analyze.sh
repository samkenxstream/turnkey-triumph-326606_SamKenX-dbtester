#!/usr/bin/env bash

dbtester analyze \
	--output=testdata/bench-01-consul-aggregated.csv \
	--bench-file-path=testdata/bench-01-consul-timeseries.csv \
	--monitor-data-file-paths=testdata/bench-01-consul-1-monitor.csv,testdata/bench-01-consul-2-monitor.csv,testdata/bench-01-consul-3-monitor.csv

dbtester analyze \
	--output=testdata/bench-01-etcd-aggregated.csv \
	--bench-file-path=testdata/bench-01-etcd-timeseries.csv \
	--monitor-data-file-paths=testdata/bench-01-etcd-1-monitor.csv,testdata/bench-01-etcd-2-monitor.csv,testdata/bench-01-etcd-3-monitor.csv

dbtester analyze \
	--output=testdata/bench-01-etcd2-aggregated.csv \
	--bench-file-path=testdata/bench-01-etcd2-timeseries.csv \
	--monitor-data-file-paths=testdata/bench-01-etcd2-1-monitor.csv,testdata/bench-01-etcd2-2-monitor.csv,testdata/bench-01-etcd2-3-monitor.csv

dbtester analyze \
	--output=testdata/bench-01-zk-aggregated.csv \
	--bench-file-path=testdata/bench-01-zk-timeseries.csv \
	--monitor-data-file-paths=testdata/bench-01-zk-1-monitor.csv,testdata/bench-01-zk-2-monitor.csv,testdata/bench-01-zk-3-monitor.csv


dbtester analyze \
	--output=testdata/bench-01-all-aggregated.csv \
	--aggregated-file-paths=testdata/bench-01-consul-aggregated.csv,testdata/bench-01-etcd-aggregated.csv,testdata/bench-01-etcd2-aggregated.csv,testdata/bench-01-zk-aggregated.csv


dbtester analyze \
	--output=testdata/bench-01-plot \
	--image-format=png \
	--file-to-plot=testdata/bench-01-all-aggregated.csv

dbtester analyze \
	--output=testdata/bench-01-plot \
	--image-format=svg \
	--file-to-plot=testdata/bench-01-all-aggregated.csv

