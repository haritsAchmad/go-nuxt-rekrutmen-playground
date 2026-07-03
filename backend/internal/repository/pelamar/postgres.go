package pelamar

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"
	"time"

	domain "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain/pelamar"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PelamarRepository struct {
	db *pgxpool.Pool
}

func NewPelamarRepository(db *pgxpool.Pool) *PelamarRepository {
	return &PelamarRepository{db: db}
}

func nullableTimeToPointer(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

func (r *PelamarRepository) GetAllPelamar(filter domain.PelamarFilterRequest) (domain.PelamarListResponse, error) {
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
	keyword := strings.TrimSpace(filter.Keyword)
	status := strings.TrimSpace(filter.Status)
	orderQuery := " ORDER BY p.id DESC"

	if filter.Sort == "oldest" {
		orderQuery = " ORDER BY p.id ASC"
	}

	whereQuery := `
		WHERE ($1 = '' OR LOWER(p.nama) LIKE '%' || LOWER($1) || '%' OR LOWER(p.email) LIKE '%' || LOWER($1) || '%' OR LOWER(COALESCE(p.no_hp, '')) LIKE '%' || LOWER($1) || '%' OR LOWER(COALESCE(l.judul, '')) LIKE '%' || LOWER($1) || '%')
		AND ($2 = '' OR p.status = $2)
		AND ($3 = 0 OR p.lowongan_id = $3)
	`

	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.lowongan_id, COALESCE(l.judul, ''), p.nama, p.email, COALESCE(p.no_hp, ''), p.status, p.created_at
		FROM pelamar p
		LEFT JOIN lowongan l ON l.id = p.lowongan_id
	`+whereQuery+orderQuery+`
		LIMIT $4 OFFSET $5
	`, keyword, status, filter.LowonganID, limit, offset)
	if err != nil {
		return domain.PelamarListResponse{}, err
	}
	defer rows.Close()

	data := []domain.Pelamar{}

	for rows.Next() {
		var pelamar domain.Pelamar
		var createdAt sql.NullTime

		err := rows.Scan(
			&pelamar.ID,
			&pelamar.LowonganID,
			&pelamar.LowonganJudul,
			&pelamar.Nama,
			&pelamar.Email,
			&pelamar.NoHP,
			&pelamar.Status,
			&createdAt,
		)
		if err != nil {
			return domain.PelamarListResponse{}, err
		}

		pelamar.CreatedAt = nullableTimeToPointer(createdAt)
		data = append(data, pelamar)
	}

	if rows.Err() != nil {
		return domain.PelamarListResponse{}, rows.Err()
	}

	var total int

	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM pelamar p
		LEFT JOIN lowongan l ON l.id = p.lowongan_id
	`+whereQuery, keyword, status, filter.LowonganID).Scan(&total)
	if err != nil {
		return domain.PelamarListResponse{}, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return domain.PelamarListResponse{
		Data: data,
		Meta: domain.PelamarPaginationMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (r *PelamarRepository) GetPelamarByID(id int) (domain.Pelamar, error) {
	ctx := context.Background()
	var pelamar domain.Pelamar
	var createdAt sql.NullTime

	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.lowongan_id, COALESCE(l.judul, ''), p.nama, p.email, COALESCE(p.no_hp, ''), p.status, p.created_at
		FROM pelamar p
		LEFT JOIN lowongan l ON l.id = p.lowongan_id
		WHERE p.id = $1
	`, id).Scan(
		&pelamar.ID,
		&pelamar.LowonganID,
		&pelamar.LowonganJudul,
		&pelamar.Nama,
		&pelamar.Email,
		&pelamar.NoHP,
		&pelamar.Status,
		&createdAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Pelamar{}, errors.New("pelamar tidak ditemukan")
	}
	if err != nil {
		return domain.Pelamar{}, err
	}

	pelamar.CreatedAt = nullableTimeToPointer(createdAt)
	return pelamar, nil
}

func (r *PelamarRepository) CreatePelamar(pelamar domain.Pelamar) (domain.Pelamar, error) {
	ctx := context.Background()
	var createdAt sql.NullTime

	err := r.db.QueryRow(ctx, `
		INSERT INTO pelamar (lowongan_id, nama, email, no_hp, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, lowongan_id, nama, email, COALESCE(no_hp, ''), status, created_at
	`, pelamar.LowonganID, pelamar.Nama, pelamar.Email, pelamar.NoHP, pelamar.Status).Scan(
		&pelamar.ID,
		&pelamar.LowonganID,
		&pelamar.Nama,
		&pelamar.Email,
		&pelamar.NoHP,
		&pelamar.Status,
		&createdAt,
	)
	if err != nil {
		return domain.Pelamar{}, err
	}

	pelamar.CreatedAt = nullableTimeToPointer(createdAt)
	return r.GetPelamarByID(pelamar.ID)
}

func (r *PelamarRepository) UpdatePelamar(id int, updatedPelamar domain.Pelamar) (domain.Pelamar, error) {
	ctx := context.Background()

	var pelamar domain.Pelamar
	var createdAt sql.NullTime

	err := r.db.QueryRow(ctx, `
		UPDATE pelamar
		SET lowongan_id = $1, nama = $2, email = $3, no_hp = $4, status = $5
		WHERE id = $6
		RETURNING id, lowongan_id, nama, email, COALESCE(no_hp, ''), status, created_at
	`, updatedPelamar.LowonganID, updatedPelamar.Nama, updatedPelamar.Email, updatedPelamar.NoHP, updatedPelamar.Status, id).Scan(
		&pelamar.ID,
		&pelamar.LowonganID,
		&pelamar.Nama,
		&pelamar.Email,
		&pelamar.NoHP,
		&pelamar.Status,
		&createdAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Pelamar{}, errors.New("pelamar tidak ditemukan")
	}
	if err != nil {
		return domain.Pelamar{}, err
	}

	pelamar.CreatedAt = nullableTimeToPointer(createdAt)
	return r.GetPelamarByID(pelamar.ID)
}

func (r *PelamarRepository) UpdatePelamarStatus(id int, status string) (domain.Pelamar, error) {
	ctx := context.Background()

	var pelamar domain.Pelamar
	var createdAt sql.NullTime

	err := r.db.QueryRow(ctx, `
		UPDATE pelamar
		SET status = $1
		WHERE id = $2
		RETURNING id, lowongan_id, nama, email, COALESCE(no_hp, ''), status, created_at
	`, status, id).Scan(
		&pelamar.ID,
		&pelamar.LowonganID,
		&pelamar.Nama,
		&pelamar.Email,
		&pelamar.NoHP,
		&pelamar.Status,
		&createdAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Pelamar{}, errors.New("pelamar tidak ditemukan")
	}
	if err != nil {
		return domain.Pelamar{}, err
	}

	pelamar.CreatedAt = nullableTimeToPointer(createdAt)
	return r.GetPelamarByID(pelamar.ID)
}

func (r *PelamarRepository) DeletePelamar(id int) error {
	ctx := context.Background()

	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM pelamar
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("pelamar tidak ditemukan")
	}

	return nil
}

func (r *PelamarRepository) BulkUpdatePelamarStatus(ids []int, status string) error {
	ctx := context.Background()

	commandTag, err := r.db.Exec(ctx, `
		UPDATE pelamar
		SET status = $1
		WHERE id = ANY($2)
	`, status, ids)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("pelamar tidak ditemukan")
	}

	return nil
}

func (r *PelamarRepository) BulkDeletePelamar(ids []int) error {
	ctx := context.Background()

	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM pelamar
		WHERE id = ANY($1)
	`, ids)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("pelamar tidak ditemukan")
	}

	return nil
}
