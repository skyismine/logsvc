#!/usr/bin/env bash

# 拉取镜像
docker pull elasticsearch:7.5.1
# 停止运行已经存在的容器
docker stop elasticsearch
# 删除已经存在的容器
docker container rm elasticsearch
#docker run -itd --net=host --name mongo -v /opt/data/mongodb:/data/db mongo --auth
docker run -itd --net=host --name elasticsearch -e "discovery.type=single-node" elasticsearch:7.5.1
