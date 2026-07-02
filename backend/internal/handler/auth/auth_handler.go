package auth

import (
	"encoding/json"
	"net/http"

	authdomain "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain/auth"
	handlerresponse "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler"
	authusecase "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/usecase/auth"
)

type AuthHandler struct {
	usecase *authusecase.AuthUsecase
}

func NewAuthHandler(usecase *authusecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Login(w, r)
	default:
		handlerresponse.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request authdomain.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handlerresponse.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	data, err := h.usecase.Login(request)
	if err != nil {
		handlerresponse.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	handlerresponse.Success(w, "Login berhasil", data)
}
