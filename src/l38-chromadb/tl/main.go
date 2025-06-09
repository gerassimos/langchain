package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
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
	split.ChunkSize = 300  // size of the chunk is number of characters
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
	ids, err := store.AddDocuments(context.Background(), docs)
	if err != nil {
		return err
	}
	fmt.Println(ids)

	ctxTODO := context.TODO()
	// query := "What is an interesting fact about the English language?"
	query := "What is an interesting fact about the English language?"
	options := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.8),
	}
	docsFromDb, err := store.SimilaritySearch(ctxTODO, query, 1, options...)
	if err != nil {
		return err
	}
	fmt.Println("len(docsFromDb):", len(docsFromDb))
	fmt.Println("docsFromDb:", docsFromDb)
	fmt.Println("End!!")

	return nil

}

func main2() {
	// Create a new Chroma vector store.
	store, errNs := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace(uuid.New().String()),
	)
	if errNs != nil {
		log.Fatalf("new: %v\n", errNs)
	}

	type meta = map[string]any

	// Add documents to the vector store.
	_, errAd := store.AddDocuments(context.Background(), []schema.Document{
		{PageContent: "Tokyo", Metadata: meta{"population": 9.7, "area": 622}},
		{PageContent: "Kyoto", Metadata: meta{"population": 1.46, "area": 828}},
		{PageContent: "Hiroshima", Metadata: meta{"population": 1.2, "area": 905}},
		{PageContent: "Kazuno", Metadata: meta{"population": 0.04, "area": 707}},
		{PageContent: "Nagoya", Metadata: meta{"population": 2.3, "area": 326}},
		{PageContent: "Toyota", Metadata: meta{"population": 0.42, "area": 918}},
		{PageContent: "Fukuoka", Metadata: meta{"population": 1.59, "area": 341}},
		{PageContent: "Paris", Metadata: meta{"population": 11, "area": 105}},
		{PageContent: "London", Metadata: meta{"population": 9.5, "area": 1572}},
		{PageContent: "Santiago", Metadata: meta{"population": 6.9, "area": 641}},
		{PageContent: "Buenos Aires", Metadata: meta{"population": 15.5, "area": 203}},
		{PageContent: "Rio de Janeiro", Metadata: meta{"population": 13.7, "area": 1200}},
		{PageContent: "Sao Paulo", Metadata: meta{"population": 22.6, "area": 1523}},
	})
	if errAd != nil {
		log.Fatalf("AddDocument: %v\n", errAd)
	}

	ctx := context.TODO()

	type exampleCase struct {
		name         string
		query        string
		numDocuments int
		options      []vectorstores.Option
	}

	type filter = map[string]any

	exampleCases := []exampleCase{
		{
			name:         "Up to 5 Cities in Japan",
			query:        "Which of these are cities are located in Japan?",
			numDocuments: 5,
			options: []vectorstores.Option{
				vectorstores.WithScoreThreshold(0.8),
			},
		},
		{
			name:         "A City in South America",
			query:        "Which of these are cities are located in South America?",
			numDocuments: 1,
			options: []vectorstores.Option{
				vectorstores.WithScoreThreshold(0.8),
			},
		},
		{
			name:         "Large Cities in South America",
			query:        "Which of these are cities are located in South America?",
			numDocuments: 100,
			options: []vectorstores.Option{
				vectorstores.WithFilters(filter{
					"$and": []filter{
						{"area": filter{"$gte": 1000}},
						{"population": filter{"$gte": 13}},
					},
				}),
			},
		},
	}

	// run the example cases
	results := make([][]schema.Document, len(exampleCases))
	for ecI, ec := range exampleCases {
		docs, errSs := store.SimilaritySearch(ctx, ec.query, ec.numDocuments, ec.options...)
		if errSs != nil {
			log.Fatalf("query1: %v\n", errSs)
		}
		results[ecI] = docs
	}

	// print out the results of the run
	fmt.Printf("Results:\n")
	for ecI, ec := range exampleCases {
		texts := make([]string, len(results[ecI]))
		for docI, doc := range results[ecI] {
			texts[docI] = doc.PageContent
		}
		fmt.Printf("%d. case: %s\n", ecI+1, ec.name)
		fmt.Printf("    result: %s\n", strings.Join(texts, ", "))
	}
}
