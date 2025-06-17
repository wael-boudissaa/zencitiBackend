package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	response := map[string]interface{}{
		"status": status,
		"data":   v,
	}
	return json.NewEncoder(r).Encode(response)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, err.Error())
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
		} else {
			log.Printf("%s request to %s with body: %s\n", r.Method, r.URL.Path, string(bodyBytes))
			// Reset the body so it can be read again later
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		next.ServeHTTP(w, r)
	})
}
