package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
	"log"
	"os"
)

const (
	_customMrklSuffix = `Begin!
	The available tables in the database are:
	addresses
	carts 
	orders
	products
	users
	orders_products

Question: {{.input}}
{{.agent_scratchpad}}`
)

func main() {
	defer db.Close()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var db *sql.DB

func init() {
	var err error
	// Initialize the database connection
	db, err = sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
}

func run() error {
	llm, err := openai.New()
	if err != nil {
		return err
	}
	ctx := context.Background()

	agentTools := []tools.Tool{
		RunSqliteQuery{},
	}

	helperMessage := `You are an AI that has access to a SQLite database. 
	The available tables in the database are:
	addresses
	carts 
	orders
	products
	users
	orders_products
`
	fmt.Printf("Helper message: %s\n", helperMessage)

	o1 := agents.WithMaxIterations(0)
	o2 := agents.WithPromptSuffix(_customMrklSuffix)

	fmt.Printf("Options: %v %v %v", o1, o2)

	//openAIOption := agents.NewOpenAIOption()
	//o2 := openAIOption.WithSystemMessage(helperMessage)
	//o3 := openAIOption.WithExtraMessages([]prompts.MessageFormatter{
	//	prompts.NewHumanMessagePromptTemplate("please be strict", nil),
	//})
	//agent := agents.NewOpenAIFunctionsAgent(llm,
	//	agentTools, o2, o3)

	//will use input variables: "input", "agent_scratchpad"
	agent := agents.NewOneShotAgent(llm,
		agentTools,
		o1, o2)

	executor := agents.NewExecutor(agent)

	question := "How many users are in the database?"
	question = "How many users have provided a shipping address?"
	//log question
	fmt.Println("Question:", question)
	answer, err := chains.Run(ctx, executor, question)
	fmt.Println(answer)
	return err
}

func runSimpleQuery() error {
	// Example usage of RunSQLiteQuery
	query := "SELECT * FROM users" // Replace with your actual query
	results, err := RunSQLiteQuery(query)
	if err != nil {
		//log.Fatalf("Error running query: %v", err)
		return fmt.Errorf("error running query: %w", err)
	}

	// Print the results
	for _, row := range results {
		fmt.Println(row)
	}
	return nil
}
