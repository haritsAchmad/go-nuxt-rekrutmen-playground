package main

import (
	"log"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/config"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/route"
)

func main() {
	cfg := config.Load()

	route.RegisterRoutes()

	address := ":" + cfg.App.Port
	log.Println("Server running on http://localhost" + address)
	log.Fatal(http.ListenAndServe(address, nil))
}
