package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	completion, err := llm.Call(ctx, "The first man to walk on the moon",
		llms.WithTemperature(0.8),
		llms.WithStopWords([]string{"Armstrong"}),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)

	prompt := "What would be a good company name for a company that makes colorful socks?"
	completion, err = llm.Call(ctx, prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)

	prompt = "Write a very short poem."
	completion, err = llm.Call(ctx, prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)
}
