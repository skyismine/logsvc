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

const (
	Trace = "trace"
	Debug = "debug"
	Info = "info"
	Warn = "warn"
	Error = "error"
	Fatal = "fatal"
	Panic = "panic"
)

var store storage.IFStorage

func logger(in *model.LogRequest, out *model.LogResponse, level string) {
	data := new(storage.Logmsg)
	data.Host = in.Host
	data.App = in.App
	data.Level = level
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
}

type LogSvc struct {}

func (h *LogSvc) Trace(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logger(in, out, Trace)
	return nil
}

func (h *LogSvc) Debug(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logger(in, out, Debug)
	return nil
}

func (h *LogSvc) Info(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logger(in, out, Info)
	return nil
}

func (h *LogSvc) Warn(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logger(in, out, Warn)
	return nil
}

func (h *LogSvc) Error(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logger(in, out, Error)
	return nil
}

func (h *LogSvc) Fatal(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logger(in, out, Fatal)
	return nil
}

func (h *LogSvc) Panic(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logger(in, out, Panic)
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
	rpcapi.RegisterLogHandler(service.Server(), new(LogSvc))

	if err := service.Run(); err != nil {
		log.Fatalln("server log run error", err)
	}
}
