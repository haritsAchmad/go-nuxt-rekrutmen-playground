package repository

import "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"

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

func GetAllLowongan() []domain.Lowongan {
	return lowonganData
}

func CreateLowongan(lowongan domain.Lowongan) domain.Lowongan {
	lowongan.ID = len(lowonganData) + 1
	lowonganData = append(lowonganData, lowongan)

	return lowongan
}