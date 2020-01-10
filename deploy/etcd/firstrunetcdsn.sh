#!/usr/bin/env bash

# For each machine
ETCD_VERSION=v3.0.0
NODE_NAME=etcd-node-single

docker stop ${NODE_NAME}
docker container rm ${NODE_NAME}

# 配置为不限制node访问
ALL_PASS_ADDR=0.0.0.0

# 获取宿主机ip地址(局域网地址)
TAURUS_HOST=`ifconfig -a|grep inet|grep -v inet6|grep -v 172.17.*|grep -v 172.18.*|grep -v 127.0.0.1|awk '{print $2}'|tr -d "addr:"`
LOCAL_ADDR=${TAURUS_HOST}

# 运行etcd容器
docker run -itd --net=host --name ${NODE_NAME} -v /opt/data/etcd:/data/etcd quay.io/coreos/etcd:${ETCD_VERSION} \
    /usr/local/bin/etcd \
    --name ${NODE_NAME} \
    --data-dir=/data/etcd/${NODE_NAME} \
    --listen-peer-urls http://${LOCAL_ADDR}:2380 \
    --advertise-client-urls http://${LOCAL_ADDR}:2379 \
    --listen-client-urls http://${ALL_PASS_ADDR}:2379
