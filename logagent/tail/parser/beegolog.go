package parser

import (
	"errors"
	"fmt"
	"logsvc/proto/model"
	"strings"
	"time"
)

type BegoologParser struct {}

func (parser *BegoologParser) Marshal(request model.LogRequest) (string, error) {
	return "", errors.New("unsupported")
}

func (parser *BegoologParser) Unmarshal(text string, request *model.LogRequest) error {
	request.Msg = text
	txts := strings.Split(text, " ")
	timestr := fmt.Sprintf("%s %s", txts[0], txts[1])
	t, err := time.Parse("2006/01/02 15:04:05.999", timestr)
	if err != nil {
		return errors.New(fmt.Sprintf("time: %s parse error: %s", timestr, err))
	}
	request.Ctime = t.Format(time.RFC3339Nano)
	if txts[2] == "[I]" {
		request.Level = Info
	} else if txts[2] == "[W]" {
		request.Level = Warn
	} else if txts[2] == "[E]" {
		request.Level = Error
	} else if txts[2] == "[D]" {
		request.Level = Debug
	}
	return nil
}
