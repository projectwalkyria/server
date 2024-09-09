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

	createTable(db)
	
	http.HandleFunc("POST /con/{id}", conPost)
	http.HandleFunc("PUT /con/{id}", conPut)
	http.HandleFunc("GET /con/{id}", conGet)
	http.HandleFunc("DELETE /con/{id}", conDelete)
	log.Fatal(http.ListenAndServe(":53072", nil))
}
