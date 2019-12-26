#!/usr/bin/env bash

CURPATH=$(pwd)
PROTO_ROOT=$CURPATH/proto

# --proto_path 和 -I 相同用来指定proto文件中import指令的搜索路径
# go proto生成的文件存储在 --go_out 指定的路径下
PROTO_MODEL=$PROTO_ROOT/model
protoc --proto_path="$PROTO_MODEL" --go_out=. "$PROTO_MODEL"/*.proto
cp logsvc/proto/model/*.go proto/model/
rm -rf logsvc

# 必须通过 --go_out=plugins=micro 的方式指定使用 micro 生成服务接口
# micro proto生成的文件存储在 --go_out 指定的路径/package/ 目录下, package为proto文件中package指令指定的值
PROTO_RPCAPI=$PROTO_ROOT/rpcapi
protoc --proto_path="$PROTO_ROOT":"$PROTO_RPCAPI" --go_out=plugins=micro:"$PROTO_ROOT" "$PROTO_RPCAPI"/*.proto
