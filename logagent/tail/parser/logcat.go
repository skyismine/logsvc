package parser

import (
	"errors"
	"logsvc/proto/model"
)

type LogcatlogParser struct {}

func (parser *LogcatlogParser) Marshal(request model.LogRequest) (string, error) {
	return "", errors.New("unsupported")
}

func (parser *LogcatlogParser) Unmarshal(text string, request *model.LogRequest) error {
	return nil
}
