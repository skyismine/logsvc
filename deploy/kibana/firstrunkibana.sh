#!/usr/bin/env bash

# 拉取镜像
docker pull kibana:7.5.1
# 停止运行已经存在的容器
docker stop kibana
# 删除已经存在的容器
docker container rm kibana
#docker run -itd --net=host --name mongo -v /opt/data/mongodb:/data/db mongo --auth
TAURUS_HOST=`ifconfig -a|grep inet|grep -v inet6|grep -v 172.17.*|grep -v 172.18.*|grep -v 127.0.0.1|awk '{print $2}'|tr -d "addr:"`
docker run -itd --net=host --name kibana --add-host "cloudbox.591ota.com:${TAURUS_HOST}" -e TZ="Asia/Shanghai" -e "ELASTICSEARCH_URL=http://cloudbox.591ota.com:9200" -e "ELASTICSEARCH_HOSTS=http://cloudbox.591ota.com:9200" -e "SERVER_HOST=0.0.0.0" kibana:7.5.1
