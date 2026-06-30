package postgres

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/domain"
	"github.com/jackc/pgx/v5"
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

func nullableTimeToPointer(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
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
	keyword := strings.TrimSpace(filter.Keyword)
	status := strings.TrimSpace(filter.Status)

	whereQuery := `
		WHERE ($1 = '' OR LOWER(judul) LIKE '%' || LOWER($1) || '%' OR LOWER(unit) LIKE '%' || LOWER($1) || '%')
		AND ($2 = '' OR status = $2)
	`

	rows, err := r.db.Query(ctx, `
		SELECT id, judul, unit, tanggal_buka, tanggal_tutup, COALESCE(deskripsi, ''), status
		FROM lowongan
	`+whereQuery+`
		ORDER BY id
		LIMIT $3 OFFSET $4
	`, keyword, status, limit, offset)
	if err != nil {
		return domain.LowonganListResponse{}, err
	}
	defer rows.Close()

	data := []domain.Lowongan{}

	for rows.Next() {
		var lowongan domain.Lowongan
		var tanggalBuka sql.NullTime
		var tanggalTutup sql.NullTime

		err := rows.Scan(
			&lowongan.ID,
			&lowongan.Judul,
			&lowongan.Unit,
			&tanggalBuka,
			&tanggalTutup,
			&lowongan.Deskripsi,
			&lowongan.Status,
		)
		if err != nil {
			return domain.LowonganListResponse{}, err
		}

		lowongan.TanggalBuka = nullableTimeToPointer(tanggalBuka)
		lowongan.TanggalTutup = nullableTimeToPointer(tanggalTutup)
		data = append(data, lowongan)
	}

	if rows.Err() != nil {
		return domain.LowonganListResponse{}, rows.Err()
	}

	var total int

	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM lowongan
	`+whereQuery, keyword, status).Scan(&total)
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

func (r *LowonganRepository) GetLowonganByID(id int) (domain.Lowongan, error) {
	ctx := context.Background()

	var lowongan domain.Lowongan
	var tanggalBuka sql.NullTime
	var tanggalTutup sql.NullTime

	err := r.db.QueryRow(ctx, `
		SELECT id, judul, unit, tanggal_buka, tanggal_tutup, COALESCE(deskripsi, ''), status
		FROM lowongan
		WHERE id = $1
	`, id).Scan(
		&lowongan.ID,
		&lowongan.Judul,
		&lowongan.Unit,
		&tanggalBuka,
		&tanggalTutup,
		&lowongan.Deskripsi,
		&lowongan.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Lowongan{}, errors.New("lowongan tidak ditemukan")
	}
	if err != nil {
		return domain.Lowongan{}, err
	}

	lowongan.TanggalBuka = nullableTimeToPointer(tanggalBuka)
	lowongan.TanggalTutup = nullableTimeToPointer(tanggalTutup)
	return lowongan, nil
}

func (r *LowonganRepository) CreateLowongan(lowongan domain.Lowongan) (domain.Lowongan, error) {
	ctx := context.Background()
	var tanggalBuka sql.NullTime
	var tanggalTutup sql.NullTime

	err := r.db.QueryRow(ctx, `
		INSERT INTO lowongan (judul, unit, tanggal_buka, tanggal_tutup, deskripsi, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, judul, unit, tanggal_buka, tanggal_tutup, COALESCE(deskripsi, ''), status
	`, lowongan.Judul, lowongan.Unit, lowongan.TanggalBuka, lowongan.TanggalTutup, lowongan.Deskripsi, lowongan.Status).Scan(
		&lowongan.ID,
		&lowongan.Judul,
		&lowongan.Unit,
		&tanggalBuka,
		&tanggalTutup,
		&lowongan.Deskripsi,
		&lowongan.Status,
	)
	if err != nil {
		return domain.Lowongan{}, err
	}

	lowongan.TanggalBuka = nullableTimeToPointer(tanggalBuka)
	lowongan.TanggalTutup = nullableTimeToPointer(tanggalTutup)
	return lowongan, nil
}

func (r *LowonganRepository) DeleteLowongan(id int) error {
	ctx := context.Background()

	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM lowongan
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("lowongan tidak ditemukan")
	}

	return nil
}

func (r *LowonganRepository) UpdateLowonganStatus(id int, status string) (domain.Lowongan, error) {
	ctx := context.Background()

	var lowongan domain.Lowongan
	var tanggalBuka sql.NullTime
	var tanggalTutup sql.NullTime

	err := r.db.QueryRow(ctx, `
		UPDATE lowongan
		SET status = $1
		WHERE id = $2
		RETURNING id, judul, unit, tanggal_buka, tanggal_tutup, COALESCE(deskripsi, ''), status
	`, status, id).Scan(
		&lowongan.ID,
		&lowongan.Judul,
		&lowongan.Unit,
		&tanggalBuka,
		&tanggalTutup,
		&lowongan.Deskripsi,
		&lowongan.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Lowongan{}, errors.New("lowongan tidak ditemukan")
	}
	if err != nil {
		return domain.Lowongan{}, err
	}

	lowongan.TanggalBuka = nullableTimeToPointer(tanggalBuka)
	lowongan.TanggalTutup = nullableTimeToPointer(tanggalTutup)
	return lowongan, nil
}

func (r *LowonganRepository) UpdateLowongan(id int, updatedLowongan domain.Lowongan) (domain.Lowongan, error) {
	ctx := context.Background()

	var lowongan domain.Lowongan
	var tanggalBuka sql.NullTime
	var tanggalTutup sql.NullTime

	err := r.db.QueryRow(ctx, `
		UPDATE lowongan
		SET judul = $1, unit = $2, tanggal_buka = $3, tanggal_tutup = $4, deskripsi = $5, status = $6
		WHERE id = $7
		RETURNING id, judul, unit, tanggal_buka, tanggal_tutup, COALESCE(deskripsi, ''), status
	`, updatedLowongan.Judul, updatedLowongan.Unit, updatedLowongan.TanggalBuka, updatedLowongan.TanggalTutup, updatedLowongan.Deskripsi, updatedLowongan.Status, id).Scan(
		&lowongan.ID,
		&lowongan.Judul,
		&lowongan.Unit,
		&tanggalBuka,
		&tanggalTutup,
		&lowongan.Deskripsi,
		&lowongan.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Lowongan{}, errors.New("lowongan tidak ditemukan")
	}
	if err != nil {
		return domain.Lowongan{}, err
	}

	lowongan.TanggalBuka = nullableTimeToPointer(tanggalBuka)
	lowongan.TanggalTutup = nullableTimeToPointer(tanggalTutup)
	return lowongan, nil
}

func (r *LowonganRepository) BulkUpdateLowonganStatus(ids []int, status string) error {
	ctx := context.Background()

	commandTag, err := r.db.Exec(ctx, `
		UPDATE lowongan
		SET status = $1
		WHERE id = ANY($2)
	`, status, ids)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("lowongan tidak ditemukan")
	}

	return nil
}

func (r *LowonganRepository) BulkDeleteLowongan(ids []int) error {
	ctx := context.Background()

	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM lowongan
		WHERE id = ANY($1)
	`, ids)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("lowongan tidak ditemukan")
	}

	return nil
}
