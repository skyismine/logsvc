#!/usr/bin/env bash

# 拉取镜像
docker pull kibana:7.5.1
# 停止运行已经存在的容器
docker stop kibana
# 删除已经存在的容器
docker container rm kibana
#docker run -itd --net=host --name mongo -v /opt/data/mongodb:/data/db mongo --auth
docker run -itd --net=host --name kibana -e "ELASTICSEARCH_URL=http://192.168.3.23:9200" -e "ELASTICSEARCH_HOSTS=http://192.168.3.23:9200" -e "SERVER_HOST=192.168.3.23" kibana:7.5.1
