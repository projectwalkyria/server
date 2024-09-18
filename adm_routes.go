package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func parseAdmContextRequest(r *http.Request) (string, error) {
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
				return "", errors.New("value is not a string")
			}
			return valueStr, nil
		}
	}
	return "", errors.New("JSON body must have 1 key-value pair")
}

func parseAdmTokenRequest(r *http.Request) (string, error) {
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
				return "", errors.New("value is not a string")
			}
			return valueStr, nil
		}
	}
	return "", errors.New("JSON body must have 1 key-value pair")
}

func parseAdmTokenGrantRequest(r *http.Request) (string, string, string, error) {
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

func validateGrant(grant string) error {
	allowedGrants := []string{"PUT", "POST", "GET", "DELETE"}

	// Variable to track if the value is found
	found := false

	// Iterate over the slice to check for the value
	for _, item := range allowedGrants {
		if item == grant {
			found = true
			break
		}
	}

	// Check if the value was not found
	if !found {
		return errors.New("Grant " + grant + " is not a valid permission.")
	}
	return nil
}

func admContextPost(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	context, err := parseAdmContextRequest(r)
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

	err = getPermission(db, authToken, "ALL", "ADM_CONTEXT_POST")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if len(context) < 8 {
		http.Error(w, errors.New("a context must to have at least 8 letters").Error(), http.StatusBadRequest)
		return
	}

	context, err = createContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	successResponse := map[string]string{
		"context": context,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successResponse)

}

func admContextGet(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	context, err := parseAdmContextRequest(r)
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

	err = getPermission(db, authToken, "ALL", "ADM_CONTEXT_GET")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	successResponse := map[string]string{
		"context": context,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)

}

func admContextDelete(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	context, err := parseAdmContextRequest(r)
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

	err = getPermission(db, authToken, "ALL", "ADM_CONTEXT_DELETE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = deleteContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func admTokenPost(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	db, err := connectSQLite()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	err = getPermission(db, authToken, "ALL", "ADM_TOKEN_POST")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var token string
	token, err = createToken(db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	successResponse := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successResponse)

}

func admTokenDelete(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := parseAdmTokenRequest(r)
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

	err = getPermission(db, authToken, "ALL", "ADM_TOKEN_DELETE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = deleteToken(db, token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func admTokenGrant(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, grant, context, err := parseAdmTokenGrantRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateGrant(grant)
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

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = getPermission(db, authToken, "ALL", "ADM_TOKEN_GRANT")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, grant, context, err = grantTokenPermission(db, token, grant, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	successResponse := map[string]string{
		"token":      token,
		"permission": grant,
		"context":    context,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successResponse)

}

func admTokenRevoke(w http.ResponseWriter, r *http.Request) {
	authToken, err := getHeaderAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, grant, context, err := parseAdmTokenGrantRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateGrant(grant)
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

	context, err = getContext(db, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = getPermission(db, authToken, "ALL", "ADM_TOKEN_REVOKE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = rovokeTokenPermission(db, token, grant, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
