#!/usr/bin/env bash

# 停止运行已经存在的容器
docker stop mongo
# 删除已经存在的容器
docker container rm mongo
#docker run -itd --net=host --name mongo -v /opt/data/mongodb:/data/db mongo --auth
docker run -itd --net=host --name mongo -v /opt/data/mongodb:/data/db mongo

#$ docker exec -it mongo mongo admin
## 创建一个名为 admin，密码为 123456 的用户。
#>  db.createUser({ user:'admin',pwd:'123456',roles:[ { role:'userAdminAnyDatabase', db: 'admin'}]});
## 尝试使用上面创建的用户信息进行连接。
#> db.auth('admin', '123456')
## 以后命令行连接使用
#$ docker exec -it mongo mongo "mongodb://admin:nimda@localhost:27017"
