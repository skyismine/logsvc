#!/usr/bin/env bash

# ./etcdctl.sh etcd-node-0 --help
# ./etcdctl.sh etcd-node-0 --endpoints=[192.168.3.23:2279] member list
DOCKERNAME=${1}
# 跳过第一个参数
shift 1
# 组合剩余参数
PARAM="$*"
docker exec ${DOCKERNAME} /bin/sh -c "export ETCDCTL_API=3 && /usr/local/bin/etcdctl ${PARAM}"
