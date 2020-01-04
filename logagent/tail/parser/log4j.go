package parser

import (
	"errors"
	"logsvc/proto/model"
)

type Log4jlogParser struct {}

func (parser *Log4jlogParser) Marshal(request model.LogRequest) (string, error) {
	return "", errors.New("unsupported")
}

func (parser *Log4jlogParser) Unmarshal(text string, request *model.LogRequest) error {
	return nil
}
