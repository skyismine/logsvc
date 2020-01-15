package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
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

// usage: ./proxy --proxy_kafka_domain 192.168.3.23:9092
func main() {
	_ = logs.SetLogger(logs.AdapterFile, `{"filename": "/var/log/logsvc/logproxy.log"}`)

	var storedomain string

	svcflags := []cli.Flag{
		cli.StringFlag{
			Name:   	 "proxy_storage_domain",
			Usage:  	 "kafka domain where proxy used",
			Value:		 "tcp://*:29000",
			Destination: &storedomain,
		},
	}

	//使用etcd做服务发现
	reg := etcdv3.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"192.168.3.23:2379",
		}
		etcdv3.Auth("root", "Erika6014")(options)
	})
	service := micro.NewService(micro.Name("cb.srv.log"), micro.Registry(reg))
	service.Options().Cmd.App().Flags = append(service.Options().Cmd.App().Flags, svcflags...)
	service.Init()
	rpcapi.RegisterLoggerHandler(service.Server(), new(LogSvc))

	store = storage.NewStorageStorageNanomsg(storedomain)

	if err := service.Run(); err != nil {
		logs.Error("server log run error", err)
	}
}
