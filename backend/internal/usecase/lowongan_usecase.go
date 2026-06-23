package usecase

import (
	"errors"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/repository"
)

func GetLowonganList() []domain.Lowongan {
	return repository.GetAllLowongan()
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