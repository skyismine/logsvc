module logsvc

go 1.13

require (
	CommonUtil v0.0.0-00010101000000-000000000000
	github.com/DataDog/zstd v1.4.4 // indirect
	github.com/Shopify/sarama v1.24.1
	github.com/astaxie/beego v1.12.0
	github.com/golang/protobuf v1.3.2
	github.com/hpcloud/tail v1.0.0
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	github.com/tidwall/pretty v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.2.0
	golang.org/x/net v0.0.0-20191112182307-2180aed22343
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/olivere/elastic.v5 v5.0.82
)

replace CommonUtil => /work/CloudBox/CommonUtil
