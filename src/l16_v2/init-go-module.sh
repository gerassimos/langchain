#!/bin/bash

module_name=${1:-gs0}
go mod init "${module_name}"
go get -u github.com/tmc/langchaingo
go get -u github.com/tmc/langchaingo/llms
go get -u github.com/tmc/langchaingo/chains
go get -u github.com/tmc/langchaingo/prompts

## Generate main.go
#cat <<EOF > main.go
#package main
#import (
#  "fmt"
#  "os"
#
#  "github.com/tmc/langchaingo/llms"
#  "github.com/tmc/langchaingo/prompts"
#)
#
#func main() {
#  fmt.Println("Hello, LangChain Go!")
#
#}
#EOF