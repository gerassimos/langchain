package main

import (
	"context"
	"fmt"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
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

	// Create a new Chroma vector store (db client).
	// storeNs:=uuid.New().String()
	storeNs := "chroma-ns-1"
	store, err := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace(storeNs),
	)
	if err != nil {
		return err
	}

	ctxTODO := context.TODO()
	// query := "What is an interesting fact about the English language?"
	query := "What is an interesting fact about the English language?"
	// query := "English"
	fmt.Println("query:", query)
	options := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.0),
	}
	docsFromDb, err := store.SimilaritySearch(ctxTODO, query, 1, options...)
	// docsFromDb, err := store.SimilaritySearch(ctxTODO, query, 1)
	if err != nil {
		return err
	}
	fmt.Println("len(docsFromDb):", len(docsFromDb))
	// fmt.Println("docsFromDb:", docsFromDb)
	for _, doc := range docsFromDb {
		fmt.Println("doc.Score", doc.Score)
		fmt.Println(doc.PageContent)
		fmt.Println()
	}

	// We can construct an LLMChain from a PromptTemplate and an LLM.
	llm, err := openai.New()
	if err != nil {
		return err
	}
	ctx := context.Background()

	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"Use the following interesting fact to answer the users question {{.fact}}.",
			[]string{"fact"},
		),
		prompts.NewHumanMessagePromptTemplate(
			"Here is the user's question: {{.question}}",
			[]string{"question"},
		),
	})

	llmChain := chains.NewLLMChain(llm, prompt)
	// Otherwise the call function must be used.
	outputValues, err := chains.Call(ctx, llmChain, map[string]any{
		"fact":     docsFromDb,
		"question": query,
	})
	if err != nil {
		return err
	}

	//loop over outputValues
	for key, value := range outputValues {
		fmt.Printf("key: %s: value: %s\n", key, value)

	}
	return nil
}
