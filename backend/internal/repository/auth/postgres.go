package auth

import (
	"context"
	"errors"
	"strings"

	authdomain "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) GetUserByEmail(email string) (authdomain.User, error) {
	ctx := context.Background()

	var user authdomain.User

	err := r.db.QueryRow(ctx, `
		SELECT id, name, email, role, password_hash, password_salt, status
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`, strings.TrimSpace(email)).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.PasswordHash,
		&user.PasswordSalt,
		&user.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return authdomain.User{}, errors.New("user tidak ditemukan")
	}
	if err != nil {
		return authdomain.User{}, err
	}

	return user, nil
}
