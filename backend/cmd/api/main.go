package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Go backend is running",
		})
	})

	http.HandleFunc("/api/lowongan", handler.LowonganHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}