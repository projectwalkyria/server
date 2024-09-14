package main

import (
	"errors"
    "encoding/json"
    "io"
    "net/http"
	"strings"
)

func getHeaderAuthToken(w http.ResponseWriter, r *http.Request) (string, error) {
	// Extract Bearer Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header missing")
	}

	// Split the "Bearer <token>" part
	tokenParts := strings.Split(authHeader, " ")

	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", errors.New("Invalid Authorization header format")
	}

	// The actual token is the second part
	return tokenParts[1], nil
}

func parseConRequest(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
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

	if len(data) == 1 {
		// Access the first key-value pair (manually in this case, as Go doesn't provide a direct way to access the first element in a map).
		for key, value := range data {
			context := r.PathValue("id") // Assuming you're extracting the 'id' from the URL query parameters

			// Assert that `value` is of type `string`. You can also handle other types based on your input.
			valueStr, ok := value.(string)
			if !ok {
				return "", "", "", errors.New("Value is not a string")
			}
			return context, key, valueStr, nil
		}
	}
	return "", "", "", errors.New("JSON body must have 1 key-value pair")
}

func conPost(w http.ResponseWriter, r *http.Request) {
	context, key, value, err := parseConRequest(w, r)



	// Connect to SQLite
	db, err := connectSQLite() // For SQLite
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	authToken, err := getHeaderAuthToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = getPermission(db, authToken, context, "POST")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	context, key, value, err = createEntry(db, context, key, value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Create a response for success
    successResponse := map[string]string{
        "context": context,
        "key": key,
		"value": value,
    }

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(successResponse)
    return
}

func conPut(w http.ResponseWriter, r *http.Request) {
	context, key, value, err := parseConRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to SQLite
	db, err := connectSQLite() // For SQLite
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()	

	authToken, err := getHeaderAuthToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = getPermission(db, authToken, context, "PUT")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	context, key, value, err = updateEntry(db, context, key, value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Create a response for success
    successResponse := map[string]string{
        "context": context,
        "key": key,
		"value": value,
    }

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(successResponse)
    return
}

func conGet(w http.ResponseWriter, r *http.Request) {
	context, key, value, err := parseConRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to SQLite
	db, err := connectSQLite() // For SQLite
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()
	
	authToken, err := getHeaderAuthToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = getPermission(db, authToken, context, "GET")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	context, key, value, err = getEntry(db, context, value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Create a response for success
    successResponse := map[string]string{
        "context": context,
        "key": key,
		"value": value,
    }

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(successResponse)
	return
}

func conDelete(w http.ResponseWriter, r *http.Request) {
	context, key, value, err := parseConRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to SQLite
	db, err := connectSQLite() // For SQLite
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()
	
	authToken, err := getHeaderAuthToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = getPermission(db, authToken, context, "DELETE")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	err = deleteEntry(db, context, value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    return
}