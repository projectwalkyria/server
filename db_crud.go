package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

// SQLite connection
func connectSQLite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS entry (
		context TEXT,
		key TEXT,
		value TEXT,
		UNIQUE (context, key)
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table created successfully!")
}

func createEntry(db *sql.DB, context string, key string, value string) (string, string, string, error) {	
	insertEntrySQL := `INSERT INTO entry(key, value, context) VALUES (?, ?, ?)`
	_, err := db.Exec(insertEntrySQL, key, value, context)
	if err != nil {
		return "", "", "", err
	}
	return context, key, value, nil
}


func getEntry(db *sql.DB, context string, key string) (string, string, string, error) {
	rows, err := db.Query("SELECT context, key, value FROM entry WHERE key = ? AND context = ?", key, context)
	if err != nil {
		return "", "", "", err
	}
	defer rows.Close()
	
	var value string
	for rows.Next() {
		rows.Scan(&context, &key, &value)
	}
	return context, key, value, nil
}

func getEntries(db *sql.DB, context string) {
	rows, err := db.Query("SELECT context, key, value FROM entry WHERE context = ?", context)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var context, key, value string
		rows.Scan(&context, &key, &value)
		fmt.Printf("entry: %s, %s, %s\n", context, key, value)
	}
}

func updateEntry(db *sql.DB, context string, key string, value string) (string, string, string, error) {
	updateEntrySQL := `UPDATE entry SET value = ? WHERE context = ? AND key = ?`
	_, err := db.Exec(updateEntrySQL, value, context, key)
	if err != nil {
		return "", "", "", err
	}
	return context, key, value, nil
}

func deleteEntry(db *sql.DB, context string, key string) error {
	deleteEntrySQL := `DELETE FROM entry WHERE context = ? AND key = ?`
	_, err := db.Exec(deleteEntrySQL, context, key)
	if err != nil {
		return err
	}
	return nil
}