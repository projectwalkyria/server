package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// connectSQLite establishes a connection to the SQLite database.
func connectSQLite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// createContextDataTable creates a table for storing key-value pairs in a specific context.
func createContextDataTable(db *sql.DB, context string) error {
	logStuff("creating context " + context + " table.")
	createTableSQL := "CREATE TABLE IF NOT EXISTS " + context + " (key TEXT UNIQUE,value TEXT);"

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// deleteContextDataTable deletes a table for a specific context.
func deleteContextDataTable(db *sql.DB, context string) error {
	logStuff("deleting context " + context + " table.")
	createTableSQL := "DROP TABLE " + context + ";"

	_, err := db.Exec(createTableSQL, context)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// createContextTable creates a table to store context names.
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

// createTokenTable creates a table to store tokens.
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

// createPermissionTable creates a table to store permissions associated with tokens and contexts.
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

// createEntry inserts a key-value pair into a specific context table.
func createEntry(db *sql.DB, context string, key string, value string) (string, string, string, error) {
	insertEntrySQL := "INSERT INTO " + context + ` (key, value) VALUES (?, ?)`
	_, err := db.Exec(insertEntrySQL, key, value)
	if err != nil {
		logStuff("error on creating entry {" + key + ":" + value + "} on context: " + context)
		logStuff(err.Error())
		return "", "", "", err
	}
	logStuff("creating entry {" + key + ":" + value + "} on context: " + context)
	return context, key, value, nil
}

// getEntry retrieves the value for a specific key from a context table.
func getEntry(db *sql.DB, context string, key string) (string, string, string, error) {
	logStuff("Searching for key " + key + " on context " + context)
	rows, err := db.Query("SELECT value FROM " + context + " WHERE key = '" + key + "';")
	if err != nil {
		logStuff("error " + err.Error())
		return "", "", "", err
	}
	defer rows.Close()
	var value string
	var found bool
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			logStuff("no entry found for key " + key + " con context " + context)
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

// updateEntry updates the value for a specific key in a context table.
func updateEntry(db *sql.DB, context string, key string, value string) (string, string, string, error) {
	updateEntrySQL := "UPDATE " + context + ` SET value = ? WHERE key = ?`
	result, err := db.Exec(updateEntrySQL, value, key)
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

// deleteEntry deletes a key-value pair from a context table.
func deleteEntry(db *sql.DB, context string, key string) error {
	deleteEntrySQL := "DELETE FROM " + context + ` WHERE key = ?`
	result, err := db.Exec(deleteEntrySQL, key)
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

// createContext inserts a new context name into the context table.
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

// getContext retrieves a context name from the context table.
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

// deleteContext deletes a context name from the context table.
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

// createToken generates a new token and inserts it into the token table.
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

// getPermission checks if a token has a specific permission in a context.
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

// deleteToken deletes a token and its associated permissions from the database.
func deleteToken(db *sql.DB, token string) error {
    // Convert the token to its SHA-256 hash.
    tokenSha := tokenToSha256(token)

    // SQL query to delete permissions associated with the token.
    deleteTokenPermissionsSQL := `DELETE FROM permission WHERE token = ?`
    _, err := db.Exec(deleteTokenPermissionsSQL, tokenSha)
    if err != nil {
        return err
    }

    // SQL query to delete the token itself.
    var result sql.Result
    deleteTokenSQL := `DELETE FROM token WHERE token = ?`
    result, err = db.Exec(deleteTokenSQL, tokenSha)
    if err != nil {
        logStuff("error on deleting token")
        return err
    }

    // Check if any rows were affected by the delete operation.
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        logStuff("error on deleting token")
        return errors.New("token doesn't exist")
    }

    logStuff("deleting token")
    return nil
}

// grantTokenPermission grants a specific permission to a token within a given context.
func grantTokenPermission(db *sql.DB, token string, permission string, context string) (string, string, string, error) {
    // Convert the token to its SHA-256 hash.
    tokenSha := tokenToSha256(token)

    // SQL query to insert the new permission.
    insertPermissionSQL := `INSERT INTO permission (token, permission, context) VALUES (?, ?, ?)`
    _, err := db.Exec(insertPermissionSQL, tokenSha, permission, context)
    if err != nil {
        logStuff("error on granting permission on context " + context + " for token")
        return "", "", "", err
    }
    logStuff("granting permission on context " + context + " for token")
    return token, permission, context, nil
}

// rovokeTokenPermission revokes a specific permission from a token within a given context.
func rovokeTokenPermission(db *sql.DB, token string, permission string, context string) error {
    // Convert the token to its SHA-256 hash.
    tokenSha := tokenToSha256(token)

    // SQL query to delete the permission.
    deletePermissionSQL := `DELETE FROM permission WHERE token = ? AND permission = ? AND context = ?`
    _, err := db.Exec(deletePermissionSQL, tokenSha, permission, context)
    if err != nil {
        logStuff("error on revoking permission of a token from the context " + context)
        return err
    }

    logStuff("revoking permission of a token from the context " + context)
    return nil
}

// createAdmToken creates an admin token with all necessary permissions if it doesn't already exist.
func createAdmToken(db *sql.DB) (string, error) {
    // Check if an admin token already exists.
    _, found, _ := admTokenExists(db)

    if !found {
        // Create a new token.
        token, err := createToken(db)
        if err != nil {
            return "", err
        }
        // Grant all necessary permissions to the new admin token.
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

// admTokenExists checks if an admin token already exists in the database.
func admTokenExists(db *sql.DB) (string, bool, error) {
    // SQL query to find any admin tokens.
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
        return "", false, errors.New("adm token does not exist")
    }
    return value, found, nil
}

// getToken checks if a token exists in the database.
func getToken(db *sql.DB, token string) error {
    // Convert the token to its SHA-256 hash.
    tokenSha := tokenToSha256(token)

    // SQL query to check if the token exists.
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
        return errors.New("token does not exist")
    }

    logStuff("checking if token exists")
    return nil
}
