package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"gopkg.in/olivere/elastic.v5"
	"strings"
)

var esclient *elastic.Client

func getMsgFromPC(partitionConsumer sarama.PartitionConsumer) {
	defer partitionConsumer.AsyncClose()
	for msg := range partitionConsumer.Messages(){
		logs.Info(fmt.Sprintf("partition:%d Offset:%d Key:%s Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value)))
		rsp, err := esclient.Index().Index("logs").Type("logs").BodyString(string(msg.Value)).Do(context.Background())
		logs.Info("esclient bodystring rsp:", rsp, "error:", err)
	}
}

func main() {
	logs.Async(1e3)
	_ = logs.SetLogger(logs.AdapterFile, `{"filename": "/var/log/logsvc/logsearcher.log"}`)

	var err error
	esclient, err = elastic.NewClient(elastic.SetSniff(false),elastic.SetURL("http://192.168.3.23:9200/"))
	if err != nil{
		logs.Error("connect es error", err)
		return
	}
	consumer,err := sarama.NewConsumer(strings.Split("192.168.3.23:9092",","),nil)
	if err != nil{
		logs.Error("failed to start consumer:", err)
		return
	}
	defer func() { _ = consumer.Close() }()
	partitionList,err := consumer.Partitions("logs")
	if err != nil {
		logs.Error("Failed to get the list of partitions:",err)
		return
	}
	logs.Info(partitionList)
	for partition := range partitionList {
		pc,err := consumer.ConsumePartition("logs",int32(partition), sarama.OffsetNewest)
		if err != nil {
			logs.Error(fmt.Sprintf("failed to start consumer for partition %d:%s\n", partition, err))
			return
		}
		go getMsgFromPC(pc)
	}

	select {}
}