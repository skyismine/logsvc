# logsvc
日志服务,支持多种输入源,分布式日志存储,支持按时间段、Level、或关键字搜索所有日志
日志方案：使用 beggo log 库将log存储在文件中,使用 logagent/tail 进行log读取并发送到kafka, logsearch和其他分析插件通过consumer kafka的消息来进行log处理

```
-----------        ---------------        ------------
|beego log|        |logagent SDK |        |其他log采集|
-----------        ---------------        ------------
     |                      |                   /
     |                      |                  /
     |                      |                 /
---------------             |                /
|logagent tail|             |               /
---------------             |              /
     \                      |             /
        \                   |            /
            \               |           /
                    ---------------
                    |    kafka    |
                    ---------------
                    /           \
                   /             \
                  /               \
                 /                 \
        -------------        ------------------- 
        |logsearcher|        |log backend dael |
        -------------        -------------------
```

logagent: 客户端日志采集工具或SDK

logproxy: 日志代理服务端,logagent将日志发送给logproxy

logsearcher: 日志查询,订阅kafka消息后使用ElasticSearch索引

kibana: http://192.168.3.23:5601/

# 使用方法
```bash
cd deploy/elasticsearch
./firstrunes.sh
```
