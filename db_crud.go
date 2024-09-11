package main

import (
	"errors"
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

func createEntryTable(db *sql.DB) {
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

	fmt.Println("Table entry created successfully!")
}

func createContextTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS context (
		name TEXT UNIQUE
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table context created successfully!")
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
	var found bool
	for rows.Next() {
		err := rows.Scan(&context, &key, &value)
		if err != nil {
			return "", "", "", err
		}
		found = true
	}

	if !found {
		return "", "", "", fmt.Errorf("no entry found for key: %s in context: %s", key, context)
	}
	
	return context, key, value, nil
}

func updateEntry(db *sql.DB, context string, key string, value string) (string, string, string, error) {
	updateEntrySQL := `UPDATE entry SET value = ? WHERE context = ? AND key = ?`
	result, err := db.Exec(updateEntrySQL, value, context, key)
	if err != nil {
		return "", "", "", err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return "", "", "", errors.New("Key-value pair doesn't exists.")
	}
	return context, key, value, nil
}

func deleteEntry(db *sql.DB, context string, key string) error {
	deleteEntrySQL := `DELETE FROM entry WHERE context = ? AND key = ?`
	result, err := db.Exec(deleteEntrySQL, context, key)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("Key-value pair doesn't exists.")
	}
	return nil
}

func createContext(db *sql.DB, name string) (string, error) {	
	insertContextSQL := `INSERT INTO context (name) VALUES (?)`
	_, err := db.Exec(insertContextSQL, name)
	if err != nil {
		return "", err
	}
	return name, nil
}


func getContext(db *sql.DB, name string) (string, error) {
	rows, err := db.Query("SELECT name FROM context WHERE name = ?", name)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	
	var value string
	var found bool
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return "", err
		}
		found = true
	}

	if !found {
		return "", fmt.Errorf("no context : %s", name)
	}

	return value, nil
}

func deleteContext(db *sql.DB, name string) error {
	deleteContextSQL := `DELETE FROM context WHERE name = ?`
	result, err := db.Exec(deleteContextSQL, name)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("Context doesn't exists.")
	}
	return nil
}