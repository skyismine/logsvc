#!/usr/bin/env bash

# For each machine
ETCD_VERSION=v3.0.0
TOKEN=my-etcd-token
CLUSTER_STATE=new
NAME_1=etcd-node-0
NAME_2=etcd-node-1
NAME_3=etcd-node-2
# 获取宿主机ip地址
TAURUS_HOST=`ifconfig -a|grep inet|grep -v inet6|grep -v 172.17.*|grep -v 127.0.0.1|awk '{print $2}'|tr -d "addr:"`
HOST_1=${TAURUS_HOST}
HOST_2=${TAURUS_HOST}
HOST_3=${TAURUS_HOST}
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2381,${NAME_3}=http://${HOST_3}:2382

# For node 1
THIS_NAME=${NAME_1}
THIS_IP=${HOST_1}
docker run -itd --net=host --name ${THIS_NAME} quay.io/coreos/etcd:${ETCD_VERSION} \
    /usr/local/bin/etcd \
    --data-dir=data.etcd --name ${THIS_NAME} \
    --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 \
    --advertise-client-urls http://${THIS_IP}:2279 --listen-client-urls http://${THIS_IP}:2279 \
    --initial-cluster ${CLUSTER} \
    --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

# For node 2
THIS_NAME=${NAME_2}
THIS_IP=${HOST_2}
docker run -itd --net=host --name ${THIS_NAME} quay.io/coreos/etcd:${ETCD_VERSION} \
    /usr/local/bin/etcd \
    --data-dir=data.etcd --name ${THIS_NAME} \
    --initial-advertise-peer-urls http://${THIS_IP}:2381 --listen-peer-urls http://${THIS_IP}:2381 \
    --advertise-client-urls http://${THIS_IP}:2280 --listen-client-urls http://${THIS_IP}:2280 \
    --initial-cluster ${CLUSTER} \
    --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

# For node 3
THIS_NAME=${NAME_3}
THIS_IP=${HOST_3}
docker run -itd --net=host --name ${THIS_NAME} quay.io/coreos/etcd:${ETCD_VERSION} \
    /usr/local/bin/etcd \
    --data-dir=data.etcd --name ${THIS_NAME} \
    --initial-advertise-peer-urls http://${THIS_IP}:2382 --listen-peer-urls http://${THIS_IP}:2382 \
    --advertise-client-urls http://${THIS_IP}:2281 --listen-client-urls http://${THIS_IP}:2281 \
    --initial-cluster ${CLUSTER} \
    --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
