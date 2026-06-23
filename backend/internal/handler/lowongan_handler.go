// backend/internal/handler/lowongan_handler.go
package handler

import (
	"encoding/json"
	"net/http"
)

type Lowongan struct {
	ID     int    `json:"id"`
	Judul  string `json:"judul"`
	Unit   string `json:"unit"`
	Status string `json:"status"`
}

type CreateLowonganRequest struct {
	Judul string `json:"judul"`
	Unit  string `json:"unit"`
}

var lowonganData = []Lowongan{
	{
		ID:     1,
		Judul:  "Staf Administrasi",
		Unit:   "Direktorat SDM",
		Status: "aktif",
	},
	{
		ID:     2,
		Judul:  "Backend Developer",
		Unit:   "Direktorat Sistem Informasi",
		Status: "aktif",
	},
}

func LowonganHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetLowongan(w, r)
	case http.MethodPost:
		CreateLowongan(w, r)
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

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    lowonganData,
	})
}

func CreateLowongan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request CreateLowonganRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Format JSON tidak valid",
		})
		return
	}

	if request.Judul == "" || request.Unit == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Judul dan unit wajib diisi",
		})
		return
	}

	newLowongan := Lowongan{
		ID:     len(lowonganData) + 1,
		Judul:  request.Judul,
		Unit:   request.Unit,
		Status: "aktif",
	}

	lowonganData = append(lowonganData, newLowongan)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Lowongan berhasil ditambahkan",
		"data":    newLowongan,
	})
}