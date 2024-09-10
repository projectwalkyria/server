package main

import (
	"errors"
    "encoding/json"
    "io"
    "net/http"
)

func parseRequest(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
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
	context, key, value, err := parseRequest(w, r)
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
	context, key, value, err := parseRequest(w, r)
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
	context, key, value, err := parseRequest(w, r)
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
	context, key, value, err := parseRequest(w, r)
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
	
	err = deleteEntry(db, context, value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Create a response for success
    successResponse := map[string]string{
        "context": context,
        "key": key,
		"value": "",
    }

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(successResponse)
    return
}