package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/textsplitter"
	"log"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	// // We can construct an LLMChain from a PromptTemplate and an LLM.
	// llm, err := openai.New()
	// if err != nil {
	// 	return err
	// }
	ctx := context.Background()

	f, err := os.Open("../resources/facts.txt")
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}

	p := documentloaders.NewText(f)

	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 300   // size of the chunk is number of characters
	split.ChunkOverlap = 30 // overlap is the number of characters that the chunks overlap
	docs, err := p.LoadAndSplit(context.Background(), split)

	if err != nil {
		fmt.Println("Error loading document: ", err)
	}

	log.Println("Document loaded: ", len(docs))

	// docs, err := p.Load(ctx)
	// if err != nil {
	// 	return err
	// }

	fmt.Println("len(docs)", len(docs))

	for _, doc := range docs {
		fmt.Println("doc.PageContent: ", doc.PageContent)
		fmt.Println("doc.Metadata: ", doc.Metadata)
		fmt.Println("doc.Score: ", doc.Score)
	}

	return nil
}
