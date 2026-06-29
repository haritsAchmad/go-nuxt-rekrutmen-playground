package domain

type Lowongan struct {
	ID     int    `json:"id"`
	Judul  string `json:"judul"`
	Unit   string `json:"unit"`
	Status string `json:"status"`
}

type CreateLowonganRequest struct {
	Judul string `json:"judul"`
	Unit  string `json:"unit"`
}

type UpdateLowonganStatusRequest struct {
	Status string `json:"status"`
}

type UpdateLowonganRequest struct {
	Judul  string `json:"judul"`
	Unit   string `json:"unit"`
	Status string `json:"status"`
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
