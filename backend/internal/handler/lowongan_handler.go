// backend/internal/handler/lowongan_handler.go
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/usecase"
)

type LowonganHandler struct {
	usecase *usecase.LowonganUsecase
}

func NewLowonganHandler(usecase *usecase.LowonganUsecase) *LowonganHandler {
	return &LowonganHandler{
		usecase: usecase,
	}
}

type apiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, statusCode int, success bool, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(apiResponse{
		Success: success,
		Message: message,
		Data:    data,
	})
}

func success(w http.ResponseWriter, message string, data any) {
	writeJSON(w, http.StatusOK, true, message, data)
}

func created(w http.ResponseWriter, message string, data any) {
	writeJSON(w, http.StatusCreated, true, message, data)
}

func errorResponse(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, false, message, nil)
}

func getIDFromQuery(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("ID tidak valid")
	}

	return id, nil
}

func (h *LowonganHandler) LowonganHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetLowongan(w, r)
	case http.MethodPost:
		h.CreateLowongan(w, r)
	case http.MethodPut:
		h.UpdateLowongan(w, r)
	case http.MethodDelete:
		h.DeleteLowongan(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (h *LowonganHandler) GetLowongan(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	filter := domain.LowonganFilterRequest{
		Keyword: r.URL.Query().Get("keyword"),
		Status:  r.URL.Query().Get("status"),
		Page:    page,
		Limit:   limit,
	}

	data, err := h.usecase.GetLowonganList(filter)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Berhasil mengambil data lowongan", data)
}

func (h *LowonganHandler) GetLowonganDetail(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ID lowongan tidak valid")
		return
	}

	data, err := h.usecase.GetLowonganDetail(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	success(w, "Berhasil mengambil detail lowongan", data)
}

func (h *LowonganHandler) LowonganDetailHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetLowonganDetail(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (h *LowonganHandler) CreateLowongan(w http.ResponseWriter, r *http.Request) {
	var request domain.CreateLowonganRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	newLowongan, err := h.usecase.CreateLowongan(request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	created(w, "Lowongan berhasil ditambahkan", newLowongan)
}

func (h *LowonganHandler) DeleteLowongan(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ID lowongan tidak valid")
		return
	}

	err = h.usecase.DeleteLowongan(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	success(w, "Lowongan berhasil dihapus", nil)
}

func (h *LowonganHandler) UpdateLowonganStatus(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ID lowongan tidak valid")
		return
	}

	var request domain.UpdateLowonganStatusRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	updatedLowongan, err := h.usecase.UpdateLowonganStatus(id, request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Status lowongan berhasil diubah", updatedLowongan)
}

func (h *LowonganHandler) UpdateLowongan(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ID lowongan tidak valid")
		return
	}

	var request domain.UpdateLowonganRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	updatedLowongan, err := h.usecase.UpdateLowongan(id, request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Lowongan berhasil diubah", updatedLowongan)
}

func (h *LowonganHandler) LowonganStatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.UpdateLowonganStatus(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (h *LowonganHandler) BulkUpdateLowonganStatus(w http.ResponseWriter, r *http.Request) {
	var request domain.BulkUpdateStatusRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	err = h.usecase.BulkUpdateLowonganStatus(request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Status lowongan terpilih berhasil diubah", nil)
}

func (h *LowonganHandler) BulkDeleteLowongan(w http.ResponseWriter, r *http.Request) {
	var request domain.BulkDeleteRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	err = h.usecase.BulkDeleteLowongan(request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Lowongan terpilih berhasil dihapus", nil)
}

func (h *LowonganHandler) LowonganBulkStatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.BulkUpdateLowonganStatus(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (h *LowonganHandler) LowonganBulkDeleteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		h.BulkDeleteLowongan(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}
