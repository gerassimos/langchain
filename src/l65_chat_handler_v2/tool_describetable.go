package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tmc/langchaingo/tools"
	"strings"
)

// SqliteDescribeTables Given a list of table names, returns the schema of those tables.
func SqliteDescribeTables(tableNames []string) (string, error) {
	//log
	fmt.Printf("SqliteDescribeTables called with tableNames: %v\n", tableNames)

	// Join table names into a single string for the SQL query
	tables := "'" + strings.Join(tableNames, "', '") + "'"
	query := fmt.Sprintf("SELECT sql FROM sqlite_master WHERE type='table' AND name IN (%s);", tables)

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Collect the results
	var result strings.Builder
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return "", err
		}
		result.WriteString(schema + "\n")
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return "", err
	}

	return result.String(), nil
}

type DescribeTables struct {
}

var _ tools.Tool = DescribeTables{}

func (c DescribeTables) Description() string {
	return `Given a string containing a comma separated list of table names, returns the schema of those tables.`
}

// Name returns the name of the tool.
func (c DescribeTables) Name() string {
	return "describe_tables"
}

// Call executes the SQLite query provided in the input string and returns the results.
// If the query execution fails, it returns an error message.
// If the query errors the error is given in the result to give the
// agent the ability to retry.
func (c DescribeTables) Call(ctx context.Context, input string) (string, error) {
	tableNames := strings.Split(input, ",")
	return SqliteDescribeTables(tableNames)
}
