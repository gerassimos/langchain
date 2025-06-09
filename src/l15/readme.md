# overview
go mod init qs1
go get -u github.com/tmc/langchaingo
go get -u github.com/tmc/langchaingo/llms
go get -u github.com/tmc/langchaingo/chains
go get -u github.com/tmc/langchaingo/prompts

export OPENAI_API_KEY=$(cat open_api_key.txt)

Reference:
https://github.com/tmc/langchaingo/blob/main/examples/llm-chain-example/llm_chain.go 