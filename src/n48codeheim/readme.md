# setup
go mod init llm_cli
go get -u github.com/spf13/cobra@latest
go get -u github.com/spf13/cobra@v1.8.1
go get -u github.com/tmc/langchaingo
go get -u github.com/tmc/langchaingo/llms

# cobra-cli
go install github.com/spf13/cobra-cli@latest
- cobra-cli init => will the structure of the app (main.go) and cmd folder
cobra-cli init
- create a subcommand called `chat`
cobra-cli add chat 

export OPENAI_API_KEY=$(cat open_api_key.txt)

# ref
https://github.com/tmc/langchaingo/tree/main/examples/openai-chat-example