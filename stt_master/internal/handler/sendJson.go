package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (e *HttpEnv) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка кодирования в JSON: %v", err)
	}
}

func (e *HttpEnv) sendJSONError(w http.ResponseWriter, statusCode int, errorMessage string) {
	e.sendJSON(w, statusCode, map[string]string{
		"message": errorMessage,
		"status":  "error",
	})
}
