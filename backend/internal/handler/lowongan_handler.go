// backend/internal/handler/lowongan_handler.go
package handler

import (
	"encoding/json"
	"net/http"
)

func GetLowongan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := []map[string]interface{}{
		{
			"id":     1,
			"judul":  "Staf Administrasi",
			"unit":   "Direktorat SDM",
			"status": "aktif",
		},
		{
			"id":     2,
			"judul":  "Backend Developer",
			"unit":   "Direktorat Sistem Informasi",
			"status": "aktif",
		},
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}