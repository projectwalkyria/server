package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

func conPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

    var data map[string]interface{}

	err = json.Unmarshal(body, &data)
    if err != nil {
        http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
        return
    }

    for key, value := range data {
        valueStr, err := json.Marshal(value)
        if err != nil {
            http.Error(w, "Failed to convert value to JSON", http.StatusInternalServerError)
            return
        }

		filename := key+".data"

		_, err = os.Stat(filename)

		if err != nil {
			err = os.WriteFile(filename, valueStr, 0666)
			if err != nil {
				http.Error(w, "Failed to write file", http.StatusInternalServerError)
				return
			}
		}
    }
}

func conPut(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

    var data map[string]interface{}

	err = json.Unmarshal(body, &data)
    if err != nil {
        http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
        return
    }

    for key, value := range data {
        valueStr, err := json.Marshal(value)
        if err != nil {
            http.Error(w, "Failed to convert value to JSON", http.StatusInternalServerError)
            return
        }

		filename := key+".data"

		_, err = os.Stat(filename)

		if err == nil {
			err = os.WriteFile(filename, valueStr, 0666)
			if err != nil {
				http.Error(w, "Failed to write file", http.StatusInternalServerError)
				return
			}
		}
    }
}

func conGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.PathValue("id"), "get")
}

func conDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.PathValue("id"), "delete")
}
