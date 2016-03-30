#!/usr/bin/env bash

gsutil -m cp -R gs://bench-20160330/ .

#######################################################################
# create aggregated csv files
dbtester analyze --output=bench-20160330/bench-01-etcd-aggregated.csv --bench-file-path=bench-20160330/bench-01-etcd-timeseries.csv --monitor-data-file-paths=bench-20160330/bench-01-etcd-1-monitor.csv,bench-20160330/bench-01-etcd-2-monitor.csv,bench-20160330/bench-01-etcd-3-monitor.csv
dbtester analyze --output=bench-20160330/bench-01-zk-aggregated.csv --bench-file-path=bench-20160330/bench-01-zk-timeseries.csv --monitor-data-file-paths=bench-20160330/bench-01-zk-1-monitor.csv,bench-20160330/bench-01-zk-2-monitor.csv,bench-20160330/bench-01-zk-3-monitor.csv

# create agg/agg
dbtester analyze --output=bench-20160330/bench-01-all-aggregated.csv --aggregated-file-paths=bench-20160330/bench-01-etcd-aggregated.csv,bench-20160330/bench-01-zk-aggregated.csv

# plot
IMAGE_TITLE="Write 300K keys, 1 client, key 64 bytes, value 256 bytes"
dbtester analyze --output=bench-20160330/bench-01-plot --image-format=png --file-to-plot=bench-20160330/bench-01-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"
dbtester analyze --output=bench-20160330/bench-01-plot --image-format=svg --file-to-plot=bench-20160330/bench-01-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"

#######################################################################
# create aggregated csv files
dbtester analyze --output=bench-20160330/bench-02-etcd-aggregated.csv --bench-file-path=bench-20160330/bench-02-etcd-timeseries.csv --monitor-data-file-paths=bench-20160330/bench-02-etcd-1-monitor.csv,bench-20160330/bench-02-etcd-2-monitor.csv,bench-20160330/bench-02-etcd-3-monitor.csv
dbtester analyze --output=bench-20160330/bench-02-zk-aggregated.csv --bench-file-path=bench-20160330/bench-02-zk-timeseries.csv --monitor-data-file-paths=bench-20160330/bench-02-zk-1-monitor.csv,bench-20160330/bench-02-zk-2-monitor.csv,bench-20160330/bench-02-zk-3-monitor.csv

# create agg/agg
dbtester analyze --output=bench-20160330/bench-02-all-aggregated.csv --aggregated-file-paths=bench-20160330/bench-02-etcd-aggregated.csv,bench-20160330/bench-02-zk-aggregated.csv

# plot
IMAGE_TITLE="Write 3M keys, 1K clients, key 64 bytes, value 256 bytes"
dbtester analyze --output=bench-20160330/bench-02-plot --image-format=png --file-to-plot=bench-20160330/bench-02-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"
dbtester analyze --output=bench-20160330/bench-02-plot --image-format=svg --file-to-plot=bench-20160330/bench-02-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"

#######################################################################
# generate README
dbtester readme --readme-dir=bench-20160330 --readme-preface=README_template
