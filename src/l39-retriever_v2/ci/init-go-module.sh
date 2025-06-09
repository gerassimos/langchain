#!/bin/bash
## The selected line set -euo pipefail is a common safety feature in shell scripts. Here's what each option does:
## -e: Causes the script to exit immediately if any command exits with a non-zero status.
## -u: Treats unset variables as an error and exits the script.
## -o pipefail: Ensures that the script exits with a failure status if any command in a pipeline fails.
set -euo pipefail
module_name=${1:-ci}
go mod init "${module_name}"
go get -u github.com/tmc/langchaingo

go get github.com/tmc/langchaingo/documentloaders@v0.1.13
go get github.com/tmc/langchaingo/vectorstores/chroma@v0.1.13