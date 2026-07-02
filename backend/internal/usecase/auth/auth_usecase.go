package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	authdomain "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain/auth"
)

type AuthRepository interface {
	GetUserByEmail(email string) (authdomain.User, error)
}

type AuthUsecase struct {
	repo      AuthRepository
	secretKey string
}

func NewAuthUsecase(repo AuthRepository, secretKey string) *AuthUsecase {
	return &AuthUsecase{
		repo:      repo,
		secretKey: secretKey,
	}
}

func (u *AuthUsecase) Login(request authdomain.LoginRequest) (authdomain.LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	password := strings.TrimSpace(request.Password)

	if email == "" || password == "" {
		return authdomain.LoginResponse{}, errors.New("email dan password wajib diisi")
	}

	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return authdomain.LoginResponse{}, errors.New("email atau password salah")
	}

	if user.Status != "aktif" {
		return authdomain.LoginResponse{}, errors.New("user tidak aktif")
	}

	if hashPassword(user.PasswordSalt, password) != user.PasswordHash {
		return authdomain.LoginResponse{}, errors.New("email atau password salah")
	}

	authUser := authdomain.AuthUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	token, err := u.generateToken(authUser)
	if err != nil {
		return authdomain.LoginResponse{}, err
	}

	return authdomain.LoginResponse{
		Token: token,
		User:  authUser,
	}, nil
}

func hashPassword(salt string, password string) string {
	hash := sha256.Sum256([]byte(salt + password))
	return hex.EncodeToString(hash[:])
}

func (u *AuthUsecase) generateToken(user authdomain.AuthUserResponse) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	payload := map[string]interface{}{
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)
	unsignedToken := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)

	signature := sign(unsignedToken, u.secretKey)

	return fmt.Sprintf("%s.%s", unsignedToken, signature), nil
}

func sign(value string, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
