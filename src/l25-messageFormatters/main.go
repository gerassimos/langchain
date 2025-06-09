package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"

	"github.com/tmc/langchaingo/httputil"

	"github.com/tmc/langchaingo/prompts"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var opts []openai.Option
	//use a custom HTTP client to log requests and responses.
	opts = append(opts, openai.WithHTTPClient(httputil.DebugHTTPClient))
	// We can construct an LLMChain from a PromptTemplate and an LLM.
	llm, err := openai.New(opts...)
	if err != nil {
		return err
	}
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	messageFormatters := []prompts.MessageFormatter{prompts.NewSystemMessagePromptTemplate("", nil)}
	messageFormatters = append(messageFormatters, prompts.NewHumanMessagePromptTemplate("{{.content}}", []string{"content"}))
	messageFormatters = append(messageFormatters, prompts.MessagesPlaceholder{
		VariableName: "messages",
	})

	prompt := prompts.NewChatPromptTemplate(messageFormatters)

	// conversationBuffer := memory.NewConversationBuffer(memory.WithChatHistory(chatHistory))
	conversationBuffer := memory.NewConversationBuffer()
	// conversationBuffer.ReturnMessages = true
	fmt.Println("InputKey: ", conversationBuffer.InputKey, "OutputKey :", conversationBuffer.OutputKey, "MemoryKey: ", conversationBuffer.MemoryKey)
	fmt.Println("ReturnMessages", conversationBuffer.ReturnMessages)

	// llmChain := chains.NewConversation(llm, conversationBuffer)
	llmChain := chains.NewLLMChain(llm, prompt)
	llmChain.Memory = conversationBuffer

	// // prepare the db with some sample data
	// if err := prepare(ctx, db); err != nil {
	// 	return err
	// }

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		// out, err := chains.Run(ctx, llmChain, input)
		out, err := chains.Call(ctx, llmChain, map[string]any{
			"content":  input,
			"messages": []string{"Hello", "World"},
		})
		if err != nil {
			return err
		}

		fmt.Println(out)
	}
}
