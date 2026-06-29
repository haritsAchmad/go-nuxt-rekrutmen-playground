package main

import (
	"log"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/database"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/config"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/route"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer db.Close()

	log.Println("✅ PostgreSQL Connected")

	route.RegisterRoutes()

	address := ":" + cfg.App.Port
	log.Println("Server running on http://localhost" + address)
	log.Fatal(http.ListenAndServe(address, nil))
}
