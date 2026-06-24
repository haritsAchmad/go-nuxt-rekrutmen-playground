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