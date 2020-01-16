package main

import (
	"CommonUtil/src/GYGUtils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/micro/cli"
	"gopkg.in/olivere/elastic.v5"
	"logsvc/logproxy/storage"
	"os"
	"strings"
	"time"
)

var msgmap map[string]map[string][]byte

func getMsgFromPC(partitionConsumer sarama.PartitionConsumer) {
	defer partitionConsumer.AsyncClose()
	for msg := range partitionConsumer.Messages(){
		logs.Info(fmt.Sprintf("partition:%d Offset:%d Key:%s Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value)))
		var pmsg storage.Logmsg
		err := json.Unmarshal(msg.Value, &pmsg)
		if err != nil {
			logs.Error("json.Unmarshal error", err, msg.Value)
			continue
		} else if pmsg.App != "screen" {
			logs.Info("not screen msg ignore")
			continue
		} else if strings.Contains(pmsg.Msg, "NotificationHeartbeat") || strings.Contains(pmsg.Msg, "NotificationDevScreenHandler") {
			elements := strings.Split(pmsg.Msg, " ")
			record, ok := msgmap[elements[4]]
			if !ok {
				logs.Info("It's not heartbeat or screen handler msg, elements[4]", pmsg.Msg)
				continue
			}
			svdmsg, ok := record[elements[6]]
			if !ok {
				logs.Info("It's not heartbeat or screen handler msg, elements[6]", pmsg.Msg)
				continue
			}
			record[elements[6]] = msg.Value
			if ok {
				var lastmsg storage.Logmsg
				err := json.Unmarshal(svdmsg, &lastmsg)
				if err != nil {
					logs.Error("lastmsg json.Unmarshal error", err, string(svdmsg))
					continue
				}
				lasttime, err := time.Parse(time.RFC3339, lastmsg.Ctime)
				if err != nil {
					logs.Error("lastmsg lasttime parse error", err, string(svdmsg))
					continue
				}
				curtime, err := time.Parse(time.RFC3339, pmsg.Ctime)
				if err != nil {
					logs.Error("lastmsg curtime parse error", err, string(svdmsg))
					continue
				}
				past := curtime.Unix()-lasttime.Unix()
				if (elements[2] == "NotificationHeartbeat" && past > 5) || (elements[2] == "NotificationHeartbeat" && past > 60) {
					rsp, err := esclient.Index().Index("lostmsg").Type("lostmsg").BodyString(string(svdmsg)).Do(context.Background())
					logs.Info("loganalysis", "esclient last bodystring rsp:", rsp, "error:", err)
					rsp, err = esclient.Index().Index("lostmsg").Type("lostmsg").BodyString(string(msg.Value)).Do(context.Background())
					logs.Info("loganalysis", "esclient current bodystring rsp:", rsp, "error:", err)
				}
			}
		} else {
			logs.Info("others screen msg ignore")
		}
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
		pc,err := consumer.ConsumePartition("logs", int32(partition), sarama.OffsetNewest)
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
		var pmsg storage.Logmsg
		err = json.Unmarshal(msg, &pmsg)
		if err != nil {
			logs.Error("json.Unmarshal error", err, msg)
			continue
		} else if pmsg.App != "scrsvc" {
			logs.Info("not screen msg ignore")
			continue
		} else if strings.Contains(pmsg.Msg, "NotificationHeartbeat") || strings.Contains(pmsg.Msg, "NotificationDevScreenHandler") {
			elements := strings.Split(pmsg.Msg, " ")
			logs.Info("elements", fmt.Sprintf("%#v", elements))
			record, ok := msgmap[elements[4]]
			if !ok {
				logs.Info("It's not heartbeat or screen handler msg, elements[4]", pmsg.Msg)
				continue
			}
			svdmsg, ok := record[elements[6]]
			if !ok {
				logs.Info("It's not heartbeat or screen handler msg, elements[6]", pmsg.Msg)
				continue
			}
			record[elements[6]] = msg
			if ok {
				var lastmsg storage.Logmsg
				err := json.Unmarshal(svdmsg, &lastmsg)
				if err != nil {
					logs.Error("lastmsg json.Unmarshal error", err, string(svdmsg))
					continue
				}
				lasttime, err := time.Parse(time.RFC3339, lastmsg.Ctime)
				if err != nil {
					logs.Error("lastmsg lasttime parse error", err, string(svdmsg))
					continue
				}
				curtime, err := time.Parse(time.RFC3339, pmsg.Ctime)
				if err != nil {
					logs.Error("lastmsg curtime parse error", err, string(svdmsg))
					continue
				}
				past := curtime.Unix()-lasttime.Unix()
				if (elements[2] == "NotificationHeartbeat" && past > 5) || (elements[2] == "NotificationHeartbeat" && past > 60) {
					rsp, err := esclient.Index().Index("lostmsg").Type("lostmsg").BodyString(string(svdmsg)).Do(context.Background())
					logs.Info("loganalysis", "esclient last bodystring rsp:", rsp, "error:", err)
					rsp, err = esclient.Index().Index("lostmsg").Type("lostmsg").BodyString(string(msg)).Do(context.Background())
					logs.Info("loganalysis", "esclient current bodystring rsp:", rsp, "error:", err)
				}
			}
		} else {
			logs.Info("others screen msg ignore")
		}
	}
}

var esclient *elastic.Client
// usage: ./analysis --analysis_es_domain="http://192.168.3.26:9200" --analysis_consumer_domain="tcp://web.njnjdjc.com:29000"
func main() {
	var err error
	var esdomain, consumerdomain, consumertype string
	app := cli.NewApp()
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "analysis_es_domain",
			Usage: "elastic domain",
			Destination: &esdomain,
		},
		cli.StringFlag{
			Name: "analysis_consumer_domain",
			Usage: "consumer domain",
			Value: "tcp://web.njnjdjc.com:29000",
			Destination: &consumerdomain,
		},
		cli.StringFlag{
			Name: "analysis_consumer_type",
			Value: "nanomsg",
			Usage: "consumer type",
			Destination: &consumertype,
		},
	}
	_ = app.Run(os.Args)

	msgmap = make(map[string]map[string][]byte)
	msgmap["NotificationHeartbeat"] = make(map[string][]byte)
	msgmap["NotificationDevScreenHandler"] = make(map[string][]byte)
	_ = logs.SetLogger(logs.AdapterFile, `{"filename": "/opt/data/log/cloudbox/loganalysis.log"}`)

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