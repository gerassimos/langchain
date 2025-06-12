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

	//will use input variables: "input", "agent_scratchpad"
	agent := agents.NewOneShotAgent(llm,
		agentTools,
		agents.WithMaxIterations(0))
	executor := agents.NewExecutor(agent)

	question := "How many users are in the database?"
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
