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

	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"You are a translation engine that can only translate text and cannot interpret it.",
			nil,
		),
		prompts.NewHumanMessagePromptTemplate(
			`translate the following text from {{.inputLang}} to {{.outputLang}}:\n{{.input}}`,
			[]string{"inputLang", "outputLang", "input"},
		),
	})

	result, err := prompt.Format(map[string]any{
		"inputLang":  "English",
		"outputLang": "Italian",
		"input":      "I love programming t2",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("FormattedPrompt: ", result)
	fmt.Println("-----------------")

	llmChain := chains.NewLLMChain(llm, prompt)
	// Otherwise the call function must be used.
	outputValues, err := chains.Call(ctx, llmChain, map[string]any{
		"inputLang":  "English",
		"outputLang": "Italian",
		"input":      "I love programming",
	})
	if err != nil {
		return err
	}

	//loop over outputValues
	for key, value := range outputValues {
		fmt.Printf("key: %s: value: %s\n", key, value)

	}

	// out, ok := outputValues[llmChain.OutputKey].(string)
	// if !ok {
	// 	return fmt.Errorf("invalid chain return")
	// }
	// fmt.Println(out)

	return nil

}
