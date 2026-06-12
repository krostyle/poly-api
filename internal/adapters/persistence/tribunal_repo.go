package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
)

type TribunalRepo struct {
	pool *pgxpool.Pool
}

func NewTribunalRepo(pool *pgxpool.Pool) *TribunalRepo {
	return &TribunalRepo{pool: pool}
}

func (r *TribunalRepo) List(ctx context.Context) ([]*domain.Tribunal, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, nombre, region, activo, created_at
		 FROM tribunales WHERE activo = true ORDER BY region ASC, nombre ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Tribunal
	for rows.Next() {
		var t domain.Tribunal
		var createdAt time.Time
		if err := rows.Scan(&t.ID, &t.Nombre, &t.Region, &t.Activo, &createdAt); err != nil {
			return nil, err
		}
		t.CreatedAt = createdAt
		result = append(result, &t)
	}
	return result, rows.Err()
}

func (r *TribunalRepo) Create(ctx context.Context, nombre, region string) (*domain.Tribunal, error) {
	id := uuid.New().String()
	now := time.Now()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tribunales (id, nombre, region, activo, created_at)
		 VALUES ($1, $2, $3, true, $4)`,
		id, nombre, region, now,
	)
	if err != nil {
		return nil, err
	}
	return &domain.Tribunal{ID: id, Nombre: nombre, Region: region, Activo: true, CreatedAt: now}, nil
}
