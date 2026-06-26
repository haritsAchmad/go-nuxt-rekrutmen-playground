package repository

import (
	"errors"
	"fmt"
	"strings"
	"math"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
)

var lowonganData = generateDummyLowongan(1000)

func generateDummyLowongan(total int) []domain.Lowongan {
	judulList := []string{
		"Staf Administrasi",
		"Backend Developer",
		"Frontend Developer",
		"QA Engineer",
		"Data Analyst",
		"DevOps Engineer",
		"UI/UX Designer",
		"System Analyst",
	}

	unitList := []string{
		"Direktorat SDM",
		"Direktorat Sistem Informasi",
		"Direktorat Keuangan",
		"Direktorat Akademik",
		"Direktorat Kemahasiswaan",
	}

	statusList := []string{
		"aktif",
		"nonaktif",
	}

	result := []domain.Lowongan{}

	for i := 1; i <= total; i++ {
		judul := judulList[(i-1)%len(judulList)]
		unit := unitList[(i-1)%len(unitList)]
		status := statusList[(i-1)%len(statusList)]

		result = append(result, domain.Lowongan{
			ID:     i,
			Judul:  fmt.Sprintf("%s %d", judul, i),
			Unit:   unit,
			Status: status,
		})
	}

	return result
}

func GetAllLowongan(filter domain.LowonganFilterRequest) domain.LowonganListResponse {
	filtered := []domain.Lowongan{}

	keyword := strings.ToLower(filter.Keyword)

	for _, lowongan := range lowonganData {
		matchKeyword := true
		matchStatus := true

		if keyword != "" {
			judul := strings.ToLower(lowongan.Judul)
			unit := strings.ToLower(lowongan.Unit)

			matchKeyword = strings.Contains(judul, keyword) || strings.Contains(unit, keyword)
		}

		if filter.Status != "" {
			matchStatus = lowongan.Status == filter.Status
		}

		if matchKeyword && matchStatus {
			filtered = append(filtered, lowongan)
		}
	}

	page := filter.Page
	limit := filter.Limit

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	total := len(filtered)
	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}

	if end > total {
		end = total
	}

	paginatedData := filtered[start:end]

	return domain.LowonganListResponse{
		Data: paginatedData,
		Meta: domain.LowonganPaginationMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}
}

func GetLowonganByID(id int) (domain.Lowongan, error) {
	for _, lowongan := range lowonganData {
		if lowongan.ID == id {
			return lowongan, nil
		}
	}

	return domain.Lowongan{}, errors.New("lowongan tidak ditemukan")
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

func UpdateLowongan(id int, updatedLowongan domain.Lowongan) (domain.Lowongan, error) {
	for index, lowongan := range lowonganData {
		if lowongan.ID == id {
			lowonganData[index].Judul = updatedLowongan.Judul
			lowonganData[index].Unit = updatedLowongan.Unit
			lowonganData[index].Status = updatedLowongan.Status

			return lowonganData[index], nil
		}
	}

	return domain.Lowongan{}, errors.New("lowongan tidak ditemukan")
}

func containsID(ids []int, id int) bool {
	for _, selectedID := range ids {
		if selectedID == id {
			return true
		}
	}

	return false
}

func BulkUpdateLowonganStatus(ids []int, status string) error {
	found := false

	for index, lowongan := range lowonganData {
		if containsID(ids, lowongan.ID) {
			lowonganData[index].Status = status
			found = true
		}
	}

	if !found {
		return errors.New("lowongan tidak ditemukan")
	}

	return nil
}

func BulkDeleteLowongan(ids []int) error {
	found := false
	newData := []domain.Lowongan{}

	for _, lowongan := range lowonganData {
		if containsID(ids, lowongan.ID) {
			found = true
			continue
		}

		newData = append(newData, lowongan)
	}

	if !found {
		return errors.New("lowongan tidak ditemukan")
	}

	lowonganData = newData

	return nil
}
