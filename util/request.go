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
