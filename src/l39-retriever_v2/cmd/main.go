package main

import (
	"context"
	"fmt"
	"log"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
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
	// We can construct an LLMChain from a PromptTemplate and an LLM.
	llm, err := openai.New()
	if err != nil {
		return err
	}
	// ctx := context.Background()

	// ----------------------------------------------
	// Retrievers.GetRelevantDocuments
	// ----------------------------------------------
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Chroma vector store (db client).
	// storeNs:=uuid.New().String()
	storeNs := "chroma-ns-1"
	store, err := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace(storeNs),
		chroma.WithEmbedder(embedder),
	)
	if err != nil {
		return err
	}

	searchQuery := "What is an interesting fact about the English language?"
	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.80), // use for precision, when you want to get only the most relevant documents
		//vectorstores.WithNameSpace(""),            // use for set a namespace in the storage
		//vectorstores.WithFilters(map[string]interface{}{"language": "en"}), // use for filter the documents
		//vectorstores.WithEmbedder(embedder), // use when you want add documents or doing similarity search
		//vectorstores.WithDeduplicater(vectorstores.NewSimpleDeduplicater()), //  This is useful to prevent wasting time on creating an embedding
	}

	retriever := vectorstores.ToRetriever(store, 3, optionsVector...)
	// search
	docs, err := retriever.GetRelevantDocuments(context.Background(), searchQuery)
	if err != nil {
		log.Fatal(err)
	}
	for _, doc := range docs {
		fmt.Println("doc.Score", doc.Score)
		fmt.Println(doc.PageContent)
		fmt.Println()
	}

	// ----------------------------------------------
	// Chain type="stuff"
	// ----------------------------------------------
	// We can use LoadStuffQA to create a chain that takes input documents and a question,
	// stuffs all the documents into the prompt of the llm and returns an answer to the
	// question. It is suitable for a small number of documents.
	stuffQAChain := chains.LoadStuffQA(llm)

	outputValues, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": docs,
		"question":        searchQuery,
	})
	if err != nil {
		return err
	}

	out, ok := outputValues["text"].(string)
	if !ok {
		return fmt.Errorf("invalid chain return")
	}
	fmt.Println(out)

	return nil
}
