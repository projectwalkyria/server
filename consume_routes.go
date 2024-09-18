package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func getHeaderAuthToken(r *http.Request) (string, error) {
	// Extract Bearer Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	// Split the "Bearer <token>" part
	tokenParts := strings.Split(authHeader, " ")

	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	// The actual token is the second part
	return tokenParts[1], nil
}

func parseConRequest(r *http.Request) (string, string, string, error) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		return "", "", "", err
	}

	defer r.Body.Close()

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", "", "", err
	}

	context := r.PathValue("id")

	if len(data) == 1 {
		// Access the first key-value pair (manually in this case, as Go doesn't provide a direct way to access the first element in a map).
		for key, value := range data {

			// Assert that `value` is of type `string`. You can also handle other types based on your input.
			valueStr, ok := value.(string)
			if !ok {
				return "", "", "", errors.New("value is not a string")
			}
			return context, key, valueStr, nil
		}
	}
	return "", "", "", errors.New("JSON body must have 1 key-value pair")
}

func conPost(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	context, key, value, err := parseConRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := connectSQLite()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	err = getPermission(db, authToken, context, "POST")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	context, err = getContext(db, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _, _, err = createEntry(db, context, key, value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func conPut(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	context := r.PathValue("id")

	db, err := connectSQLite()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	err = getPermission(db, authToken, context, "PUT")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = getContext(db, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	context, key, value, err := parseConRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _, _, err = updateEntry(db, context, key, value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func conGet(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	context := r.PathValue("id")

	db, err := connectSQLite()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	err = getPermission(db, authToken, context, "PUT")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = getContext(db, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	context, _, value, err := parseConRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var key string
	context, key, value, err = getEntry(db, context, value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a response for success
	successResponse := map[string]string{
		"context": context,
		"key":     key,
		"value":   value,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)

}

func conDelete(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	context := r.PathValue("id")

	db, err := connectSQLite()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	err = getPermission(db, authToken, context, "PUT")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = getContext(db, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	context, _, key, err := parseConRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = deleteEntry(db, context, key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
}
