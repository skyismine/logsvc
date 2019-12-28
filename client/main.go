package main

import (
	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"
	"logsvc/proto/model"
	"logsvc/proto/rpcapi"
	"time"
)

func main() {
	reg := etcdv3.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"192.168.3.23:2279",
			"192.168.3.23:2280",
			"192.168.3.23:2281",
		}
	})

	service := micro.NewService(micro.Registry(reg))
	service.Init()
	logsvcclient := rpcapi.NewLogClient("cb.srv.log", service.Client())
	rsp, err := logsvcclient.Info(context.Background(), &model.LogRequest{App: "LogClient", Tag: "Client", Msg: "this is a log message", Ctime: time.Now().Format(time.RFC3339Nano)})
	if err != nil {
		log.Fatalln("call srv error", err)
	}
	log.Println("call srv response", rsp)
}
