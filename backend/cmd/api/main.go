package main

import (
	"log"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/database"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/config"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/repository/postgres"
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

	repo := postgres.NewLowonganRepository(db)

	// result, err := repo.GetAllLowongan(domain.LowonganFilterRequest{})
	result, err := repo.GetAllLowongan(domain.LowonganFilterRequest{
		Page:  2,
		Limit: 5,
	})
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("Total data: %d\n", len(result.Data))
	for _, item := range result.Data {
		log.Println(item.ID, item.Judul)
	}

	route.RegisterRoutes()

	address := ":" + cfg.App.Port
	log.Println("Server running on http://localhost" + address)
	log.Fatal(http.ListenAndServe(address, nil))
}
