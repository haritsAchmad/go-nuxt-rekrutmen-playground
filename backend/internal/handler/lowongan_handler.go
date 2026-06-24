// backend/internal/handler/lowongan_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/usecase"
)

func LowonganHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetLowongan(w, r)
	case http.MethodPost:
		CreateLowongan(w, r)
	case http.MethodDelete:
		DeleteLowongan(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method tidak diizinkan",
		})
	}
}

func GetLowongan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := usecase.GetLowonganList()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func CreateLowongan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request domain.CreateLowonganRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Format JSON tidak valid",
		})
		return
	}

	newLowongan, err := usecase.CreateLowongan(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Lowongan berhasil ditambahkan",
		"data":    newLowongan,
	})
}

func DeleteLowongan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idText := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idText)

	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "ID lowongan tidak valid",
		})
		return
	}

	err = usecase.DeleteLowongan(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Lowongan berhasil dihapus",
	})
}