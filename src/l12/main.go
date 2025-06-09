package main

import (
	"context"
	"fmt"
	"os"

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

	prompt := prompts.NewPromptTemplate(
		"Write a very short {{.language}} function that will {{.task}}",
		[]string{"language", "task"},
	)
	llmChain := chains.NewLLMChain(llm, prompt)

	//# Use the Call function from chains
	outputValues, err := chains.Call(ctx, llmChain, map[string]any{
		"language": "Java",
		"task":     "return a list of numbers",
	})
	if err != nil {
		return err
	}

	//# Calling from the chain directly
	// outputValues, err := llmChain.Call(ctx, map[string]any{
	// 	"language":  "Java",
	// 	"task":  "return a list of numbers",
	// })
	// if err != nil {
	// 	return err
	// }

	//At this point the outputValues is a Map containing one key-value pair
	//The key is [text] and
	//The value is the [generated text]

	//loop over outputValues
	// for key, value := range outputValues {
	//   fmt.Printf("key: %s: value: %s\n", key, value)

	// }

	out, ok := outputValues[llmChain.OutputKey].(string)
	if !ok {
		return fmt.Errorf("invalid chain return")
	}
	fmt.Println(out)

	return nil
}
