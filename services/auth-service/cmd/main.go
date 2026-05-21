package main

import (
	"encoding/json"
	"net/http"
)

func healthhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/health", healthhandler)
	http.ListenAndServe(":8080", nil)
}
