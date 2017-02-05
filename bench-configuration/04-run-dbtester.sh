#!/usr/bin/env bash
set -e

# agent; specify network interface, disk device on the machine
# this starts the database on host machine when signalled from 'control'
nohup dbtester agent --network-interface ens4  --disk-device sda1  --agent-port :3500 &

# control; specify 'control' configuration file
# this starts client stressing
nohup dbtester control -c config.yaml > $HOME/control.log 2>&1 &

# analyze; get all data from remote machines
# specify 'analyze' configuration file
# this aggregates, generates all graphs, texts
dbtester analyze --config analyze.yaml
