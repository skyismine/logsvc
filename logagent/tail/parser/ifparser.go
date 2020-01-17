package parser

import "logsvc/proto/model"

const (
	Trace = "trace"
	Debug = "debug"
	Info = "info"
	Warn = "warn"
	Error = "error"
	Fatal = "fatal"
	Panic = "panic"
)

type IFParser interface {
	Marshal(request model.LogRequest) (string, error)
	Unmarshal(text string, request *model.LogRequest) error
}

var PManager map[string]IFParser

func RegisterParser(app string, parser IFParser) {
	PManager[app] = parser
}

func UnregisterParser(app string) {
	delete(PManager, app)
}

func init() {
	PManager = make(map[string]IFParser)
	RegisterParser("syslog", &SyslogParser{})
	RegisterParser("nginx", &NginxlogParser{})
	RegisterParser("logcat", &LogcatlogParser{})
	RegisterParser("apache", &ApachelogParser{})
	RegisterParser("gostd", &GostdlogParser{})
	RegisterParser("beego", &BegoologParser{})
	RegisterParser("log4j", &Log4jlogParser{})
}