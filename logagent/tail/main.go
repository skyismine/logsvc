package main

import (
	"context"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"logsvc/logagent/tail/parser"
	"logsvc/proto/model"
	"logsvc/proto/rpcapi"
	"os"
	"time"
)

// usage: ./tail --log_file /work/CloudBox/logsvc/Bin/screen.log --log_app scrsvc --log_type gostd --log_seek 2
func main() {
	_ = logs.SetLogger(logs.AdapterConsole)

	var filename, app, logtype string
	var seek int

	svcflags := []cli.Flag{
		cli.StringFlag{
			Name:   	 "log_file",
			Usage:  	 "log file for tail",
			Destination: &filename,
		},
		cli.StringFlag{
			Name:   	 "log_app",
			Usage:  	 "log app identification",
			Destination: &app,
		},
		cli.StringFlag{
			Name:   	 "log_type",
			Usage:  	 "log type identification",
			Destination: &logtype,
		},
		cli.IntFlag{
			Name:        "log_seek",
			Usage:       "log file seek offset: 0 seek relative to the origin of the file, 1 seek relative to the current offset, 2 seek relative to the end",
			Destination: &seek,
		},
	}
	reg := etcdv3.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"192.168.3.23:2379",
		}
		etcdv3.Auth("root", "Erika6014")(options)
	})

	service := micro.NewService(micro.Registry(reg))
	service.Options().Cmd.App().Flags = append(service.Options().Cmd.App().Flags, svcflags...)
	service.Init()
	logsvcclient := rpcapi.NewLoggerClient("cb.srv.log", service.Client())
	logs.Info("tail log_file:", filename, ", log_app:", app, ", log_type:", logtype)
	tails, err := tail.TailFile(filename, tail.Config {
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{Offset: 0, Whence: seek},
		MustExist: false,
		Poll: true,
	})
	if err !=nil{
		logs.Error("tail file err:",err)
		return
	}
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	var msg *tail.Line
	var ok bool
	var line int
	for {
		msg, ok = <-tails.Lines
		if !ok {
			logs.Info("tail file close reopen,filenam:%s\n", filename)
			time.Sleep(100*time.Millisecond)
			continue
		}
		line++
		logs.Info("tail line", line)
		if logparser, ok := parser.PManager[logtype]; ok {
			var logmsg model.LogRequest
			err := logparser.Unmarshal(msg.Text, &logmsg)
			if err != nil {
				logs.Info("tail parser.Unmarshal error", err)
				continue
			}
			logmsg.App = app
			logmsg.Host = host
			_, err = logsvcclient.Log(context.Background(), &logmsg)
			if err != nil {
				logs.Error("call srv error", err)
				return
			}
		} else {
			logs.Info("tail no log parser for app", app)
		}
	}
}