# logsvc
日志服务,支持多种输入源,分布式日志存储,支持按时间段、Level、或关键字搜索所有日志

logagent: 客户端日志采集工具或SDK

logproxy: 日志代理服务端,logagent将日志发送给logproxy

logsearcher: 日志查询,订阅kafka消息后使用ElasticSearch索引
