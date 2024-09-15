package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
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

func createTokenTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS token (
		token TEXT UNIQUE
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table token created successfully!")
}

func createPermissionTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS permission (
		token TEXT,
		permission TEXT,
		context TEXT,
		UNIQUE (token, permission, context)
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table permission created successfully!")
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
		return "", "", "", errors.New("key-value pair doesn't exists")
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
		return errors.New("key-value pair doesn't exists")
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
		return errors.New("context doesn't exists")
	}
	return nil
}

func createToken(db *sql.DB) (string, error) {
	newUUID := uuid.New().String()
	insertTokenSQL := `INSERT INTO token (token) VALUES (?)`
	_, err := db.Exec(insertTokenSQL, newUUID)
	if err != nil {
		return "", err
	}
	return newUUID, nil
}

func getPermission(db *sql.DB, token string, context string, reqType string) error {
	rows, err := db.Query(
		"SELECT token FROM permission WHERE token = ? AND context = ? AND permission = ?;",
		token, context, reqType)
	if err != nil {
		return err
	}
	defer rows.Close()

	var found bool
	for rows.Next() {
		found = true
	}

	if !found {
		return errors.New("not authorized")
	}

	return nil
}

func deleteToken(db *sql.DB, token string) error {
	deleteTokenPermissionsSQL := `DELETE FROM permission WHERE token = ?`
	_, err := db.Exec(deleteTokenPermissionsSQL, token)
	if err != nil {
		return err
	}

	var result sql.Result
	deleteTokenSQL := `DELETE FROM token WHERE token = ?`
	result, err = db.Exec(deleteTokenSQL, token)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("token doesn't exists")
	}
	return nil
}

func grantTokenPermission(db *sql.DB, token string, permission string, context string) (string, string, string, error) {
	insertPermissionSQL := `INSERT INTO permission (token, permission, context) VALUES (?, ?, ?)`
	_, err := db.Exec(insertPermissionSQL, token, permission, context)
	if err != nil {
		return "", "", "", err
	}
	return token, permission, context, nil
}

func rovokeTokenPermission(db *sql.DB, token string, permission string, context string) error {
	deletePermissionSQL := `DELETE FROM permission WHERE token = ? AND permission = ? AND context = ?`
	result, err := db.Exec(deletePermissionSQL, token, permission, context)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("token doesn't exists")
	}
	return nil
}

func createAdmToken(db *sql.DB) (string, error) {
	var token string
	var found bool
	token, found, _ = admTokenExists(db)

	if !found {
		token, err := createToken(db)
		if err != nil {
			return "", err
		}
		token, _, _, err = grantTokenPermission(db, token, "ADM_TOKEN_POST", "ALL")
		if err != nil {
			return "", err
		}
		token, _, _, err = grantTokenPermission(db, token, "ADM_TOKEN_DELETE", "ALL")
		if err != nil {
			return "", err
		}
		token, _, _, err = grantTokenPermission(db, token, "ADM_TOKEN_GRANT", "ALL")
		if err != nil {
			return "", err
		}
		token, _, _, err = grantTokenPermission(db, token, "ADM_TOKEN_REVOKE", "ALL")
		if err != nil {
			return "", err
		}
		token, _, _, err = grantTokenPermission(db, token, "ADM_CONTEXT_POST", "ALL")
		if err != nil {
			return "", err
		}
		token, _, _, err = grantTokenPermission(db, token, "ADM_CONTEXT_GET", "ALL")
		if err != nil {
			return "", err
		}
		token, _, _, err = grantTokenPermission(db, token, "ADM_CONTEXT_DELETE", "ALL")
		if err != nil {
			return "", err
		}
		return token, err
	}
	return token, nil
}

func admTokenExists(db *sql.DB) (string, bool, error) {
	rows, err := db.Query("SELECT token FROM permission WHERE permission like 'ADM_TOKEN%'")
	if err != nil {
		return "", false, err
	}
	defer rows.Close()
	var value string
	var found bool
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return "", false, err
		}
		found = true
	}
	if !found {
		return "", false, errors.New("adm token does not exists")
	}
	return value, found, nil
}
