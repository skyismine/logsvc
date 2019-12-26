package main

import (
	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/sirupsen/logrus"
	"log"
	"logsvc/proto/model"
	"logsvc/proto/rpcapi"
)

func logfunc(in *model.LogRequest, out *model.LogResponse, level logrus.Level) {
	fields := make(logrus.Fields)
	for key, value := range in.Field {
		fields[key] = value
	}
	logrus.WithFields(logrus.Fields{"app": in.App, "tag": in.Tag}).WithFields(fields).Log(level, in.Msg)
	out.Msg = "logsvc success"
}

type LogSvc struct {}

func (h *LogSvc) Trace(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logfunc(in, out, logrus.TraceLevel)
	return nil
}

func (h *LogSvc) Debug(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logfunc(in, out, logrus.DebugLevel)
	return nil
}

func (h *LogSvc) Info(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logfunc(in, out, logrus.InfoLevel)
	return nil
}

func (h *LogSvc) Warn(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logfunc(in, out, logrus.WarnLevel)
	return nil
}

func (h *LogSvc) Error(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logfunc(in, out, logrus.ErrorLevel)
	return nil
}

func (h *LogSvc) Fatal(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logfunc(in, out, logrus.FatalLevel)
	return nil
}

func (h *LogSvc) Panic(ctx context.Context, in *model.LogRequest, out *model.LogResponse) error {
	logfunc(in, out, logrus.PanicLevel)
	return nil
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

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
