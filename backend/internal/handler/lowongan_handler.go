// backend/internal/handler/lowongan_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/request"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/response"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/usecase"
)

func LowonganHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetLowongan(w, r)
	case http.MethodPost:
		CreateLowongan(w, r)
	case http.MethodPut:
		UpdateLowongan(w, r)
	case http.MethodDelete:
		DeleteLowongan(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func GetLowongan(w http.ResponseWriter, r *http.Request) {
	filter := domain.LowonganFilterRequest{
		Keyword: r.URL.Query().Get("keyword"),
		Status:  r.URL.Query().Get("status"),
	}

	data, err := usecase.GetLowonganList(filter)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, "Berhasil mengambil data lowongan", data)
}

func CreateLowongan(w http.ResponseWriter, r *http.Request) {
	var request domain.CreateLowonganRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	newLowongan, err := usecase.CreateLowongan(request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, "Lowongan berhasil ditambahkan", newLowongan)
}

func DeleteLowongan(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIDFromQuery(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID lowongan tidak valid")
		return
	}

	err = usecase.DeleteLowongan(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, "Lowongan berhasil dihapus", nil)
}

func UpdateLowonganStatus(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIDFromQuery(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID lowongan tidak valid")
		return
	}

	var request domain.UpdateLowonganStatusRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	updatedLowongan, err := usecase.UpdateLowonganStatus(id, request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, "Status lowongan berhasil diubah", updatedLowongan)
}

func UpdateLowongan(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetIDFromQuery(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID lowongan tidak valid")
		return
	}

	var request domain.UpdateLowonganRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	updatedLowongan, err := usecase.UpdateLowongan(id, request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, "Lowongan berhasil diubah", updatedLowongan)
}

func LowonganStatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		UpdateLowonganStatus(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}