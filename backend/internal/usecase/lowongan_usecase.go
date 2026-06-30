package usecase

import (
	"errors"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/repository"
)

type LowonganRepository interface {
	GetAllLowongan(filter domain.LowonganFilterRequest) (domain.LowonganListResponse, error)
	GetLowonganByID(id int) (domain.Lowongan, error)
	CreateLowongan(lowongan domain.Lowongan) (domain.Lowongan, error)
	DeleteLowongan(id int) error
	UpdateLowonganStatus(id int, status string) (domain.Lowongan, error)
	UpdateLowongan(id int, lowongan domain.Lowongan) (domain.Lowongan, error)
	BulkUpdateLowonganStatus(ids []int, status string) error
	BulkDeleteLowongan(ids []int) error
}

var lowonganRepository LowonganRepository = repository.NewLowonganMemoryRepository()

func SetLowonganRepository(repo LowonganRepository) {
	if repo == nil {
		return
	}

	lowonganRepository = repo
}

func GetLowonganList(filter domain.LowonganFilterRequest) (domain.LowonganListResponse, error) {
	if filter.Status != "" && filter.Status != "aktif" && filter.Status != "nonaktif" {
		return domain.LowonganListResponse{}, errors.New("status filter harus aktif atau nonaktif")
	}

	return lowonganRepository.GetAllLowongan(filter)
}

func GetLowonganDetail(id int) (domain.Lowongan, error) {
	if id <= 0 {
		return domain.Lowongan{}, errors.New("ID lowongan tidak valid")
	}

	return lowonganRepository.GetLowonganByID(id)
}

func CreateLowongan(request domain.CreateLowonganRequest) (domain.Lowongan, error) {
	if request.Judul == "" || request.Unit == "" {
		return domain.Lowongan{}, errors.New("judul dan unit wajib diisi")
	}

	lowongan := domain.Lowongan{
		Judul:  request.Judul,
		Unit:   request.Unit,
		Status: "aktif",
	}

	return lowonganRepository.CreateLowongan(lowongan)
}

func DeleteLowongan(id int) error {
	return lowonganRepository.DeleteLowongan(id)
}

func UpdateLowonganStatus(id int, request domain.UpdateLowonganStatusRequest) (domain.Lowongan, error) {
	if request.Status != "aktif" && request.Status != "nonaktif" {
		return domain.Lowongan{}, errors.New("status harus aktif atau nonaktif")
	}

	return lowonganRepository.UpdateLowonganStatus(id, request.Status)
}

func UpdateLowongan(id int, request domain.UpdateLowonganRequest) (domain.Lowongan, error) {
	if request.Judul == "" || request.Unit == "" {
		return domain.Lowongan{}, errors.New("judul dan unit wajib diisi")
	}

	if request.Status != "aktif" && request.Status != "nonaktif" {
		return domain.Lowongan{}, errors.New("status harus aktif atau nonaktif")
	}

	lowongan := domain.Lowongan{
		Judul:  request.Judul,
		Unit:   request.Unit,
		Status: request.Status,
	}

	return lowonganRepository.UpdateLowongan(id, lowongan)
}

func BulkUpdateLowonganStatus(request domain.BulkUpdateStatusRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu lowongan")
	}

	if request.Status != "aktif" && request.Status != "nonaktif" {
		return errors.New("status harus aktif atau nonaktif")
	}

	return lowonganRepository.BulkUpdateLowonganStatus(request.IDs, request.Status)
}

func BulkDeleteLowongan(request domain.BulkDeleteRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu lowongan")
	}

	return lowonganRepository.BulkDeleteLowongan(request.IDs)
}
