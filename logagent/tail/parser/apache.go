package parser

import (
	"errors"
	"logsvc/proto/model"
)

type ApachelogParser struct {}

func (parser *ApachelogParser) Marshal(request model.LogRequest) (string, error) {
	return "", errors.New("unsupported")
}

func (parser *ApachelogParser) Unmarshal(text string, request *model.LogRequest) error {
	return nil
}
