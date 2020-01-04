package parser

import (
	"errors"
	"logsvc/proto/model"
)

type SyslogParser struct {}

func (parser *SyslogParser) Marshal(request model.LogRequest) (string, error) {
	return "", errors.New("unsupported")
}

func (parser *SyslogParser) Unmarshal(text string, request *model.LogRequest) error {
	return nil
}
