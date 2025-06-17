#!/bin/bash

set -euo pipefail
module_name=${1:-l65_test_border_msg}
go mod init "${module_name}"
go get -u "github.com/charmbracelet/lipgloss"