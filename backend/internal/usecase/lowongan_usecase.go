package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
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

type LowonganUsecase struct {
	repo LowonganRepository
}

func NewLowonganUsecase(repo LowonganRepository) *LowonganUsecase {
	return &LowonganUsecase{
		repo: repo,
	}
}

func parseOptionalDate(value string, fieldName string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	parsedDate, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, errors.New(fieldName + " harus berformat YYYY-MM-DD")
	}

	return &parsedDate, nil
}

func validateLowonganDates(tanggalBuka *time.Time, tanggalTutup *time.Time) error {
	if tanggalBuka != nil && tanggalTutup != nil && tanggalBuka.After(*tanggalTutup) {
		return errors.New("tanggal buka tidak boleh melewati tanggal tutup")
	}

	return nil
}

func (u *LowonganUsecase) GetLowonganList(filter domain.LowonganFilterRequest) (domain.LowonganListResponse, error) {
	if filter.Status != "" && filter.Status != "aktif" && filter.Status != "nonaktif" {
		return domain.LowonganListResponse{}, errors.New("status filter harus aktif atau nonaktif")
	}

	return u.repo.GetAllLowongan(filter)
}

func (u *LowonganUsecase) GetLowonganDetail(id int) (domain.Lowongan, error) {
	if id <= 0 {
		return domain.Lowongan{}, errors.New("ID lowongan tidak valid")
	}

	return u.repo.GetLowonganByID(id)
}

func (u *LowonganUsecase) CreateLowongan(request domain.CreateLowonganRequest) (domain.Lowongan, error) {
	if request.Judul == "" || request.Unit == "" {
		return domain.Lowongan{}, errors.New("judul dan unit wajib diisi")
	}

	tanggalBuka, err := parseOptionalDate(request.TanggalBuka, "tanggal buka")
	if err != nil {
		return domain.Lowongan{}, err
	}

	tanggalTutup, err := parseOptionalDate(request.TanggalTutup, "tanggal tutup")
	if err != nil {
		return domain.Lowongan{}, err
	}

	if err := validateLowonganDates(tanggalBuka, tanggalTutup); err != nil {
		return domain.Lowongan{}, err
	}

	lowongan := domain.Lowongan{
		Judul:        strings.TrimSpace(request.Judul),
		Unit:         strings.TrimSpace(request.Unit),
		TanggalBuka:  tanggalBuka,
		TanggalTutup: tanggalTutup,
		Deskripsi:    strings.TrimSpace(request.Deskripsi),
		Status:       "aktif",
	}

	return u.repo.CreateLowongan(lowongan)
}

func (u *LowonganUsecase) DeleteLowongan(id int) error {
	return u.repo.DeleteLowongan(id)
}

func (u *LowonganUsecase) UpdateLowonganStatus(id int, request domain.UpdateLowonganStatusRequest) (domain.Lowongan, error) {
	if request.Status != "aktif" && request.Status != "nonaktif" {
		return domain.Lowongan{}, errors.New("status harus aktif atau nonaktif")
	}

	return u.repo.UpdateLowonganStatus(id, request.Status)
}

func (u *LowonganUsecase) UpdateLowongan(id int, request domain.UpdateLowonganRequest) (domain.Lowongan, error) {
	if request.Judul == "" || request.Unit == "" {
		return domain.Lowongan{}, errors.New("judul dan unit wajib diisi")
	}

	if request.Status != "aktif" && request.Status != "nonaktif" {
		return domain.Lowongan{}, errors.New("status harus aktif atau nonaktif")
	}

	tanggalBuka, err := parseOptionalDate(request.TanggalBuka, "tanggal buka")
	if err != nil {
		return domain.Lowongan{}, err
	}

	tanggalTutup, err := parseOptionalDate(request.TanggalTutup, "tanggal tutup")
	if err != nil {
		return domain.Lowongan{}, err
	}

	if err := validateLowonganDates(tanggalBuka, tanggalTutup); err != nil {
		return domain.Lowongan{}, err
	}

	lowongan := domain.Lowongan{
		Judul:        strings.TrimSpace(request.Judul),
		Unit:         strings.TrimSpace(request.Unit),
		TanggalBuka:  tanggalBuka,
		TanggalTutup: tanggalTutup,
		Deskripsi:    strings.TrimSpace(request.Deskripsi),
		Status:       request.Status,
	}

	return u.repo.UpdateLowongan(id, lowongan)
}

func (u *LowonganUsecase) BulkUpdateLowonganStatus(request domain.BulkUpdateStatusRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu lowongan")
	}

	if request.Status != "aktif" && request.Status != "nonaktif" {
		return errors.New("status harus aktif atau nonaktif")
	}

	return u.repo.BulkUpdateLowonganStatus(request.IDs, request.Status)
}

func (u *LowonganUsecase) BulkDeleteLowongan(request domain.BulkDeleteRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu lowongan")
	}

	return u.repo.BulkDeleteLowongan(request.IDs)
}
