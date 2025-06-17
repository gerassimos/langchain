#!/bin/bash

set -euo pipefail
module_name=${1:-l65v2}
go mod init "${module_name}"
go get -u github.com/tmc/langchaingo/llms