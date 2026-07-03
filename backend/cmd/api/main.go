package main

import (
	"log"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/database"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/config"
	authhandler "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler/auth"
	lowonganhandler "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler/lowongan"
	pelamarhandler "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler/pelamar"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/middleware"
	authrepository "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/repository/auth"
	lowonganrepository "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/repository/lowongan"
	pelamarrepository "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/repository/pelamar"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/route"
	authusecase "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/usecase/auth"
	lowonganusecase "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/usecase/lowongan"
	pelamarusecase "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/usecase/pelamar"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("✅ PostgreSQL Connected")

	lowonganRepo := lowonganrepository.NewLowonganRepository(db)
	lowonganUsecase := lowonganusecase.NewLowonganUsecase(lowonganRepo)
	lowonganHandler := lowonganhandler.NewLowonganHandler(lowonganUsecase)

	pelamarRepo := pelamarrepository.NewPelamarRepository(db)
	pelamarUsecase := pelamarusecase.NewPelamarUsecase(pelamarRepo)
	pelamarHandler := pelamarhandler.NewPelamarHandler(pelamarUsecase)

	authRepo := authrepository.NewAuthRepository(db)
	authUsecase := authusecase.NewAuthUsecase(authRepo, cfg.Auth.SecretKey)
	authHandler := authhandler.NewAuthHandler(authUsecase)

	route.RegisterRoutes(lowonganHandler, pelamarHandler, authHandler, cfg.Auth.SecretKey)

	serverHandler := middleware.Recovery(
		middleware.Logger(
			middleware.CORS(http.DefaultServeMux),
		),
	)

	address := ":" + cfg.App.Port
	log.Println("Server running on http://localhost" + address)
	log.Fatal(http.ListenAndServe(address, serverHandler))
}
