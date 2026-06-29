package postgres

import (
	"context"
	"math"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LowonganRepository struct {
	db *pgxpool.Pool
}

func NewLowonganRepository(db *pgxpool.Pool) *LowonganRepository {
	return &LowonganRepository{
		db: db,
	}
}

func (r *LowonganRepository) GetAllLowongan(filter domain.LowonganFilterRequest) (domain.LowonganListResponse, error) {
	ctx := context.Background()

	page := filter.Page
	limit := filter.Limit

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	rows, err := r.db.Query(ctx, `
		SELECT id, judul, unit, status
		FROM rekrutmen_playground.lowongan
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return domain.LowonganListResponse{}, err
	}
	defer rows.Close()

	data := []domain.Lowongan{}

	for rows.Next() {
		var lowongan domain.Lowongan

		err := rows.Scan(
			&lowongan.ID,
			&lowongan.Judul,
			&lowongan.Unit,
			&lowongan.Status,
		)
		if err != nil {
			return domain.LowonganListResponse{}, err
		}

		data = append(data, lowongan)
	}

	if rows.Err() != nil {
		return domain.LowonganListResponse{}, rows.Err()
	}

	var total int

	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM lowongan
	`).Scan(&total)
	if err != nil {
		return domain.LowonganListResponse{}, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return domain.LowonganListResponse{
		Data: data,
		Meta: domain.LowonganPaginationMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}
