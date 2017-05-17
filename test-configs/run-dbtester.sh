#!/usr/bin/env bash
set -e

# agent; specify network interface, disk device of host machine,
# this starts the database on host machine, when 'control' signals
nohup dbtester agent --network-interface ens4  --disk-device sda  --agent-port :3500 &

# control; specify 'control' configuration file (client number, key number, key-value size),
# this starts database stressing, and shuts down the database when done
nohup dbtester control --database-id etcd__tip --config config.yaml > $HOME/control.log 2>&1 &
nohup dbtester control --database-id zookeeper__r3_5_3_beta --config config.yaml > $HOME/control.log 2>&1 &
nohup dbtester control --database-id consul__v0_8_3 --config config.yaml > $HOME/control.log 2>&1 &

# analyze; get all data from remote machines
# and specify 'analyze' configuration file,
# this aggregates data, generates all graphs, texts
dbtester analyze --config config.yaml
