package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJson(r *http.Request, v any) error {

	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func ParseJsonList(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	defer r.Body.Close()

	// Decode the entire JSON array into the provided slice.
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	return nil
}

func WriteJson(r http.ResponseWriter, status int, v any) error {
	r.Header().Add("Content-Type", "application/json")
	r.WriteHeader(status)
	return json.NewEncoder(r).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, map[string]string{"error": err.Error()})
}
