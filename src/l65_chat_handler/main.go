package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"
	"log"
	"os"
)

const ()

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
	db, err = sql.Open("sqlite3", "../resources/db.sqlite")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
}

func run() error {
	var opts []openai.Option
	//use a custom HTTP client to log requests and responses.
	//opts = append(opts, openai.WithHTTPClient(httputil.DebugHTTPClient))

	opts = append(opts, openai.WithCallback(callbacks.LogHandler{}))
	// We can construct an LLMChain from a PromptTemplate and an LLM.
	llm, err := openai.New(opts...)
	//llm, err := openai.New()
	if err != nil {
		return err
	}
	//llm.CallbacksHandler = callbacks.LogHandler{}
	ctx := context.Background()

	agentTools := []tools.Tool{
		RunSqliteQuery{},
		DescribeTables{},
	}

	customMessageAboutSqlTables, err := listSQLiteTables()
	_customMrklSuffix := "Take into account the following information about the database:\n" +
		customMessageAboutSqlTables + "\n\n" +
		"Begin!" + "\n\n" +
		"Question: {{.input}}" + "\n" +
		"{{.agent_scratchpad}}"

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
	printTemplate(agent)
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
func printTemplate(agent *agents.OneShotZeroAgent) {
	chain := agent.Chain
	if llmChain, ok := chain.(*chains.LLMChain); ok {
		prompt := llmChain.Prompt
		p := prompt.(prompts.PromptTemplate)
		t := p.Template
		fmt.Printf("Template:\n %s", t)
	}
}
