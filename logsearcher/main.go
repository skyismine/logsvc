package main

import (
	"CommonUtil/src/GYGUtils"
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/micro/cli"
	"gopkg.in/olivere/elastic.v5"
	"os"
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

func consumerKafka(domain string) {
	consumer,err := sarama.NewConsumer(strings.Split(domain,","),nil)
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
}

func consumerNanomsg(domain string) {
	node := GYGUtils.SubNode(domain, "", nil)
	for {
		msg, _, err := GYGUtils.GSocketRecv(node)
		if err != nil {
			logs.Error("consumerNanomsg sub RecvMessage error", err.Error())
			continue
		}
		rsp, err := esclient.Index().Index("logs").Type("logs").BodyString(string(msg)).Do(context.Background())
		logs.Info("esclient bodystring rsp:", rsp, "error:", err)
	}
}

func main() {
	var err error
	var esdomain, consumerdomain, consumertype string
	app := cli.NewApp()
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "searcher_es_domain",
			Usage: "elastic domain",
			Destination: &esdomain,
		},
		cli.StringFlag{
			Name: "searcher_consumer_domain",
			Usage: "consumer domain",
			Value: "web.njnjdjc.com:29000",
			Destination: &consumerdomain,
		},
		cli.StringFlag{
			Name: "searcher_consumer_type",
			Value: "nanomsg",
			Usage: "consumer domain",
			Destination: &consumertype,
		},
	}
	_ = app.Run(os.Args)

	_ = logs.SetLogger(logs.AdapterConsole)
	esclient, err = elastic.NewClient(elastic.SetSniff(false),elastic.SetURL(esdomain))
	if err != nil{
		logs.Error("connect es error", err)
		return
	}

	if consumertype == "kafka" {
		consumerKafka(consumerdomain)
	} else {
		consumerNanomsg(consumerdomain)
	}
}