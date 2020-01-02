#!/usr/bin/env bash

go build -o Bin/client logsvc/logagent

go build -o Bin/proxy logsvc/logproxy

go build -o Bin/searcher logsvc/logsearcher
