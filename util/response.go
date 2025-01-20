package util

import (
	"encoding/json"
	"net/http"
)

// Response is a struct for holding the response data.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSONResponse writes a JSON response to the http.ResponseWriter.
func JSONResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	response := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
