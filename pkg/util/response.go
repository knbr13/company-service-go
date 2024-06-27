package util

import (
	"encoding/json"
	"log"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Failed to encode payload: %s\n", err.Error())
		return
	}
}

func ErrJsonResponse(w http.ResponseWriter, statusCode int, errMsg string) {
	JsonResponse(w, statusCode, map[string]any{
		"error": errMsg,
	})
}
