#!/usr/bin/env bash

go build -o Bin/client logsvc/logagent/example

go build -o Bin/logtail logsvc/logagent/tail

go build -o Bin/logproxy logsvc/logproxy

go build -o Bin/logsearcher logsvc/logsearcher

go build -o Bin/loganalysis logsvc/loganalysis
