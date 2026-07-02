package lowongan

import (
	"errors"
	"strings"
	"time"

	lowongandomain "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain/lowongan"
)

type LowonganRepository interface {
	GetAllLowongan(filter lowongandomain.LowonganFilterRequest) (lowongandomain.LowonganListResponse, error)
	GetLowonganByID(id int) (lowongandomain.Lowongan, error)
	CreateLowongan(lowongan lowongandomain.Lowongan) (lowongandomain.Lowongan, error)
	DeleteLowongan(id int) error
	UpdateLowonganStatus(id int, status string) (lowongandomain.Lowongan, error)
	UpdateLowongan(id int, lowongan lowongandomain.Lowongan) (lowongandomain.Lowongan, error)
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

func (u *LowonganUsecase) GetLowonganList(filter lowongandomain.LowonganFilterRequest) (lowongandomain.LowonganListResponse, error) {
	if filter.Status != "" && filter.Status != "aktif" && filter.Status != "nonaktif" {
		return lowongandomain.LowonganListResponse{}, errors.New("status filter harus aktif atau nonaktif")
	}

	if filter.Sort != "" && filter.Sort != "newest" && filter.Sort != "oldest" {
		return lowongandomain.LowonganListResponse{}, errors.New("sort harus newest atau oldest")
	}

	return u.repo.GetAllLowongan(filter)
}

func (u *LowonganUsecase) GetLowonganDetail(id int) (lowongandomain.Lowongan, error) {
	if id <= 0 {
		return lowongandomain.Lowongan{}, errors.New("ID lowongan tidak valid")
	}

	return u.repo.GetLowonganByID(id)
}

func (u *LowonganUsecase) CreateLowongan(request lowongandomain.CreateLowonganRequest) (lowongandomain.Lowongan, error) {
	if request.Judul == "" || request.Unit == "" {
		return lowongandomain.Lowongan{}, errors.New("judul dan unit wajib diisi")
	}

	tanggalBuka, err := parseOptionalDate(request.TanggalBuka, "tanggal buka")
	if err != nil {
		return lowongandomain.Lowongan{}, err
	}

	tanggalTutup, err := parseOptionalDate(request.TanggalTutup, "tanggal tutup")
	if err != nil {
		return lowongandomain.Lowongan{}, err
	}

	if err := validateLowonganDates(tanggalBuka, tanggalTutup); err != nil {
		return lowongandomain.Lowongan{}, err
	}

	lowongan := lowongandomain.Lowongan{
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

func (u *LowonganUsecase) UpdateLowonganStatus(id int, request lowongandomain.UpdateLowonganStatusRequest) (lowongandomain.Lowongan, error) {
	if request.Status != "aktif" && request.Status != "nonaktif" {
		return lowongandomain.Lowongan{}, errors.New("status harus aktif atau nonaktif")
	}

	return u.repo.UpdateLowonganStatus(id, request.Status)
}

func (u *LowonganUsecase) UpdateLowongan(id int, request lowongandomain.UpdateLowonganRequest) (lowongandomain.Lowongan, error) {
	if request.Judul == "" || request.Unit == "" {
		return lowongandomain.Lowongan{}, errors.New("judul dan unit wajib diisi")
	}

	if request.Status != "aktif" && request.Status != "nonaktif" {
		return lowongandomain.Lowongan{}, errors.New("status harus aktif atau nonaktif")
	}

	tanggalBuka, err := parseOptionalDate(request.TanggalBuka, "tanggal buka")
	if err != nil {
		return lowongandomain.Lowongan{}, err
	}

	tanggalTutup, err := parseOptionalDate(request.TanggalTutup, "tanggal tutup")
	if err != nil {
		return lowongandomain.Lowongan{}, err
	}

	if err := validateLowonganDates(tanggalBuka, tanggalTutup); err != nil {
		return lowongandomain.Lowongan{}, err
	}

	lowongan := lowongandomain.Lowongan{
		Judul:        strings.TrimSpace(request.Judul),
		Unit:         strings.TrimSpace(request.Unit),
		TanggalBuka:  tanggalBuka,
		TanggalTutup: tanggalTutup,
		Deskripsi:    strings.TrimSpace(request.Deskripsi),
		Status:       request.Status,
	}

	return u.repo.UpdateLowongan(id, lowongan)
}

func (u *LowonganUsecase) BulkUpdateLowonganStatus(request lowongandomain.BulkUpdateStatusRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu lowongan")
	}

	if request.Status != "aktif" && request.Status != "nonaktif" {
		return errors.New("status harus aktif atau nonaktif")
	}

	return u.repo.BulkUpdateLowonganStatus(request.IDs, request.Status)
}

func (u *LowonganUsecase) BulkDeleteLowongan(request lowongandomain.BulkDeleteRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("minimal pilih satu lowongan")
	}

	return u.repo.BulkDeleteLowongan(request.IDs)
}
