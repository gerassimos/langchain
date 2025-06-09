package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		return err
	}

	// opts := []openai.Option{
	// 	openai.WithModel("gpt-3.5-turbo-0125"),
	// 	openai.WithEmbeddingModel("text-embedding-3-large"),
	// }
	// llm, err := openai.New(opts...)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	embedings, err := llm.CreateEmbedding(ctx, []string{"ola"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(embedings)
	//len(embedings):  1536 => with default model
	//len(embedings):  3072 => with gpt-3.5-turbo-0125 model - EmbeddingModel "text-embedding-3-large"
	fmt.Println("len(embedings): ", len(embedings[0]))

	return nil
}
