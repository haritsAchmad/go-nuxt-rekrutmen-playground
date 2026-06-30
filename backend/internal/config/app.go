package config

import "fmt"

type AppConfig struct {
	Port string
	Env  string
}

type AuthConfig struct {
	SecretKey string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Schema   string
	SSLMode  string
}

type Config struct {
	App      AppConfig
	Auth     AuthConfig
	Database DatabaseConfig
}

func Load() Config {
	return Config{
		App: AppConfig{
			Port: GetString("APP_PORT", "8080"),
			Env:  GetString("APP_ENV", "local"),
		},
		Auth: AuthConfig{
			SecretKey: GetString("AUTH_SECRET_KEY", "playground-secret-key"),
		},
		Database: DatabaseConfig{
			Host:     GetString("DB_HOST", "localhost"),
			Port:     GetInt("DB_PORT", 5432),
			User:     GetString("DB_USER", "postgres"),
			Password: GetString("DB_PASSWORD", "admin"),
			Name:     GetString("DB_NAME", "playground"),
			Schema:   GetString("DB_SCHEMA", "rekrutmen_playground"),
			SSLMode:  GetString("DB_SSLMODE", "disable"),
		},
	}
}

func (db DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		db.Host,
		db.Port,
		db.User,
		db.Password,
		db.Name,
		db.SSLMode,
		db.Schema,
	)
}
