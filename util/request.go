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
