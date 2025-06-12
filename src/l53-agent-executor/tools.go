package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tmc/langchaingo/tools"
	"strings"
)

// RunSQLiteQuery executes a given SQLite query and returns the results.
func RunSQLiteQuery(query string) (string, error) {
	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Parse the results
	var results [][]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return "", err
		}

		// Append the row to the results
		results = append(results, values)
	}
	// convert results to a formated string

	return formatResults(results), nil
}

// formatResults converts a two-dimensional slice of interface{} to a formatted string.
func formatResults(results [][]interface{}) string {
	var resultString string
	for _, row := range results {
		for _, value := range row {
			resultString += fmt.Sprintf("%v\t", value) // Convert each value to a string and add a tab
		}
		resultString += "\n" // Add a newline after each row
	}
	return resultString
}

// RunSqliteQuery is a tool that can execute SQLite queries.
type RunSqliteQuery struct {
}

var _ tools.Tool = RunSqliteQuery{}

// Description returns a string describing the calculator tool.
func (c RunSqliteQuery) Description() string {
	return `Run a sqlite query.`
}

// Name returns the name of the tool.
func (c RunSqliteQuery) Name() string {
	return "run_sqlite_query"
}

// Call executes the SQLite query provided in the input string and returns the results.
// If the query execution fails, it returns an error message.
// If the query errors the error is given in the result to give the
// agent the ability to retry.
func (c RunSqliteQuery) Call(ctx context.Context, input string) (string, error) {

	fmt.Printf("Running SQLite query: %s\n", input)
	// trim double quotes from input
	input = strings.Trim(input, "\"")
	result, err := RunSQLiteQuery(input)
	if err != nil {
		return "", fmt.Errorf("error running RunSQLiteQuery with input \"%s\": %w", input, err)
	}
	//log result
	fmt.Printf("SQLite query result: %s\n", result)

	// trim result to remove trailing newlines and tabs
	result = strings.TrimSpace(result)

	return result, nil
}
