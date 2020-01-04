package parser

import (
	"errors"
	"logsvc/proto/model"
)

type NginxlogParser struct {}

func (parser *NginxlogParser) Marshal(request model.LogRequest) (string, error) {
	return "", errors.New("unsupported")
}

func (parser *NginxlogParser) Unmarshal(text string, request *model.LogRequest) error {
	return nil
}
