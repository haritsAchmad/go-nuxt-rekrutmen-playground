package pelamar

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	domain "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain/pelamar"
	handlerresponse "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler"
	pelamarusecase "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/usecase/pelamar"
)

type PelamarHandler struct {
	usecase *pelamarusecase.PelamarUsecase
}

func NewPelamarHandler(usecase *pelamarusecase.PelamarUsecase) *PelamarHandler {
	return &PelamarHandler{usecase: usecase}
}

func success(w http.ResponseWriter, message string, data any) {
	handlerresponse.Success(w, message, data)
}

func created(w http.ResponseWriter, message string, data any) {
	handlerresponse.Created(w, message, data)
}

func errorResponse(w http.ResponseWriter, statusCode int, message string) {
	handlerresponse.Error(w, statusCode, message)
}

func getIDFromQuery(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("ID tidak valid")
	}

	return id, nil
}

func (h *PelamarHandler) PelamarHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetPelamar(w, r)
	case http.MethodPost:
		h.CreatePelamar(w, r)
	case http.MethodPut:
		h.UpdatePelamar(w, r)
	case http.MethodDelete:
		h.DeletePelamar(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (h *PelamarHandler) GetPelamar(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	lowonganID, _ := strconv.Atoi(r.URL.Query().Get("lowonganId"))

	filter := domain.PelamarFilterRequest{
		Keyword:    r.URL.Query().Get("keyword"),
		Status:     r.URL.Query().Get("status"),
		Sort:       r.URL.Query().Get("sort"),
		LowonganID: lowonganID,
		Page:       page,
		Limit:      limit,
	}

	data, err := h.usecase.GetPelamarList(filter)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Berhasil mengambil data pelamar", data)
}

func (h *PelamarHandler) GetPelamarDetail(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ID pelamar tidak valid")
		return
	}

	data, err := h.usecase.GetPelamarDetail(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	success(w, "Berhasil mengambil detail pelamar", data)
}

func (h *PelamarHandler) PelamarDetailHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetPelamarDetail(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (h *PelamarHandler) CreatePelamar(w http.ResponseWriter, r *http.Request) {
	var request domain.CreatePelamarRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	newPelamar, err := h.usecase.CreatePelamar(request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	created(w, "Pelamar berhasil ditambahkan", newPelamar)
}

func (h *PelamarHandler) UpdatePelamar(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ID pelamar tidak valid")
		return
	}

	var request domain.UpdatePelamarRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	updatedPelamar, err := h.usecase.UpdatePelamar(id, request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Pelamar berhasil diubah", updatedPelamar)
}

func (h *PelamarHandler) DeletePelamar(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ID pelamar tidak valid")
		return
	}

	err = h.usecase.DeletePelamar(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	success(w, "Pelamar berhasil dihapus", nil)
}

func (h *PelamarHandler) UpdatePelamarStatus(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ID pelamar tidak valid")
		return
	}

	var request domain.UpdatePelamarStatusRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	updatedPelamar, err := h.usecase.UpdatePelamarStatus(id, request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Status pelamar berhasil diubah", updatedPelamar)
}

func (h *PelamarHandler) PelamarStatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.UpdatePelamarStatus(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (h *PelamarHandler) BulkUpdatePelamarStatus(w http.ResponseWriter, r *http.Request) {
	var request domain.BulkUpdateStatusRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	err = h.usecase.BulkUpdatePelamarStatus(request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Status pelamar terpilih berhasil diubah", nil)
}

func (h *PelamarHandler) BulkDeletePelamar(w http.ResponseWriter, r *http.Request) {
	var request domain.BulkDeleteRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	err = h.usecase.BulkDeletePelamar(request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	success(w, "Pelamar terpilih berhasil dihapus", nil)
}

func (h *PelamarHandler) PelamarBulkStatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.BulkUpdatePelamarStatus(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (h *PelamarHandler) PelamarBulkDeleteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		h.BulkDeletePelamar(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}
