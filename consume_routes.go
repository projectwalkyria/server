package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// getHeaderAuthToken extracts the Bearer token from the Authorization header in an HTTP request.
// It returns the token string or an error if the header is missing or improperly formatted.
func getHeaderAuthToken(r *http.Request) (string, error) {
	// Get the value of the Authorization header from the request.
	authHeader := r.Header.Get("Authorization")
	
	// If the Authorization header is empty, return an error indicating it's missing.
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	// Split the Authorization header into two parts: the scheme ("Bearer") and the token.
	tokenParts := strings.Split(authHeader, " ")

	// Check that the header is correctly formatted as "Bearer <token>".
	// If the length isn't 2 or the scheme isn't "Bearer", return an error for invalid format.
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	// Return the token (the second part of the split header) and no error.
	return tokenParts[1], nil
}

// parseConRequest parses an HTTP request and extracts a context, key, and value from the request body.
// It expects a JSON body with exactly one key-value pair and returns an error if this condition isn't met.
func parseConRequest(r *http.Request) (string, string, string, error) {
	// Read the entire body of the request.
	body, err := io.ReadAll(r.Body)

	// If reading the body fails, return an empty response and the error.
	if err != nil {
		return "", "", "", err
	}

	// Ensure the request body is closed after the function completes.
	defer r.Body.Close()

	// Create a map to hold the parsed JSON data.
	var data map[string]interface{}

	// Unmarshal the JSON body into the `data` map.
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", "", "", err // Return an error if the JSON parsing fails.
	}

	// Extract the "id" from the request path as the context (assuming r.PathValue is a valid method).
	context := r.PathValue("id")

	// Ensure the JSON body contains exactly one key-value pair.
	if len(data) == 1 {
		// Since Go maps do not have a direct way to access the first element, we manually iterate over the map.
		for key, value := range data {

			// Type assert the value to a string. If the value is not a string, return an error.
			valueStr, ok := value.(string)
			if !ok {
				return "", "", "", errors.New("value is not a string") // Handle the case where the value isn't a string.
			}

			// If everything is successful, return the context (id), the key, and the string value.
			return context, key, valueStr, nil
		}
	}
	
	// Return an error if the JSON body does not contain exactly one key-value pair.
	return "", "", "", errors.New("JSON body must have 1 key-value pair")
}

// conPost handles HTTP POST requests to create a new entry in the database.
// It performs authentication, parses the request, checks permissions, and then inserts the entry into the database.
func conPost(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization token from the request header.
	authToken, err := getHeaderAuthToken(r)
	// If the token is missing or invalid, return a 401 Unauthorized error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Parse the context, key, and value from the request body.
	context, key, value, err := parseConRequest(r)
	// If the request body is invalid or improperly formatted, return a 400 Bad Request error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Establish a connection to the SQLite database.
	db, err := connectSQLite()
	// If the connection to the database fails, return a 400 Bad Request error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Ensure the database connection is closed when the function completes.
	defer db.Close()

	// Verify that the authenticated user has permission to perform the POST operation in the given context.
	err = getPermission(db, authToken, context, "POST")
	// If the user lacks the necessary permissions, return a 401 Unauthorized error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Retrieve and verify the actual context from the database (e.g., normalizing or ensuring its existence).
	context, err = getContext(db, context)
	// If the context is invalid or not found, return a 400 Bad Request error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new entry in the database using the context, key, and value.
	_, _, _, err = createEntry(db, context, key, value)
	// If the entry creation fails, return a 400 Bad Request error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the response header to indicate the content type is JSON.
	w.Header().Set("Content-Type", "application/json")
	// Respond with a 201 Created status code, indicating that the entry was successfully created.
	w.WriteHeader(http.StatusCreated)
}

// conPut handles HTTP PUT requests to update an existing entry in the database.
func conPut(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization token from the request header.
	authToken, err := getHeaderAuthToken(r)
	// If the token is missing or invalid, return a 401 Unauthorized error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Extract the context (id) from the request path.
	context := r.PathValue("id")

	// Establish a connection to the SQLite database.
	db, err := connectSQLite()
	// If the connection to the database fails, return a 400 Bad Request error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Ensure the database connection is closed when the function completes.
	defer db.Close()

	// Verify that the authenticated user has permission to perform the PUT operation in the given context.
	err = getPermission(db, authToken, context, "PUT")
	// If the user lacks the necessary permissions, return a 401 Unauthorized error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Retrieve and verify the actual context from the database.
	_, err = getContext(db, context)
	// If the context is invalid or not found, return a 400 Bad Request error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the context, key, and value from the request body.
	context, key, value, err := parseConRequest(r)
	// If the request body is invalid or improperly formatted, return a 400 Bad Request error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the existing entry in the database using the context, key, and value.
	_, _, _, err = updateEntry(db, context, key, value)
	// If the entry update fails, return a 400 Bad Request error.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the response header to indicate the content type is JSON.
	w.Header().Set("Content-Type", "application/json")
	// Respond with a 200 OK status code, indicating that the entry was successfully updated.
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

	err = getPermission(db, authToken, context, "GET")
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
	_, key, value, err = getEntry(db, context, value)

	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Create a response for success
	successResponse := map[string]string{
		key: value,
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
