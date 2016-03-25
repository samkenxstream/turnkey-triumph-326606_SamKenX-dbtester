#!/usr/bin/env bash

gsutil -m cp -R gs://bench-20160325/ .

#######################################################################
# create aggregated csv files
dbtester analyze --output=bench-20160325/bench-01-consul-aggregated.csv --bench-file-path=bench-20160325/bench-01-consul-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-01-consul-1-monitor.csv,bench-20160325/bench-01-consul-2-monitor.csv,bench-20160325/bench-01-consul-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-01-etcd-aggregated.csv --bench-file-path=bench-20160325/bench-01-etcd-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-01-etcd-1-monitor.csv,bench-20160325/bench-01-etcd-2-monitor.csv,bench-20160325/bench-01-etcd-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-01-etcd2-aggregated.csv --bench-file-path=bench-20160325/bench-01-etcd2-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-01-etcd2-1-monitor.csv,bench-20160325/bench-01-etcd2-2-monitor.csv,bench-20160325/bench-01-etcd2-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-01-zk-aggregated.csv --bench-file-path=bench-20160325/bench-01-zk-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-01-zk-1-monitor.csv,bench-20160325/bench-01-zk-2-monitor.csv,bench-20160325/bench-01-zk-3-monitor.csv

# create agg/agg
dbtester analyze --output=bench-20160325/bench-01-all-aggregated.csv --aggregated-file-paths=bench-20160325/bench-01-consul-aggregated.csv,bench-20160325/bench-01-etcd-aggregated.csv,bench-20160325/bench-01-etcd2-aggregated.csv,bench-20160325/bench-01-zk-aggregated.csv

# plot
IMAGE_TITLE="Write 300K keys, 1 client, key 64 bytes, value 256 bytes"
dbtester analyze --output=bench-20160325/bench-01-plot --image-format=png --file-to-plot=bench-20160325/bench-01-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"
dbtester analyze --output=bench-20160325/bench-01-plot --image-format=svg --file-to-plot=bench-20160325/bench-01-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"


#######################################################################
# create aggregated csv files
dbtester analyze --output=bench-20160325/bench-02-consul-aggregated.csv --bench-file-path=bench-20160325/bench-02-consul-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-02-consul-1-monitor.csv,bench-20160325/bench-02-consul-2-monitor.csv,bench-20160325/bench-02-consul-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-02-etcd-aggregated.csv --bench-file-path=bench-20160325/bench-02-etcd-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-02-etcd-1-monitor.csv,bench-20160325/bench-02-etcd-2-monitor.csv,bench-20160325/bench-02-etcd-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-02-etcd2-aggregated.csv --bench-file-path=bench-20160325/bench-02-etcd2-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-02-etcd2-1-monitor.csv,bench-20160325/bench-02-etcd2-2-monitor.csv,bench-20160325/bench-02-etcd2-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-02-zk-aggregated.csv --bench-file-path=bench-20160325/bench-02-zk-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-02-zk-1-monitor.csv,bench-20160325/bench-02-zk-2-monitor.csv,bench-20160325/bench-02-zk-3-monitor.csv

# create agg/agg
dbtester analyze --output=bench-20160325/bench-02-all-aggregated.csv --aggregated-file-paths=bench-20160325/bench-02-consul-aggregated.csv,bench-20160325/bench-02-etcd-aggregated.csv,bench-20160325/bench-02-etcd2-aggregated.csv,bench-20160325/bench-02-zk-aggregated.csv

# plot
IMAGE_TITLE="Write 1M keys, 10 clients, key 64 bytes, value 256 bytes"
dbtester analyze --output=bench-20160325/bench-02-plot --image-format=png --file-to-plot=bench-20160325/bench-02-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"
dbtester analyze --output=bench-20160325/bench-02-plot --image-format=svg --file-to-plot=bench-20160325/bench-02-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"


#######################################################################
# create aggregated csv files
dbtester analyze --output=bench-20160325/bench-03-consul-aggregated.csv --bench-file-path=bench-20160325/bench-03-consul-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-03-consul-1-monitor.csv,bench-20160325/bench-03-consul-2-monitor.csv,bench-20160325/bench-03-consul-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-03-etcd-aggregated.csv --bench-file-path=bench-20160325/bench-03-etcd-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-03-etcd-1-monitor.csv,bench-20160325/bench-03-etcd-2-monitor.csv,bench-20160325/bench-03-etcd-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-03-etcd2-aggregated.csv --bench-file-path=bench-20160325/bench-03-etcd2-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-03-etcd2-1-monitor.csv,bench-20160325/bench-03-etcd2-2-monitor.csv,bench-20160325/bench-03-etcd2-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-03-zk-aggregated.csv --bench-file-path=bench-20160325/bench-03-zk-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-03-zk-1-monitor.csv,bench-20160325/bench-03-zk-2-monitor.csv,bench-20160325/bench-03-zk-3-monitor.csv

# create agg/agg
dbtester analyze --output=bench-20160325/bench-03-all-aggregated.csv --aggregated-file-paths=bench-20160325/bench-03-consul-aggregated.csv,bench-20160325/bench-03-etcd-aggregated.csv,bench-20160325/bench-03-etcd2-aggregated.csv,bench-20160325/bench-03-zk-aggregated.csv

# plot
IMAGE_TITLE="Write 3M keys, 500 clients, key 64 bytes, value 256 bytes"
dbtester analyze --output=bench-20160325/bench-03-plot --image-format=png --file-to-plot=bench-20160325/bench-03-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"
dbtester analyze --output=bench-20160325/bench-03-plot --image-format=svg --file-to-plot=bench-20160325/bench-03-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"


#######################################################################
# create aggregated csv files
dbtester analyze --output=bench-20160325/bench-04-consul-aggregated.csv --bench-file-path=bench-20160325/bench-04-consul-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-04-consul-1-monitor.csv,bench-20160325/bench-04-consul-2-monitor.csv,bench-20160325/bench-04-consul-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-04-etcd-aggregated.csv --bench-file-path=bench-20160325/bench-04-etcd-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-04-etcd-1-monitor.csv,bench-20160325/bench-04-etcd-2-monitor.csv,bench-20160325/bench-04-etcd-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-04-etcdmulti-aggregated.csv --bench-file-path=bench-20160325/bench-04-etcdmulti-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-04-etcdmulti-1-monitor.csv,bench-20160325/bench-04-etcdmulti-2-monitor.csv,bench-20160325/bench-04-etcdmulti-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-04-etcd2-aggregated.csv --bench-file-path=bench-20160325/bench-04-etcd2-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-04-etcd2-1-monitor.csv,bench-20160325/bench-04-etcd2-2-monitor.csv,bench-20160325/bench-04-etcd2-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-04-zk-aggregated.csv --bench-file-path=bench-20160325/bench-04-zk-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-04-zk-1-monitor.csv,bench-20160325/bench-04-zk-2-monitor.csv,bench-20160325/bench-04-zk-3-monitor.csv

# create agg/agg
dbtester analyze --output=bench-20160325/bench-04-all-aggregated.csv --aggregated-file-paths=bench-20160325/bench-04-consul-aggregated.csv,bench-20160325/bench-04-etcd-aggregated.csv,bench-20160325/bench-04-etcdmulti-aggregated.csv,bench-20160325/bench-04-etcd2-aggregated.csv,bench-20160325/bench-04-zk-aggregated.csv

# plot
IMAGE_TITLE="Write 3M keys, 1000 clients, key 64 bytes, value 256 bytes"
MULTI_TAG_TITLE="100-conns-1k-clients"
dbtester analyze --output=bench-20160325/bench-04-plot --image-format=png --file-to-plot=bench-20160325/bench-04-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)" --multi-tag-title="$(echo $MULTI_TAG_TITLE)"
dbtester analyze --output=bench-20160325/bench-04-plot --image-format=svg --file-to-plot=bench-20160325/bench-04-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)" --multi-tag-title="$(echo $MULTI_TAG_TITLE)"


#######################################################################
# create aggregated csv files
dbtester analyze --output=bench-20160325/bench-05-consul-aggregated.csv --bench-file-path=bench-20160325/bench-05-consul-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-05-consul-1-monitor.csv,bench-20160325/bench-05-consul-2-monitor.csv,bench-20160325/bench-05-consul-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-05-etcd-aggregated.csv --bench-file-path=bench-20160325/bench-05-etcd-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-05-etcd-1-monitor.csv,bench-20160325/bench-05-etcd-2-monitor.csv,bench-20160325/bench-05-etcd-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-05-etcd2-aggregated.csv --bench-file-path=bench-20160325/bench-05-etcd2-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-05-etcd2-1-monitor.csv,bench-20160325/bench-05-etcd2-2-monitor.csv,bench-20160325/bench-05-etcd2-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-05-zk-aggregated.csv --bench-file-path=bench-20160325/bench-05-zk-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-05-zk-1-monitor.csv,bench-20160325/bench-05-zk-2-monitor.csv,bench-20160325/bench-05-zk-3-monitor.csv

# create agg/agg
dbtester analyze --output=bench-20160325/bench-05-all-aggregated.csv --aggregated-file-paths=bench-20160325/bench-05-consul-aggregated.csv,bench-20160325/bench-05-etcd-aggregated.csv,bench-20160325/bench-05-etcd2-aggregated.csv,bench-20160325/bench-05-zk-aggregated.csv

# plot
IMAGE_TITLE="Read 1M keys, 1 client, key 64 bytes, value 1 kb"
dbtester analyze --output=bench-20160325/bench-05-plot --image-format=png --file-to-plot=bench-20160325/bench-05-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"
dbtester analyze --output=bench-20160325/bench-05-plot --image-format=svg --file-to-plot=bench-20160325/bench-05-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"


#######################################################################
# create aggregated csv files
dbtester analyze --output=bench-20160325/bench-06-consul-aggregated.csv --bench-file-path=bench-20160325/bench-06-consul-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-06-consul-1-monitor.csv,bench-20160325/bench-06-consul-2-monitor.csv,bench-20160325/bench-06-consul-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-06-etcd-aggregated.csv --bench-file-path=bench-20160325/bench-06-etcd-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-06-etcd-1-monitor.csv,bench-20160325/bench-06-etcd-2-monitor.csv,bench-20160325/bench-06-etcd-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-06-etcd2-aggregated.csv --bench-file-path=bench-20160325/bench-06-etcd2-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-06-etcd2-1-monitor.csv,bench-20160325/bench-06-etcd2-2-monitor.csv,bench-20160325/bench-06-etcd2-3-monitor.csv
dbtester analyze --output=bench-20160325/bench-06-zk-aggregated.csv --bench-file-path=bench-20160325/bench-06-zk-timeseries.csv --monitor-data-file-paths=bench-20160325/bench-06-zk-1-monitor.csv,bench-20160325/bench-06-zk-2-monitor.csv,bench-20160325/bench-06-zk-3-monitor.csv

# create agg/agg
dbtester analyze --output=bench-20160325/bench-06-all-aggregated.csv --aggregated-file-paths=bench-20160325/bench-06-consul-aggregated.csv,bench-20160325/bench-06-etcd-aggregated.csv,bench-20160325/bench-06-etcd2-aggregated.csv,bench-20160325/bench-06-zk-aggregated.csv

# plot
IMAGE_TITLE="Read 1M keys, 1000 clients, key 64 bytes, value 1 kb"
dbtester analyze --output=bench-20160325/bench-06-plot --image-format=png --file-to-plot=bench-20160325/bench-06-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"
dbtester analyze --output=bench-20160325/bench-06-plot --image-format=svg --file-to-plot=bench-20160325/bench-06-all-aggregated.csv --image-title="$(echo $IMAGE_TITLE)"


#######################################################################
# generate README
dbtester readme --readme-dir=bench-20160325 --readme-preface=README_template
