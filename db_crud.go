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
	logStuff("creating entry table.")
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
	logStuff("creating context table")
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
	logStuff("create token table")
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
	logStuff("create permission table")
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
		logStuff("error on creating entry {" + key + ":" + value + "} on context: " + context)
		logStuff(err.Error())
		return "", "", "", err
	}
	logStuff("creating entry {" + key + ":" + value + "} on context: " + context)
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
		logStuff("error on returning entry " + key + " from context " + context)
		return "", "", "", fmt.Errorf("no entry found for key: %s in context: %s", key, context)
	}

	logStuff("returning entry " + key + " from context " + context)
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
		logStuff("error on updating entry " + key + " with value " + value + " on context :" + context)
		return "", "", "", errors.New("key-value pair doesn't exists")
	}
	logStuff("updating entry " + key + " with value " + value + " on context :" + context)
	return context, key, value, nil
}

func deleteEntry(db *sql.DB, context string, key string) error {
	deleteEntrySQL := `DELETE FROM entry WHERE context = ? AND key = ?`
	result, err := db.Exec(deleteEntrySQL, context, key)
	if err != nil {
		logStuff("error on deleting entry " + key + " from context " + context)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logStuff("error on deleting entry " + key + " from context " + context)
		return errors.New("key-value pair doesn't exists")
	}
	logStuff("deleting entry " + key + " from context " + context)
	return nil
}

func createContext(db *sql.DB, name string) (string, error) {
	insertContextSQL := `INSERT INTO context (name) VALUES (?)`
	_, err := db.Exec(insertContextSQL, name)
	if err != nil {
		logStuff("error on creating context " + name)
		return "", err
	}
	logStuff("creating context " + name)
	return name, nil
}

func getContext(db *sql.DB, name string) (string, error) {
	rows, err := db.Query("SELECT name FROM context WHERE name = ?", name)
	if err != nil {
		logStuff("error on returning context " + name)
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
		logStuff("error on returning context " + name)
		return "", fmt.Errorf("no context : %s", name)
	}

	logStuff("returning context " + name)
	return value, nil
}

func deleteContext(db *sql.DB, name string) error {
	deleteContextSQL := `DELETE FROM context WHERE name = ?`
	result, err := db.Exec(deleteContextSQL, name)
	if err != nil {
		logStuff("error on deleting context " + name)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logStuff("error on deleting context " + name)
		return errors.New("context doesn't exists")
	}
	logStuff("deleting context " + name)
	return nil
}

func createToken(db *sql.DB) (string, error) {
	newUUID := uuid.New().String()
	tokenSha := tokenToSha256(newUUID)
	insertTokenSQL := `INSERT INTO token (token) VALUES (?)`
	_, err := db.Exec(insertTokenSQL, tokenSha)
	if err != nil {
		logStuff("error on creating token")
		return "", err
	}
	logStuff("creating token")
	return newUUID, nil
}

func getPermission(db *sql.DB, token string, context string, reqType string) error {
	tokenSha := tokenToSha256(token)

	rows, err := db.Query(
		"SELECT token FROM permission WHERE token = ? AND context = ? AND permission = ?;",
		tokenSha, context, reqType)
	if err != nil {
		logStuff("error on checking permission for context " + context)
		return err
	}
	defer rows.Close()

	var found bool
	for rows.Next() {
		found = true
	}

	if !found {
		logStuff("error on checking permission for context " + context)
		return errors.New("not authorized")
	}
	logStuff("checking permission for context " + context)
	return nil
}

func deleteToken(db *sql.DB, token string) error {
	tokenSha := tokenToSha256(token)

	deleteTokenPermissionsSQL := `DELETE FROM permission WHERE token = ?`
	_, err := db.Exec(deleteTokenPermissionsSQL, tokenSha)
	if err != nil {
		return err
	}

	var result sql.Result
	deleteTokenSQL := `DELETE FROM token WHERE token = ?`
	result, err = db.Exec(deleteTokenSQL, tokenSha)
	if err != nil {
		logStuff("error on deleting token")
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logStuff("error on deleting token")
		return errors.New("token doesn't exists")
	}

	logStuff("deleting token")
	return nil
}

func grantTokenPermission(db *sql.DB, token string, permission string, context string) (string, string, string, error) {
	tokenSha := tokenToSha256(token)

	insertPermissionSQL := `INSERT INTO permission (token, permission, context) VALUES (?, ?, ?)`
	_, err := db.Exec(insertPermissionSQL, tokenSha, permission, context)
	if err != nil {
		logStuff("error on granting permission on context " + context + " for token")
		return "", "", "", err
	}
	logStuff("granting permission on context " + context + " for token")
	return token, permission, context, nil
}

func rovokeTokenPermission(db *sql.DB, token string, permission string, context string) error {
	tokenSha := tokenToSha256(token)

	deletePermissionSQL := `DELETE FROM permission WHERE token = ? AND permission = ? AND context = ?`
	_, err := db.Exec(deletePermissionSQL, tokenSha, permission, context)
	if err != nil {
		logStuff("error on revoking permission of a token from the context " + context)
		return err
	}

	logStuff("revoking permission of a token from the context " + context)
	return nil
}

func createAdmToken(db *sql.DB) (string, error) {
	_, found, _ := admTokenExists(db)

	if !found {
		token, err := createToken(db)
		if err != nil {
			return "", err
		}
		_, _, _, err = grantTokenPermission(db, token, "ADM_TOKEN_POST", "ALL")
		if err != nil {
			return "", err
		}
		_, _, _, err = grantTokenPermission(db, token, "ADM_TOKEN_DELETE", "ALL")
		if err != nil {
			return "", err
		}
		_, _, _, err = grantTokenPermission(db, token, "ADM_TOKEN_GRANT", "ALL")
		if err != nil {
			return "", err
		}
		_, _, _, err = grantTokenPermission(db, token, "ADM_TOKEN_REVOKE", "ALL")
		if err != nil {
			return "", err
		}
		_, _, _, err = grantTokenPermission(db, token, "ADM_CONTEXT_POST", "ALL")
		if err != nil {
			return "", err
		}
		_, _, _, err = grantTokenPermission(db, token, "ADM_CONTEXT_GET", "ALL")
		if err != nil {
			return "", err
		}
		_, _, _, err = grantTokenPermission(db, token, "ADM_CONTEXT_DELETE", "ALL")
		if err != nil {
			return "", err
		}
		return token, err
	}
	return "", nil
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

func getToken(db *sql.DB, token string) error {
	tokenSha := tokenToSha256(token)

	rows, err := db.Query("SELECT token FROM token WHERE token = ? ", tokenSha)
	if err != nil {
		logStuff("error on checking if token exists")
		return err
	}
	defer rows.Close()

	var found bool
	for rows.Next() {
		found = true
	}

	if !found {
		logStuff("error on checking if token exists")
		return errors.New("token not exists")
	}

	logStuff("checking if token exists")
	return nil
}
