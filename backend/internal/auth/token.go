package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Claims struct {
	UserID int
	Name   string
	Email  string
	Role   string
	Exp    int64
}

func ValidateToken(token string, secretKey string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, errors.New("token tidak valid")
	}

	unsignedToken := parts[0] + "." + parts[1]
	expectedSignature := sign(unsignedToken, secretKey)

	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return Claims{}, errors.New("signature token tidak valid")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, errors.New("payload token tidak valid")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return Claims{}, errors.New("payload token tidak valid")
	}

	claims := Claims{}

	if sub, ok := payload["sub"].(float64); ok {
		claims.UserID = int(sub)
	} else if sub, ok := payload["sub"].(string); ok {
		userID, _ := strconv.Atoi(sub)
		claims.UserID = userID
	}

	if name, ok := payload["name"].(string); ok {
		claims.Name = name
	}

	if email, ok := payload["email"].(string); ok {
		claims.Email = email
	}

	if role, ok := payload["role"].(string); ok {
		claims.Role = role
	}

	if exp, ok := payload["exp"].(float64); ok {
		claims.Exp = int64(exp)
	}

	if claims.UserID <= 0 || claims.Email == "" || claims.Role == "" {
		return Claims{}, errors.New("claims token tidak lengkap")
	}

	if claims.Exp <= time.Now().Unix() {
		return Claims{}, errors.New("token sudah kedaluwarsa")
	}

	return claims, nil
}

func sign(value string, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
