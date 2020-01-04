package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"
	"logsvc/logproxy/storage"
	"logsvc/proto/model"
	"logsvc/proto/rpcapi"
	"time"
)

var store storage.IFStorage

type LogSvc struct {}

func (h *LogSvc) Log(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	data := new(storage.Logmsg)
	data.Host = in.Host
	data.App = in.App
	data.Level = in.Level
	data.Tag = in.Tag
	data.Msg = in.Msg
	data.Ctime = in.Ctime
	data.Stime = time.Now().Format(time.RFC3339Nano)
	err := store.Save(data)
	if err != nil {
		out.Msg = fmt.Sprintf("Error: %v", err)
	} else {
		out.Msg = "Success"
	}
	return nil
}

func main() {
	store = storage.NewStorageKafka("192.168.3.23")
	//使用etcd做服务发现
	reg := etcdv3.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"192.168.3.23:2279",
			"192.168.3.23:2280",
			"192.168.3.23:2281",
		}
	})

	service := micro.NewService(micro.Name("cb.srv.log"), micro.Registry(reg))
	service.Init()
	rpcapi.RegisterLoggerHandler(service.Server(), new(LogSvc))

	if err := service.Run(); err != nil {
		log.Fatalln("server log run error", err)
	}
}
