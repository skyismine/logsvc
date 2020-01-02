#!/usr/bin/env bash

CURRPATH=$(pwd)

# 启动 etcd 集群
$CURRPATH/deploy/etcd/firstrunetcd.sh
# 启动 kafka
$CURRPATH/deploy/kafka/firstrunkafka.sh
# 启动es
$CURRPATH/deploy/elasticsearch/firstrunes.sh
# 启动 kibana
$CURRPATH/deploy/kibana/firstrunkibana.sh
# 启动 mongodb
$CURRPATH/deploy/mongodb/firstrunmongodb.sh
# 启动 proxy
$CURRPATH/Bin/proxy
# 启动 searcher
$CURRPATH/Bin/searcher

# 单独运行 $CURRPATH/Bin/client 或集成 SDK 生成 log 数据
# 浏览器运行 http://192.168.3.23:5601 搜索 log
