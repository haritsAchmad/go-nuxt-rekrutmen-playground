package usecase

import (
	"errors"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/repository"
)

func GetLowonganList(filter domain.LowonganFilterRequest) ([]domain.Lowongan, error) {
	if filter.Status != "" && filter.Status != "aktif" && filter.Status != "nonaktif" {
		return []domain.Lowongan{}, errors.New("status filter harus aktif atau nonaktif")
	}

	data := repository.GetAllLowongan(filter)

	return data, nil
}

func GetLowonganDetail(id int) (domain.Lowongan, error) {
	if id <= 0 {
		return domain.Lowongan{}, errors.New("ID lowongan tidak valid")
	}

	return repository.GetLowonganByID(id)
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

	newLowongan := repository.CreateLowongan(lowongan)

	return newLowongan, nil
}

func DeleteLowongan(id int) error {
	err := repository.DeleteLowongan(id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateLowonganStatus(id int, request domain.UpdateLowonganStatusRequest) (domain.Lowongan, error) {
	if request.Status != "aktif" && request.Status != "nonaktif" {
		return domain.Lowongan{}, errors.New("status harus aktif atau nonaktif")
	}

	updatedLowongan, err := repository.UpdateLowonganStatus(id, request.Status)
	if err != nil {
		return domain.Lowongan{}, err
	}

	return updatedLowongan, nil
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

	updatedLowongan, err := repository.UpdateLowongan(id, lowongan)
	if err != nil {
		return domain.Lowongan{}, err
	}

	return updatedLowongan, nil
}

func BulkUpdateLowonganStatus(request domain.BulkUpdateStatusRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu lowongan")
	}

	if request.Status != "aktif" && request.Status != "nonaktif" {
		return errors.New("status harus aktif atau nonaktif")
	}

	return repository.BulkUpdateLowonganStatus(request.IDs, request.Status)
}

func BulkDeleteLowongan(request domain.BulkDeleteRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu lowongan")
	}

	return repository.BulkDeleteLowongan(request.IDs)
}