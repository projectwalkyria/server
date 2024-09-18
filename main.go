package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Connect to SQLite
	db, err := connectSQLite() // For SQLite
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Create Table

	createEntryTable(db)
	createContextTable(db)
	createTokenTable(db)
	createPermissionTable(db)

	var token string

	token, _ = createAdmToken(db)

	if token != "" {
		fmt.Println("Master adm token : " + token)
	} else {
		fmt.Println("ADM TOKEN ALREADY CREATED")
	}

	http.HandleFunc("POST /con/{id}", conPost)
	http.HandleFunc("PUT /con/{id}", conPut)
	http.HandleFunc("GET /con/{id}", conGet)
	http.HandleFunc("DELETE /con/{id}", conDelete)

	http.HandleFunc("POST /adm/token", admTokenPost)
	// CREATE : check if adm token exists
	http.HandleFunc("DELETE /adm/token", admTokenDelete)

	http.HandleFunc("POST /adm/context", admContextPost)
	http.HandleFunc("GET /adm/context", admContextGet)
	// CREATE : check contexts, return all contexts
	http.HandleFunc("DELETE /adm/context", admContextDelete)

	http.HandleFunc("POST /adm/token/grant", admTokenGrant)
	// CREATE : check a token grant on context
	http.HandleFunc("DELETE /adm/token/revoke", admTokenRevoke)

	log.Fatal(http.ListenAndServe(":53072", nil))
}
