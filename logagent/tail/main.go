package main

import (
	"context"
	"github.com/hpcloud/tail"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"
	"logsvc/logagent/tail/parser"
	"logsvc/proto/model"
	"logsvc/proto/rpcapi"
	"os"
	"time"
)

// usage: ./tail --log_file /work/CloudBox/logsvc/Bin/screen.log --log_app scrsvc --log_type gostd --log_seek 2
func main() {
	var filename, app, logtype string
	var seek int
	reg := etcdv3.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"192.168.3.23:2279",
			"192.168.3.23:2280",
			"192.168.3.23:2281",
		}
	})

	service := micro.NewService(micro.Registry(reg))
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
	service.Options().Cmd.App().Flags = append(service.Options().Cmd.App().Flags, svcflags...)
	service.Init()
	logsvcclient := rpcapi.NewLoggerClient("cb.srv.log", service.Client())
	log.Println("tail log_file:", filename, ", log_app:", app, ", log_type:", logtype)
	tails, err := tail.TailFile(filename, tail.Config {
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{Offset: 0, Whence: seek},
		MustExist: false,
		Poll: true,
	})
	if err !=nil{
		log.Fatalln("tail file err:",err)
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
			log.Printf("tail file close reopen,filenam:%s\n", filename)
			time.Sleep(100*time.Millisecond)
			continue
		}
		line++
		log.Println("tail line", line)
		if logparser, ok := parser.PManager[logtype]; ok {
			var logmsg model.LogRequest
			err := logparser.Unmarshal(msg.Text, &logmsg)
			if err != nil {
				log.Println("tail parser.Unmarshal error", err)
				continue
			}
			logmsg.App = app
			logmsg.Host = host
			_, err = logsvcclient.Log(context.Background(), &logmsg)
			if err != nil {
				log.Fatalln("call srv error", err)
			}
		} else {
			log.Println("tail no log parser for app", app)
		}
	}
}