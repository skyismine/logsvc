#!/usr/bin/env bash

# 运行 logproxy
nohup ./logproxy &

# 运行 loganalysis
nohup ./loganalysis --analysis_es_domain="http://172.31.235.243:9200" --analysis_consumer_domain="tcp://172.31.235.243:29000" &

# 运行 logsearcher
nohup ./logsearcher --searcher_es_domain="http://172.31.235.243:9200" --searcher_consumer_domain="tcp://172.31.235.243:29000" &

# 运行 logtail
nohup ./logtail --log_file /opt/data/log/cloudbox/scrsvc.log --log_app scrsvc --log_type beego --log_seek 0 &
