package usecase

import (
	"errors"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
)

var lowonganData = []domain.Lowongan{
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

func GetLowonganList() []domain.Lowongan {
	return lowonganData
}

func CreateLowongan(request domain.CreateLowonganRequest) (domain.Lowongan, error) {
	if request.Judul == "" || request.Unit == "" {
		return domain.Lowongan{}, errors.New("judul dan unit wajib diisi")
	}

	newLowongan := domain.Lowongan{
		ID:     len(lowonganData) + 1,
		Judul:  request.Judul,
		Unit:   request.Unit,
		Status: "aktif",
	}

	lowonganData = append(lowonganData, newLowongan)

	return newLowongan, nil
}