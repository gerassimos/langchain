#!/bin/bash
export OPENAI_API_KEY=$(cat ../openai-api-key.txt)
go run main.go cmd_executor.go