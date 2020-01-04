package parser

import (
	"errors"
	"fmt"
	"logsvc/proto/model"
	"strings"
	"time"
)

type GostdlogParser struct {}

func (parser *GostdlogParser) Marshal(request model.LogRequest) (string, error) {
	return "", errors.New("unsupported")
}

func (parser *GostdlogParser) Unmarshal(text string, request *model.LogRequest) error {
	request.Msg = text
	request.Level = Info
	txts := strings.Split(text, " ")
	timestr := fmt.Sprintf("%sT%s", strings.Replace(txts[0], "/", "-", -1), txts[1])
	t, err := time.Parse("2006-01-02T15:04:05", timestr)
	if err != nil {
		return errors.New(fmt.Sprintf("time: %s parse error: %s", timestr, err))
	}
	request.Ctime = t.Format(time.RFC3339)
	return nil
}
