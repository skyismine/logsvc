#!/usr/bin/env bash

go build -o Bin/client logsvc/logagent/example

go build -o Bin/tail logsvc/logagent/tail

go build -o Bin/proxy logsvc/logproxy

go build -o Bin/searcher logsvc/logsearcher
