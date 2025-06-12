#!/bin/bash
export OPENAI_API_KEY=$(cat ../openai-api-key.txt)
export CHROMA_URL="http://localhost:8000"
go run .