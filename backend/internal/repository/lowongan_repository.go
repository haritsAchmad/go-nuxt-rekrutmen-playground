package repository

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

func GetAllLowongan() []domain.Lowongan {
	return lowonganData
}

func CreateLowongan(lowongan domain.Lowongan) domain.Lowongan {
	maxID := 0

	for _, item := range lowonganData {
		if item.ID > maxID {
			maxID = item.ID
		}
	}

	lowongan.ID = maxID + 1
	lowonganData = append(lowonganData, lowongan)

	return lowongan
}

func DeleteLowongan(id int) error {
	for index, lowongan := range lowonganData {
		if lowongan.ID == id {
			lowonganData = append(lowonganData[:index], lowonganData[index+1:]...)
			return nil
		}
	}

	return errors.New("lowongan tidak ditemukan")
}

func UpdateLowonganStatus(id int, status string) (domain.Lowongan, error) {
	for index, lowongan := range lowonganData {
		if lowongan.ID == id {
			lowonganData[index].Status = status
			return lowonganData[index], nil
		}
	}

	return domain.Lowongan{}, errors.New("lowongan tidak ditemukan")
}
