package utils

import (
	"encoding/json"
	"net/http"
)

func SetError(err error, message string) map[string]string {
	if err != nil {
		return map[string]string{"error": message, "details": err.Error()}
	}
	return map[string]string{"error": message}
}

func WriteJSONError(w http.ResponseWriter, status int, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(SetError(err, message))
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
