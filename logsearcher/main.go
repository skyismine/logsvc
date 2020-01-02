package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"gopkg.in/olivere/elastic.v5"
	"log"
	"strings"
)

var esclient *elastic.Client

func getMsgFromPC(partitionConsumer sarama.PartitionConsumer) {
	defer partitionConsumer.AsyncClose()
	for msg := range partitionConsumer.Messages(){
		log.Println(fmt.Sprintf("partition:%d Offset:%d Key:%s Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value)))
		rsp, err := esclient.Index().Index("logs").Type("logs").BodyString(string(msg.Value)).Do(context.Background())
		log.Println("esclient bodystring rsp:", rsp, "error:", err)
	}
}

func main() {
	var err error
	esclient, err = elastic.NewClient(elastic.SetSniff(false),elastic.SetURL("http://192.168.3.23:9200/"))
	if err != nil{
		log.Fatalln("connect es error", err)
		return
	}
	consumer,err := sarama.NewConsumer(strings.Split("192.168.3.23:9092",","),nil)
	if err != nil{
		log.Fatalln("failed to start consumer:", err)
		return
	}
	defer func() { _ = consumer.Close() }()
	partitionList,err := consumer.Partitions("logs")
	if err != nil {
		log.Fatalln("Failed to get the list of partitions:",err)
		return
	}
	log.Println(partitionList)
	for partition := range partitionList {
		pc,err := consumer.ConsumePartition("logs",int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Fatalln(fmt.Sprintf("failed to start consumer for partition %d:%s\n", partition, err))
		}
		go getMsgFromPC(pc)
	}

	select {}
}