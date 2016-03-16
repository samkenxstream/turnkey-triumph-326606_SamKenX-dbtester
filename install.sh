#!/usr/bin/env bash
set -e

WORKING_DIR=$HOME
if [ -n "$1" ]; then
  WORKING_DIR=$1
fi;
echo Setting working directory to $WORKING_DIR
cd $WORKING_DIR


######################
echo Installing etcd #
######################
go get github.com/coreos/etcd


#########################
echo Installing Ansible #
#########################
sudo apt-get install -y software-properties-common
sudo apt-add-repository -y ppa:ansible/ansible
sudo apt-get -y update
sudo apt-get install -y ansible
echo "---
- name: a play that runs entirely on the ansible host
  hosts: 127.0.0.1
  connection: local
  tasks:
  - name: Install Linux utils
    become: yes
    apt: name={{item}} state=latest
    with_items:
      - bash
      - curl
      - git
      - tar
      - iptables
      - iproute2

  - name: Install add-apt-repostory
    become: yes
    apt: name=software-properties-common state=latest

  - name: Add Oracle Java Repository
    become: yes
    apt_repository: repo='ppa:webupd8team/java'

  - name: Accept Java 8 License
    become: yes
    debconf: name='oracle-java8-installer' question='shared/accepted-oracle-license-v1-1' value='true' vtype='select'

  - name: Install Oracle Java 8
    become: yes
    apt: name={{item}} state=latest
    with_items:
      - oracle-java8-installer
      - ca-certificates
      - oracle-java8-set-default

  - name: Print Java version
    command: java -version
    register: result
  - debug:
      var: result.stderr

  - name: Print JDK version
    command: javac -version
    register: result
  - debug:
      var: result.stderr
" > $WORKING_DIR/dbtest_install.yml
ansible-playbook $WORKING_DIR/dbtest_install.yml


###########################
echo Installing Zookeeper #
###########################
ZOOKEEPER_VERSION=3.4.8
sudo rm -rf $WORKING_DIR/zookeeper
sudo curl -sf -o /tmp/zookeeper-$ZOOKEEPER_VERSION.tar.gz -L https://www.apache.org/dist/zookeeper/zookeeper-$ZOOKEEPER_VERSION/zookeeper-$ZOOKEEPER_VERSION.tar.gz
sudo tar -xzf /tmp/zookeeper-$ZOOKEEPER_VERSION.tar.gz -C /tmp/
sudo mv /tmp/zookeeper-$ZOOKEEPER_VERSION /tmp/zookeeper
sudo mv /tmp/zookeeper $WORKING_DIR/
sudo chmod -R 777 $WORKING_DIR/zookeeper/
mkdir -p $WORKING_DIR/zookeeper/data.zk
touch $WORKING_DIR/zookeeper/data.zk/myid
sudo chmod -R 777 $WORKING_DIR/zookeeper/data.zk/
