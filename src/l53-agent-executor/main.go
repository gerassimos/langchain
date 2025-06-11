package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
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
