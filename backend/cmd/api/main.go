package main

import (
	"log"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/route"
)

func main() {
	route.RegisterRoutes()

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}