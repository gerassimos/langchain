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
	db, err = sql.Open("sqlite3", "../resources/db.sqlite")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
}

func run() error {
	var opts []openai.Option
	//use a custom HTTP client to log requests and responses.
	//opts = append(opts, openai.WithHTTPClient(httputil.DebugHTTPClient))

	opts = append(opts, openai.WithCallback(LogInBoxHandler{}))
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
		WriteHtmlReport{},
	}

	sqlTables, err := listSQLiteTables()
	if err != nil {
		return fmt.Errorf("error listing SQLite tables: %w", err)
	}

	customMrklPrefix := "Today is {{.today}}.\n" +
		"You are an AI that has access to a SQLite database.\n" +
		"The database contains the following tables:\n" + sqlTables + "\n\n" +
		"Do not make any assumptions about what tables exist or what columns exist. " +
		"Instead, use the 'describe_tables' tool\n" +
		"Answer the following questions as best you can.\n" +
		"You have access to the following tools:\n" +
		"{{.tool_descriptions}}"

	o1 := agents.WithMaxIterations(0)
	//o2 := agents.WithPromptSuffix(_customMrklSuffix)
	o2 := agents.WithPromptPrefix(customMrklPrefix)

	//will use input variables: "input", "agent_scratchpad"
	agent := agents.NewOneShotAgent(llm,
		agentTools,
		o1, o2)

	executor := agents.NewExecutor(agent)

	// Memory is not working. The {{.history}} input is missing in the prompt.
	// I believe to make it work we need to create a custom prompt template that includes the history.
	// Use of the option `agents.WithPrompt()`
	//conversationBuffer := memory.NewConversationBuffer()
	//agentOption := agents.WithMemory(conversationBuffer)
	//executor := agents.NewExecutor(agent, agentOption)

	//q1 := "How many users are in the database?"
	//q2 := "How many users have provided a shipping address?"
	//q3 := "How many orders are there? Write the result to a html report."
	// NOT working => question:
	// NOT working => "Summarize the top 5 most popular products. Write the result to a report html file."
	//log question

	q1 := "How many users are in the database? Write the result to a html report."
	q2 := "Same question for orders."

	//array of questions
	questions := []string{q1, q2}

	for _, question := range questions {
		fmt.Println("Question:", question)
		answer, err := chains.Run(ctx, executor, question)

		if err != nil {
			return err
		}
		fmt.Println(answer)
	}
	return nil
}
