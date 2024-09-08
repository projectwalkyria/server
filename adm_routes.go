package main

import (
    "fmt"
    "net/http"
)

func adm(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "adm1")
}