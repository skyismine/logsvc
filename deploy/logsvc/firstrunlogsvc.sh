#!/usr/bin/env bash

# 运行 logproxy
nohup ./logproxy > proxy.log 2>&1 &

# 运行 logtail
nohup ./tail --log_file /work/CloudBox/logsvc/Bin/screen.log --log_app scrsvc --log_type gostd --log_seek 2 > tail.log 2>&1 &

# 运行 logsearcher
nohup ./logsearcher --searcher_es_domain="http://192.168.3.26:9200" > searcher.log 2>&1 &

# 运行 loganalysis
nohup ./loganalysis --analysis_es_domain="http://192.168.3.26:9200" > analysis.log 2>&1 &
