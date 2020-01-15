package storage

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type StorageKafka struct {
	kafkaclient sarama.SyncProducer
}

func NewStorageKafka(addr string) *StorageKafka {
	var err error
	kafka := new(StorageKafka)
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	kafka.kafkaclient, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil{
		logs.Error("producer close err:", err)
		return nil
	}
	return kafka
}

func (store *StorageKafka) Save(msg *Logmsg) error {
	pmsg := &sarama.ProducerMessage{}
	pmsg.Topic = "logs"
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	pmsg.Value = sarama.StringEncoder(string(data))
	_, _, err = store.kafkaclient.SendMessage(pmsg)
	if err != nil{
		return nil
	}
	return nil
}

func (store *StorageKafka) Close() {
	_ = store.kafkaclient.Close()
}
