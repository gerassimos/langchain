package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	// We can construct an LLMChain from a PromptTemplate and an LLM.
	llm, err := openai.New()
	if err != nil {
		return err
	}
	ctx := context.Background()

	reader := bufio.NewReader(os.Stdin)

	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewHumanMessagePromptTemplate(
			`{{.content}}`,
			[]string{"content"},
		),
	})

	llmChain := chains.NewLLMChain(llm, prompt)

	// out, ok := outputValues[llmChain.OutputKey].(string)
	// if !ok {
	// 	return fmt.Errorf("invalid chain return")
	// }
	// fmt.Println(out)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		outputValues, err := chains.Call(ctx, llmChain, map[string]any{
			"content": input,
		})
		if err != nil {
			return err
		}
		//loop over outputValues
		// for key, value := range outputValues {
		// 	fmt.Printf("key: %s: value: %s\n", key, value)
		// }

		out, ok := outputValues[llmChain.OutputKey].(string)
		if !ok {
			return fmt.Errorf("invalid chain return")
		}
		fmt.Println(out)
	}
}
