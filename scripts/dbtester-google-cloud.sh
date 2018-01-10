#!/usr/bin/env bash
set -e

##################################################
# create compute instances

gcloud compute instances create \
  bench-agent-1 \
  --custom-cpu=16 \
  --custom-memory=60 \
  --image-family=ubuntu-1710 \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=300 \
  --boot-disk-type="pd-ssd" \
  --network dbtester \
  --zone us-west1-a \
  --maintenance-policy=MIGRATE \
  --restart-on-failure

gcloud compute instances create \
  bench-agent-2 \
  --custom-cpu=16 \
  --custom-memory=60 \
  --image-family=ubuntu-1710 \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=300 \
  --boot-disk-type="pd-ssd" \
  --network dbtester \
  --zone us-west1-a \
  --maintenance-policy=MIGRATE \
  --restart-on-failure

gcloud compute instances create \
  bench-agent-3 \
  --custom-cpu=16 \
  --custom-memory=60 \
  --image-family=ubuntu-1710 \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=300 \
  --boot-disk-type="pd-ssd" \
  --network dbtester \
  --zone us-west1-a \
  --maintenance-policy=MIGRATE \
  --restart-on-failure

export GCP_KEY_PATH=/etc/gcp-key-etcd-development.json
gcloud compute instances create \
  bench-tester \
  --custom-cpu=16 \
  --custom-memory=60 \
  --image-family=ubuntu-1710 \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=300 \
  --boot-disk-type="pd-ssd" \
  --network dbtester \
  --zone us-west1-a \
  --maintenance-policy=MIGRATE \
  --restart-on-failure \
  --metadata-from-file gcp-key-etcd=${GCP_KEY_PATH}
##################################################


##################################################
gcloud compute ssh --zone=us-west1-a bench-agent-1
gcloud compute ssh --zone=us-west1-a bench-agent-2
gcloud compute ssh --zone=us-west1-a bench-agent-3
gcloud compute ssh --zone=us-west1-a bench-tester

gcloud compute instances list

<<COMMENT
10.138.0.2
10.138.0.3
10.138.0.4
10.138.0.5
COMMENT

<<COMMENT

sudo apt update -y

sudo apt install -y \
  build-essential \
  apt-utils \
  gcc \
  bash \
  bash-completion \
  tar \
  unzip \
  curl \
  wget \
  git \
  apt-transport-https \
  software-properties-common \
  libssl-dev \
  ntpdate

sudo apt upgrade -y

sudo apt autoremove
sudo apt autoclean

sudo service ntp stop
sudo ntpdate time.google.com
sudo service ntp start

COMMENT
##################################################


##################################################
ulimit -n
# ulimit -v unlimited

sudo vi /etc/security/limits.conf;

# add the following lines
* soft nofile 120000
* hard nofile 120000

sudo reboot
##################################################


sudo service ntp stop
sudo ntpdate time.google.com
sudo service ntp start


##################################################
sudo rm -rf ${HOME}/*

# etcd v3.2.0
# GO_VERSION=1.8.3

# etcd v3.3.0
GO_VERSION=1.9.2

sudo rm -f /usr/local/go/bin/go && sudo rm -rf /usr/local/go && sudo rm -f /bin/go

GOOGLE_URL=https://storage.googleapis.com/golang
DOWNLOAD_URL=${GOOGLE_URL}

sudo curl -s ${DOWNLOAD_URL}/go$GO_VERSION.linux-amd64.tar.gz | sudo tar -v -C /usr/local/ -xz

if grep -q GOPATH "$(echo $HOME)/.bashrc"; then
  echo "bashrc already has GOPATH";
else
  echo "adding GOPATH to bashrc";
  echo "export GOPATH=$(echo $HOME)/go" >> ${HOME}/.bashrc;
  PATH_VAR=$PATH":/usr/local/go/bin:$(echo $HOME)/go/bin";
  echo "export PATH=$(echo $PATH_VAR)" >> ${HOME}/.bashrc;
  source ${HOME}/.bashrc;
fi

mkdir -p ${GOPATH}/bin/
source ${HOME}/.bashrc
go version
##################################################


##################################################
GIT_PATH=github.com/coreos/etcd

USER_NAME=coreos
# BRANCH_NAME=release-3.2
BRANCH_NAME=release-3.3

rm -rf ${GOPATH}/src/${GIT_PATH}
mkdir -p ${GOPATH}/src/github.com/coreos

git clone https://github.com/${USER_NAME}/etcd \
  --branch ${BRANCH_NAME} \
  ${GOPATH}/src/${GIT_PATH}

cd ${GOPATH}/src/${GIT_PATH}

# v3.2.0
# git reset --hard 66722b1ada68fcd5227db853ee92003169a975c8

# v3.3.0 RC+
# git reset --hard 58c402a47bc5ad46fa748ce666257b45b06b1444

GO_BUILD_FLAGS="-v" ./build

${GOPATH}/src/${GIT_PATH}/bin/etcd --version
${GOPATH}/src/${GIT_PATH}/bin/etcdctl --version

cp ${GOPATH}/src/${GIT_PATH}/bin/etcd ${GOPATH}/bin/etcd
sudo cp ${GOPATH}/src/${GIT_PATH}/bin/etcd /etcd

cp ${GOPATH}/src/${GIT_PATH}/bin/etcdctl ${GOPATH}/bin/etcdctl
sudo cp ${GOPATH}/src/${GIT_PATH}/bin/etcdctl /etcdctl

${GOPATH}/bin/etcd --version
ETCDCTL_API=3 ${GOPATH}/bin/etcdctl version
etcd --version
ETCDCTL_API=3 etcdctl version
##################################################

# reinstall Go 1.9+ for context imports

##################################################
USER_NAME=coreos
BRANCH_NAME=master
cd ${HOME}
rm -rf ${HOME}/go/src/github.com/coreos/dbtester
git clone https://github.com/$USER_NAME/dbtester --branch $BRANCH_NAME ${HOME}/go/src/github.com/coreos/dbtester

cd ${HOME}
go install -v ./go/src/github.com/coreos/dbtester/cmd/dbtester

dbtester -h
dbtester agent -h
dbtester control -h
##################################################


##################################################
# agent on each machine
# specify network interface, disk device of host machine,
# this starts the database on host machine, when 'control' signals


sudo service ntp stop
sudo ntpdate time.google.com
sudo service ntp start

rm -f ${HOME}/agent.log
nohup dbtester agent \
  --agent-log ${HOME}/agent.log \
  --network-interface ens4 \
  --disk-device sda \
  --agent-port :3500 &

sleep 10s
cat ${HOME}/agent.log
##################################################


##################################################
# control on tester machine
# specify 'control' configuration file
# (client number, key number, key-value size)
# starts/shuts down database agents, send stress requests through RPCs

curl -L http://metadata.google.internal/computeMetadata/v1/instance/attributes/gcp-key-etcd -H 'Metadata-Flavor:Google' > /tmp/gcp-key-etcd-development.json
sudo mv /tmp/gcp-key-etcd-development.json /etc/gcp-key-etcd-development.json
head -10 /etc/gcp-key-etcd-development.json

# copy the tester configuration from git repository
cp ${HOME}/go/src/github.com/coreos/dbtester/test-configs/write-too-many-keys.yaml ${HOME}/config.yaml
cat ${HOME}/config.yaml



sudo service ntp stop
sudo ntpdate time.google.com
sudo service ntp start

nohup dbtester control \
  --database-id etcd__v3_3 \
  --config config.yaml > ${HOME}/client-control.log 2>&1 &

sleep 10s

tail -f ${HOME}/client-control.log

<<COMMENT
nohup dbtester control \
  --database-id etcd__v3_2 \
  --config config.yaml > ${HOME}/client-control.log 2>&1 &

nohup dbtester control \
  --database-id etcd__v3_3 \
  --config config.yaml > ${HOME}/client-control.log 2>&1 &

nohup dbtester control \
  --database-id zookeeper__r3_5_3_beta \
  --config config.yaml > ${HOME}/client-control.log 2>&1 &

nohup dbtester control \
  --database-id consul__v1_0_2 \
  --config config.yaml > ${HOME}/client-control.log 2>&1 &
COMMENT
##################################################


##################################################
# analyze; get all data from remote machines
# and specify 'analyze' configuration file,
# this aggregates data, generates all graphs, texts

cd ${HOME}/go/src/github.com/coreos/dbtester
go install -v ./cmd/dbtester


gsutil -m cp -R gs://dbtester-results/2018Q1-01-etcd .

cp ./test-configs/write-too-many-keys.yaml ./2018Q1-01-etcd/write-too-many-keys/

dbtester analyze --config 2018Q1-01-etcd/write-too-many-keys/write-too-many-keys.yaml

gsutil -m cp -R 2018Q1-01-etcd gs://dbtester-results/
gsutil -m acl ch -u allUsers:R -r gs://dbtester-results/2018Q1-01-etcd
##################################################


<<COMMENT
########################
etcd v3.2.0 (1,000 clients)
etcd v3.3.0 (1,000 clients)
Zookeeper r3.5.3-beta (500 clients)
Consul v1.0.2 (500 clients)
########################

########################
curl -X PURGE https://camo.githubusercontent.com/a6e057c6a9cff6a8d49f4c5b83a2b471bfc86817/68747470733a2f2f73746f726167652e676f6f676c65617069732e636f6d2f64627465737465722d726573756c74732f3230313751322d30312d657463642d7a6f6f6b65657065722d636f6e73756c2f30312d77726974652d314d2d6b6579732d636c69656e742d7661726961626c652f4156472d4c4154454e43592d4d532d42592d4b45592e737667
########################

########################
plotly

etcd v3.2.0 (Go 1.8.3)
etcd v3.3.0 (Go 1.9.2)
Zookeeper r3.5.3-beta (Java 8)
Consul v1.0.2 (Go 1.9.2)

Graph BoxPlot
Traces
Outlier, Show Statistics Mean

Filter
>0

Layout->Title and Fonts
Write 2M 256-byte key-value pairs, Latency
Write 500K 256-byte key-value pairs, Latency (1-connection)
font 30

Axes->Title->X
empty

Axes->Titles->Y
Latency (millisecond)
font 18

Axes->Tick Labels->X
Axes->Tick Labels->Y
font 17

Axes->Range
Auto Range, Log

Legend
font 22

Layout->Margins
Top 200px
########################
COMMENT

