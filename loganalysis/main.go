package main

import (
	"CommonUtil/src/GYGUtils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/micro/cli"
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
			logs.Error("json.Unmarshal error", err)
			continue
		} else if pmsg.App != "screen" {
			logs.Info("not screen msg ignore")
			continue
		} else if strings.Contains(pmsg.Msg, "NotificationHeartbeat") || strings.Contains(pmsg.Msg, "NotificationDevScreenHandler") {
			elements := strings.Split(pmsg.Msg, " ")
			record := msgmap[elements[2]]
			svdmsg, ok := record[elements[4]]
			record[elements[4]] = msg.Value
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
				buf := bytes.NewBuffer(svdmsg)
				buf.WriteString("\n")
				buf.Write(msg.Value)
				if elements[2] == "NotificationHeartbeat" && past > 5  {
					_, _ = hbfile.Write(buf.Bytes())
				} else if elements[2] == "NotificationHeartbeat" && past > 60 {
					_, _ = dshfile.Write(buf.Bytes())
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
			logs.Error("json.Unmarshal error", err)
			continue
		} else if pmsg.App != "scrsvc" {
			logs.Info("not screen msg ignore")
			continue
		} else if strings.Contains(pmsg.Msg, "NotificationHeartbeat") || strings.Contains(pmsg.Msg, "NotificationDevScreenHandler") {
			elements := strings.Split(pmsg.Msg, " ")
			record := msgmap[elements[2]]
			svdmsg, ok := record[elements[4]]
			record[elements[4]] = msg
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
				buf := bytes.NewBuffer(svdmsg)
				buf.WriteString("\n")
				buf.Write(msg)
				if elements[2] == "NotificationHeartbeat" && past > 5  {
					_, _ = hbfile.Write(buf.Bytes())
				} else if elements[2] == "NotificationHeartbeat" && past > 60 {
					_, _ = dshfile.Write(buf.Bytes())
				}
			}
		} else {
			logs.Info("others screen msg ignore")
		}
	}
}

var hbfile *os.File
var dshfile *os.File

func main() {
	var err error
	var consumerdomain, consumertype string
	app := cli.NewApp()
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "analysis_consumer_domain",
			Usage: "consumer domain",
			Value: "web.njnjdjc.com:29000",
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
	_ = logs.SetLogger(logs.AdapterConsole)

	hbfile, err = os.Create("hbfile.log")
	if err != nil {
		logs.Error("NotificationHeartbeat log file create error", err)
		return
	}
	defer func() { _ = hbfile.Close() }()
	dshfile, err := os.Create("dshfile.log")
	if err != nil {
		logs.Error("NotificationDevScreenHandler log file create error", err)
		return
	}
	defer func() { _ = dshfile.Close() }()

	if consumertype == "kafka" {
		consumerKafka(consumerdomain)
	} else {
		consumerNanomsg(consumerdomain)
	}
}