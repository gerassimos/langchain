package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// RunSQLiteQuery executes a given SQLite query and returns the results.
func RunSQLiteQuery(query string) ([][]interface{}, error) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the results
	var results [][]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
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
			return nil, err
		}

		// Append the row to the results
		results = append(results, values)
	}

	return results, nil
}
