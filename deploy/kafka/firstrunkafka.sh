#!/usr/bin/env bash

docker stop zookeeper
docker container rm zookeeper

docker stop kafka
docker container rm kafka

docker-compose -f ./kafka-zk.yml up -d
