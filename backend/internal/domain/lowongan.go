package domain

import "time"

type Lowongan struct {
	ID           int        `json:"id"`
	Judul        string     `json:"judul"`
	Unit         string     `json:"unit"`
	TanggalBuka  *time.Time `json:"tanggalBuka,omitempty"`
	TanggalTutup *time.Time `json:"tanggalTutup,omitempty"`
	Deskripsi    string     `json:"deskripsi"`
	Status       string     `json:"status"`
}

type CreateLowonganRequest struct {
	Judul        string `json:"judul"`
	Unit         string `json:"unit"`
	TanggalBuka  string `json:"tanggalBuka"`
	TanggalTutup string `json:"tanggalTutup"`
	Deskripsi    string `json:"deskripsi"`
}

type UpdateLowonganStatusRequest struct {
	Status string `json:"status"`
}

type UpdateLowonganRequest struct {
	Judul        string `json:"judul"`
	Unit         string `json:"unit"`
	TanggalBuka  string `json:"tanggalBuka"`
	TanggalTutup string `json:"tanggalTutup"`
	Deskripsi    string `json:"deskripsi"`
	Status       string `json:"status"`
}

type LowonganFilterRequest struct {
	Keyword string
	Status  string
	Page    int
	Limit   int
}

type BulkUpdateStatusRequest struct {
	IDs    []int  `json:"ids"`
	Status string `json:"status"`
}

type BulkDeleteRequest struct {
	IDs []int `json:"ids"`
}

type LowonganPaginationMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

type LowonganListResponse struct {
	Data []Lowongan             `json:"data"`
	Meta LowonganPaginationMeta `json:"meta"`
}
