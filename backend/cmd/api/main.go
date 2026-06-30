package main

import (
	"log"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/database"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/config"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/middleware"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/repository/postgres"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/route"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/usecase"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("✅ PostgreSQL Connected")

	lowonganRepo := postgres.NewLowonganRepository(db)
	lowonganUsecase := usecase.NewLowonganUsecase(lowonganRepo)
	lowonganHandler := handler.NewLowonganHandler(lowonganUsecase)

	route.RegisterRoutes(lowonganHandler)

	serverHandler := middleware.Recovery(
		middleware.Logger(
			middleware.CORS(http.DefaultServeMux),
		),
	)

	address := ":" + cfg.App.Port
	log.Println("Server running on http://localhost" + address)
	log.Fatal(http.ListenAndServe(address, serverHandler))
}
