package main

import (
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

// a function that lists all tables in the SQLite database
func listSQLiteTables() (string, error) {

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return "", err
		}
		tables = append(tables, tableName)
	}

	return "The available tables in the database are:\n" + strings.Join(tables, "\n"), nil
}
