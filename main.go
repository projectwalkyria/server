package main

import (
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/adm", adm)
    http.HandleFunc("POST /con/{id}", conPost)
    http.HandleFunc("PUT /con/{id}", conPut)
    http.HandleFunc("GET /con/{id}", conGet)
    http.HandleFunc("DELETE /con/{id}", conDelete)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
