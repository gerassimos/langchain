package main

import (
	"context"
	"fmt"
	"log"
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

	prompt1 := prompts.NewPromptTemplate(
		"Write a very short {{.language}} function that will {{.task}}",
		[]string{"language", "task"},
	)
	prompt2 := prompts.NewPromptTemplate(
		"Write a test for the following {{.language}} code:\n {{.code}}",
		[]string{"language", "code"},
	)

	chain1 := chains.NewLLMChain(llm, prompt1)
	chain1.OutputKey = "code"
	chain2 := chains.NewLLMChain(llm, prompt2)
	chain2.OutputKey = "test"

	sequentialChain, err := chains.NewSequentialChain([]chains.Chain{chain1, chain2},
		[]string{"language", "task"}, // input keys
		[]string{"code", "test"})     // output keys
	if err != nil {
		log.Fatal(err)
	}

	outputValues, err := sequentialChain.Call(ctx, map[string]any{
		"language": "Java",
		"task":     "return a list of numbers",
	})
	if err != nil {
		return err
	}

	// outputValues, err := chains.Call(ctx, sequentialChain, map[string]any{
	// 	"language":  "Java",
	// 	"task":  "return a list of numbers",
	// })
	// if err != nil {
	// 	return err
	// }

	fmt.Println(outputValues)

	return nil
}
