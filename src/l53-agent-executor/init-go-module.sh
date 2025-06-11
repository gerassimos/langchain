#!/bin/bash

set -euo pipefail
module_name=${1:-l53-00}
go mod init "${module_name}"
#go get -u github.com/tmc/langchaingo
#go get -u github.com/tmc/langchaingo/llms
#go get -u github.com/tmc/langchaingo/chains
#go get -u github.com/tmc/langchaingo/prompts
go get -u github.com/mattn/go-sqlite3
#go get -u github.com/tmc/langchaingo/httputil