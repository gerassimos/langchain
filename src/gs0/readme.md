# overview
go mod init qs1
go get -u github.com/tmc/langchaingo
go get -u github.com/tmc/langchaingo/llms

export OPENAI_API_KEY=$(cat open_api_key.txt)