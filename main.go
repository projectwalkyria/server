package main

import (
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
	createAdmToken(db)
	
	http.HandleFunc("POST /con/{id}", conPost)
	http.HandleFunc("PUT /con/{id}", conPut)
	http.HandleFunc("GET /con/{id}", conGet)
	http.HandleFunc("DELETE /con/{id}", conDelete)

	http.HandleFunc("POST /adm/context", admContextPost)
	http.HandleFunc("GET /adm/context", admContextGet)
	http.HandleFunc("DELETE /adm/context", admContextDelete)

	http.HandleFunc("POST /adm/token", admTokenPost)
	http.HandleFunc("DELETE /adm/token", admTokenDelete)

	http.HandleFunc("POST /adm/token/grant", admTokenGrant)
	http.HandleFunc("DELETE /adm/token/revoke", admTokenRevoke)

	log.Fatal(http.ListenAndServe(":53072", nil))
}
