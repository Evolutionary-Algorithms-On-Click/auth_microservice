package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Validate and decode JSON data in the request to a map.
func Body(req *http.Request) (map[string]any, error) {
	if req.Method != "POST" {
		return nil, fmt.Errorf("%v not allowed", req.Method)
	}

	body := json.NewDecoder(req.Body)
	var data map[string]any
	if err := body.Decode(&data); err != nil {
		return nil, errors.New("invalid JSON body")
	}

	return data, nil
}

func FromJson[T any](data map[string]any) (*T, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result *T
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//skipping CSRF authentication for routes where user is not authenticated
		if r.URL.Path == "/api/login" || r.URL.Path == "/api/register" || r.URL.Path == "/api/password/reset" || r.URL.Path == "/api/password/verify" || r.URL.Path == "/api/verify" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("csrf_token")
		if err != nil {
			http.Error(w, "CSRF cookie not found", http.StatusForbidden)
			return
		}

		csrfToken := r.Header.Get("X-CSRF-Token")
		if csrfToken == "" || csrfToken != cookie.Value {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
