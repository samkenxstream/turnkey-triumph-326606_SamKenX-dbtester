#!/usr/bin/env bash
set -e

<<COMMENT
curl https://.../dbtester_agent.sh | bash -s f /mnt/ssd0
COMMENT

# clean page cache
echo "echo 1 > /proc/sys/vm/drop_caches" | sudo sh

if [ "$1" = "f" ]; then
  echo "Formatting following https://cloud.google.com/compute/docs/disks/local-ssd"
  sudo mkdir -p /mnt/ssd0
  ls /dev/disk/by-id/
  sudo mkfs.ext4 -F /dev/disk/by-id/google-local-ssd-0
  sudo mount -o discard,defaults /dev/disk/by-id/google-local-ssd-0 /mnt/ssd0
  sudo chmod a+w /mnt/ssd0
  echo '/dev/disk/by-id/google-local-ssd-0 /mnt/ssd0 ext4 defaults 1 1' | sudo tee -a /etc/fstab
  echo 'Hello, World!' > /mnt/ssd0/hello.txt
  cat /mnt/ssd0/hello.txt
  df -h
  rm -f /mnt/ssd0/hello.txt
  ls /mnt/ssd0/
else
  echo "No formatting"
fi;

WORKING_DIR=$HOME
if [ -n "$2" ]; then
  WORKING_DIR=$2
fi;
echo Setting working directory to $WORKING_DIR
cd $WORKING_DIR
sleep 3s

#############################
echo Installing Go programs #
#############################
GO_VERSION="1.6" && cd /usr/local && sudo rm -rf ./go && sudo curl -s https://storage.googleapis.com/golang/go$GO_VERSION.linux-amd64.tar.gz | sudo tar -v -C /usr/local/ -xz && cd $HOME;
echo "export GOPATH=$(echo $HOME)/go" >> $HOME/.bashrc
PATH_VAR=$PATH":/usr/local/go/bin:$(echo $HOME)/go/bin"
echo "export PATH=$(echo $PATH_VAR)" >> $HOME/.bashrc
export GOPATH=$(echo $HOME)/go
PATH_VAR=$PATH":/usr/local/go/bin:$(echo $HOME)/go/bin"
export PATH=$(echo $PATH_VAR)

go get -v -u -f github.com/gyuho/psn
psn ps-kill --force --clean-up
psn ps-kill --force --program etcd
psn ps-kill --force --program dbtester
psn ps-kill --force --program java
psn ss-kill --local-port 2181
psn ss-kill --local-port 2378
psn ss-kill --local-port 2379
psn ss-kill --local-port 2380
psn ss-kill --local-port 2888
psn ss-kill --local-port 3500
psn ss-kill --local-port 3888
rm -rf $WORKING_DIR/*.etcd
rm -rf $WORKING_DIR/*.log
rm -rf $WORKING_DIR/failure_archive
go get -v -u -f github.com/coreos/etcd
go get -v -u -f github.com/coreos/dbtester/bench
go get -v -u -f github.com/coreos/dbtester

if [ -n "$3" ]; then
  echo "Running node_exporter..."
  psn ps-kill --force --program node_exporter
  psn ss-kill --local-port 9100
  go get -v -u -f github.com/prometheus/node_exporter
  nohup node_exporter -web.listen-address=":9100" > $WORKING_DIR/node_exporter.log 2>&1 &
  sleep 5s
  psn ps --program node_exporter
  psn ss --program node_exporter
else
  echo "No node_exporter..."
fi;

#########################
echo Installing Ansible #
#########################
sudo apt-get -y update
sudo apt-get -y upgrade
sudo apt-get -y autoremove
sudo apt-get -y autoclean
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
" > $WORKING_DIR/dbtester_install.yml
ansible-playbook $WORKING_DIR/dbtester_install.yml

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
