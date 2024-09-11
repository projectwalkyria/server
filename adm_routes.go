package main

import (
	"errors"
    "encoding/json"
    "io"
    "net/http"
)

func parseAdmContextRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	
	if err != nil {
		return "", err
	}

	defer r.Body.Close()

    var data map[string]interface{}

	err = json.Unmarshal(body, &data)
    if err != nil {
		return "", err
    }

	if len(data) == 1 {
		// Access the first key-value pair (manually in this case, as Go doesn't provide a direct way to access the first element in a map).
		for _, value := range data {
			// Assert that `value` is of type `string`. You can also handle other types based on your input.
			valueStr, ok := value.(string)
			if !ok {
				return "", errors.New("Value is not a string")
			}
			return valueStr, nil
		}
	}
	return "", errors.New("JSON body must have 1 key-value pair")
}

func parseAdmTokenRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	
	if err != nil {
		return "", err
	}

	defer r.Body.Close()

    var data map[string]interface{}

	err = json.Unmarshal(body, &data)
    if err != nil {
		return "", err
    }

	if len(data) == 1 {
		// Access the first key-value pair (manually in this case, as Go doesn't provide a direct way to access the first element in a map).
		for _, value := range data {
			// Assert that `value` is of type `string`. You can also handle other types based on your input.
			valueStr, ok := value.(string)
			if !ok {
				return "", errors.New("Value is not a string")
			}
			return valueStr, nil
		}
	}
	return "", errors.New("JSON body must have 1 key-value pair")
}

func parseAdmTokenGrantRequest(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
	body, err := io.ReadAll(r.Body)
	
	if err != nil {
		return "", "", "", err
	}

	defer r.Body.Close()

    var data map[string]interface{}

    // Unmarshal the JSON string into the map
    err = json.Unmarshal([]byte(body), &data)
    if err != nil {
        return "", "", "", err
    }

    // Type assert the "token" field as a string
    token, ok := data["token"].(string)
    if !ok {
		return "", "", "", errors.New("key token not defined")
    }

    // Type assert the "grant" field as a string
    grant, ok := data["grant"].(string)
    if !ok {
		return "", "", "", errors.New("key grant not defined")
    }

    // Type assert the "grant" field as a string
    context, ok := data["context"].(string)
    if !ok {
		return "", "", "", errors.New("key context not defined")
    }

	return token, grant, context, nil
}

func admContextPost(w http.ResponseWriter, r *http.Request) {
	context, err := parseAdmContextRequest(w, r)
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
	 
	context, err = createContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Create a response for success
    successResponse := map[string]string{
        "context": context,
    }

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(successResponse)
    return
}

func admContextGet(w http.ResponseWriter, r *http.Request) {
	context, err := parseAdmContextRequest(w, r)
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

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Create a response for success
    successResponse := map[string]string{
        "context": context,
    }

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(successResponse)
	return
}

func admContextDelete(w http.ResponseWriter, r *http.Request) {
	context, err := parseAdmContextRequest(w, r)
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
	
	err = deleteContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Create a response for success
    successResponse := map[string]string{
        "context": "",
    }

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(successResponse)
    return
}

func admTokenPost(w http.ResponseWriter, r *http.Request) {
 	// Connect to SQLite
	db, err := connectSQLite() // For SQLite
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()
	
	var token string
	token, err = createToken(db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Create a response for success
    successResponse := map[string]string{
        "token": token,
    }

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(successResponse)
    return
}

func admTokenDelete(w http.ResponseWriter, r *http.Request) {
	token, err := parseAdmTokenRequest(w, r)
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
	
	err = deleteToken(db, token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Create a response for success
    successResponse := map[string]string{
        "token": "",
    }

    // Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(successResponse)
    return
}

func admTokenGrant(w http.ResponseWriter, r *http.Request) {
	token, grant, context, err := parseAdmTokenGrantRequest(w, r)
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

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, grant, context, err = grantTokenPermission(db, token, grant, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a response for success
	successResponse := map[string]string{
		"token": token,
		"permission": grant,
		"context": context,
	}

	// Set header and return success response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successResponse)
	return
}

func admTokenRevoke(w http.ResponseWriter, r *http.Request) {
	token, grant, context, err := parseAdmTokenGrantRequest(w, r)
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

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rovokeTokenPermission(db, token, grant, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set header and return success response as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    return
}
