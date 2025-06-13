#!/bin/bash

set -euo pipefail
module_name="mrkl-agent-example"
go mod init "${module_name}"

go get -u github.com/tmc/langchaingo/llms