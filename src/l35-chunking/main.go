package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/documentloaders"

	"github.com/tmc/langchaingo/textsplitter"
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
	// DefaultOptions
	//	Separators:    []string{"\n\n", "\n", " ", ""},
	// Set specific separators, discard default separators
	// split.Separators = []string{"\n"} // separators are the characters that the text is split on
	// split.Separators = []string{"#"} // separators are the characters that the text is split on
	docs, err := p.LoadAndSplit(ctx, split)

	if err != nil {
		return err
	}

	for _, doc := range docs {
		fmt.Println(doc.PageContent)
		fmt.Println("")
	}

	fmt.Println("len(docs):", len(docs))

	return nil
}
