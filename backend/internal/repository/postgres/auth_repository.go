package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
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

func (r *AuthRepository) GetUserByEmail(email string) (domain.User, error) {
	ctx := context.Background()

	var user domain.User

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
		return domain.User{}, errors.New("user tidak ditemukan")
	}
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
