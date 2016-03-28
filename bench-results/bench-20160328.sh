#!/usr/bin/env bash

gsutil -m cp -R gs://bench-20160328/ .

#######################################################################
# create aggregated csv files
dbtester analyze --output=bench-20160328/bench-01-etcd-aggregated.csv --bench-file-path=bench-20160328/bench-01-etcd-timeseries.csv --monitor-data-file-paths=bench-20160328/bench-01-etcd-1-monitor.csv,bench-20160328/bench-01-etcd-2-monitor.csv,bench-20160328/bench-01-etcd-3-monitor.csv
dbtester analyze --output=bench-20160328/bench-02-etcd-aggregated.csv --bench-file-path=bench-20160328/bench-02-etcd-timeseries.csv --monitor-data-file-paths=bench-20160328/bench-02-etcd-1-monitor.csv,bench-20160328/bench-02-etcd-2-monitor.csv,bench-20160328/bench-02-etcd-3-monitor.csv

# create agg/agg
dbtester analyze --same-database --output=bench-20160328/bench-all-aggregated.csv --aggregated-file-paths=bench-20160328/bench-01-etcd-aggregated.csv,bench-20160328/bench-02-etcd-aggregated.csv

# plot
IMAGE_TITLE="Write 3M keys, 700 conns, 1500 clients, key 64 bytes, value 256 bytes"
dbtester analyze --same-database --output=bench-20160328/bench-plot --image-format=png --file-to-plot=bench-20160328/bench-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"
dbtester analyze --same-database --output=bench-20160328/bench-plot --image-format=svg --file-to-plot=bench-20160328/bench-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"

#######################################################################
# generate README
dbtester readme --readme-dir=bench-20160328 --readme-preface=README_template
