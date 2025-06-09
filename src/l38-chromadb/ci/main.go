package main

import (
	"context"
	"fmt"
	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/chroma"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {

	ctx := context.Background()

	f, err := os.Open("../../resources/facts.txt")
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}

	p := documentloaders.NewText(f)

	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 300 // size of the chunk is number of characters
	// split.ChunkSize = 100   // size of the chunk is number of characters
	split.ChunkOverlap = 3 // overlap is the number of characters that the chunks overlap
	// DefaultOptions
	//	Separators:    []string{"\n\n", "\n", " ", ""},
	// Set specific separators, discard default separators
	// split.Separators = []string{"\n"} // separators are the characters that the text is split on
	// split.Separators = []string{"#"} // separators are the characters that the text is split on
	docs, err := p.LoadAndSplit(ctx, split)
	//update docs with metadata
	// for i := range docs {
	// 	docs[i].Metadata = map[string]interface{}{
	// 		"source": "facts.txt",
	// 	}
	// }
	fmt.Println("=====================================")
	fmt.Println(docs)
	for _, doc := range docs {
		fmt.Println()
		fmt.Println(doc.PageContent)
		fmt.Println()
	}

	// for loop over the docs

	fmt.Println("=====================================")
	fmt.Println("len(docs):", len(docs))

	if err != nil {
		return err
	}

	// Create a new Chroma vector store (db client).
	// storeNs:=uuid.New().String()
	storeNs := "chroma-ns-1"
	fmt.Println("storeNs:", storeNs)
	store, err := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace(storeNs),
	)
	if err != nil {
		return err
	}
	// Add documents to the vector store.
	// Will actually calculate embeddings for the documents by calling the OpenAI API.
	ids, err := store.AddDocuments(context.Background(), docs)
	if err != nil {
		return err
	}
	fmt.Println(ids)
	fmt.Println("End!!")
	return nil
}
