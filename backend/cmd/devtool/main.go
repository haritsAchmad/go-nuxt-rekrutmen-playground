package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/database"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/config"
)

var allowedRoles = map[string]bool{
	"superadmin": true,
	"admin":      true,
	"viewer":     true,
}

func main() {
	command := flag.String("command", "create-user", "command yang dijalankan")
	name := flag.String("name", "", "nama user")
	email := flag.String("email", "", "email user")
	password := flag.String("password", "", "password user")
	role := flag.String("role", "viewer", "role user: superadmin, admin, viewer")
	flag.Parse()

	if *command != "create-user" {
		log.Fatalf("command tidak dikenal: %s", *command)
	}

	input := createUserInput{
		Name:     strings.TrimSpace(*name),
		Email:    strings.ToLower(strings.TrimSpace(*email)),
		Password: strings.TrimSpace(*password),
		Role:     strings.ToLower(strings.TrimSpace(*role)),
	}

	if err := input.Validate(); err != nil {
		log.Fatal(err)
	}

	cfg := config.Load()

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("failed connect database: ", err)
	}
	defer db.Close()

	salt, err := generateSalt(16)
	if err != nil {
		log.Fatal("failed generate salt: ", err)
	}

	passwordHash := hashPassword(salt, input.Password)

	var userID int
	err = db.QueryRow(context.Background(), `
		INSERT INTO users (name, email, role, password_hash, password_salt, status)
		VALUES ($1, $2, $3, $4, $5, 'aktif')
		RETURNING id
	`, input.Name, input.Email, input.Role, passwordHash, salt).Scan(&userID)
	if err != nil {
		log.Fatal("failed create user: ", err)
	}

	fmt.Println("✅ User berhasil dibuat")
	fmt.Println("ID   :", userID)
	fmt.Println("Nama :", input.Name)
	fmt.Println("Email:", input.Email)
	fmt.Println("Role :", input.Role)
}

type createUserInput struct {
	Name     string
	Email    string
	Password string
	Role     string
}

func (input createUserInput) Validate() error {
	if input.Name == "" {
		return fmt.Errorf("name wajib diisi")
	}

	if input.Email == "" || !strings.Contains(input.Email, "@") {
		return fmt.Errorf("email tidak valid")
	}

	if len(input.Password) < 6 {
		return fmt.Errorf("password minimal 6 karakter")
	}

	if !allowedRoles[input.Role] {
		return fmt.Errorf("role tidak valid. Pilihan: superadmin, admin, viewer")
	}

	return nil
}

func generateSalt(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func hashPassword(salt string, password string) string {
	hash := sha256.Sum256([]byte(salt + password))
	return hex.EncodeToString(hash[:])
}
