package pelamar

import "time"

type Pelamar struct {
	ID             int        `json:"id"`
	LowonganID     int        `json:"lowonganId"`
	LowonganJudul  string     `json:"lowonganJudul"`
	Nama           string     `json:"nama"`
	Email          string     `json:"email"`
	NoHP           string     `json:"noHp"`
	Status         string     `json:"status"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
}

type CreatePelamarRequest struct {
	LowonganID int    `json:"lowonganId"`
	Nama       string `json:"nama"`
	Email      string `json:"email"`
	NoHP       string `json:"noHp"`
}

type UpdatePelamarRequest struct {
	LowonganID int    `json:"lowonganId"`
	Nama       string `json:"nama"`
	Email      string `json:"email"`
	NoHP       string `json:"noHp"`
	Status     string `json:"status"`
}

type UpdatePelamarStatusRequest struct {
	Status string `json:"status"`
}

type PelamarFilterRequest struct {
	Keyword    string
	Status     string
	Sort       string
	LowonganID int
	Page       int
	Limit      int
}

type BulkUpdateStatusRequest struct {
	IDs    []int  `json:"ids"`
	Status string `json:"status"`
}

type BulkDeleteRequest struct {
	IDs []int `json:"ids"`
}

type PelamarPaginationMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

type PelamarListResponse struct {
	Data []Pelamar             `json:"data"`
	Meta PelamarPaginationMeta `json:"meta"`
}
