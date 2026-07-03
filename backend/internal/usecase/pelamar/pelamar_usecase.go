package pelamar

import (
	"errors"
	"strings"

	pelamardomain "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain/pelamar"
)

type PelamarRepository interface {
	GetAllPelamar(filter pelamardomain.PelamarFilterRequest) (pelamardomain.PelamarListResponse, error)
	GetPelamarByID(id int) (pelamardomain.Pelamar, error)
	CreatePelamar(pelamar pelamardomain.Pelamar) (pelamardomain.Pelamar, error)
	UpdatePelamar(id int, pelamar pelamardomain.Pelamar) (pelamardomain.Pelamar, error)
	UpdatePelamarStatus(id int, status string) (pelamardomain.Pelamar, error)
	DeletePelamar(id int) error
	BulkUpdatePelamarStatus(ids []int, status string) error
	BulkDeletePelamar(ids []int) error
}

type PelamarUsecase struct {
	repo PelamarRepository
}

func NewPelamarUsecase(repo PelamarRepository) *PelamarUsecase {
	return &PelamarUsecase{repo: repo}
}

func isValidStatus(status string) bool {
	return status == "baru" || status == "diproses" || status == "diterima" || status == "ditolak"
}

func (u *PelamarUsecase) GetPelamarList(filter pelamardomain.PelamarFilterRequest) (pelamardomain.PelamarListResponse, error) {
	if filter.Status != "" && !isValidStatus(filter.Status) {
		return pelamardomain.PelamarListResponse{}, errors.New("status filter harus baru, diproses, diterima, atau ditolak")
	}

	if filter.Sort != "" && filter.Sort != "newest" && filter.Sort != "oldest" {
		return pelamardomain.PelamarListResponse{}, errors.New("sort harus newest atau oldest")
	}

	return u.repo.GetAllPelamar(filter)
}

func (u *PelamarUsecase) GetPelamarDetail(id int) (pelamardomain.Pelamar, error) {
	if id <= 0 {
		return pelamardomain.Pelamar{}, errors.New("ID pelamar tidak valid")
	}

	return u.repo.GetPelamarByID(id)
}

func (u *PelamarUsecase) CreatePelamar(request pelamardomain.CreatePelamarRequest) (pelamardomain.Pelamar, error) {
	if request.LowonganID <= 0 {
		return pelamardomain.Pelamar{}, errors.New("lowongan wajib dipilih")
	}

	if strings.TrimSpace(request.Nama) == "" || strings.TrimSpace(request.Email) == "" {
		return pelamardomain.Pelamar{}, errors.New("nama dan email wajib diisi")
	}

	pelamar := pelamardomain.Pelamar{
		LowonganID: request.LowonganID,
		Nama:       strings.TrimSpace(request.Nama),
		Email:      strings.TrimSpace(request.Email),
		NoHP:       strings.TrimSpace(request.NoHP),
		Status:     "baru",
	}

	return u.repo.CreatePelamar(pelamar)
}

func (u *PelamarUsecase) UpdatePelamar(id int, request pelamardomain.UpdatePelamarRequest) (pelamardomain.Pelamar, error) {
	if id <= 0 {
		return pelamardomain.Pelamar{}, errors.New("ID pelamar tidak valid")
	}

	if request.LowonganID <= 0 {
		return pelamardomain.Pelamar{}, errors.New("lowongan wajib dipilih")
	}

	if strings.TrimSpace(request.Nama) == "" || strings.TrimSpace(request.Email) == "" {
		return pelamardomain.Pelamar{}, errors.New("nama dan email wajib diisi")
	}

	if !isValidStatus(request.Status) {
		return pelamardomain.Pelamar{}, errors.New("status harus baru, diproses, diterima, atau ditolak")
	}

	pelamar := pelamardomain.Pelamar{
		LowonganID: request.LowonganID,
		Nama:       strings.TrimSpace(request.Nama),
		Email:      strings.TrimSpace(request.Email),
		NoHP:       strings.TrimSpace(request.NoHP),
		Status:     request.Status,
	}

	return u.repo.UpdatePelamar(id, pelamar)
}

func (u *PelamarUsecase) UpdatePelamarStatus(id int, request pelamardomain.UpdatePelamarStatusRequest) (pelamardomain.Pelamar, error) {
	if id <= 0 {
		return pelamardomain.Pelamar{}, errors.New("ID pelamar tidak valid")
	}

	if !isValidStatus(request.Status) {
		return pelamardomain.Pelamar{}, errors.New("status harus baru, diproses, diterima, atau ditolak")
	}

	return u.repo.UpdatePelamarStatus(id, request.Status)
}

func (u *PelamarUsecase) DeletePelamar(id int) error {
	if id <= 0 {
		return errors.New("ID pelamar tidak valid")
	}

	return u.repo.DeletePelamar(id)
}

func (u *PelamarUsecase) BulkUpdatePelamarStatus(request pelamardomain.BulkUpdateStatusRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu pelamar")
	}

	if !isValidStatus(request.Status) {
		return errors.New("status harus baru, diproses, diterima, atau ditolak")
	}

	return u.repo.BulkUpdatePelamarStatus(request.IDs, request.Status)
}

func (u *PelamarUsecase) BulkDeletePelamar(request pelamardomain.BulkDeleteRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu pelamar")
	}

	return u.repo.BulkDeletePelamar(request.IDs)
}
