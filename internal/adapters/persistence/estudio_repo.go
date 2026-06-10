package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
)

type EstudioRepo struct {
	pool *pgxpool.Pool
}

func NewEstudioRepo(pool *pgxpool.Pool) *EstudioRepo {
	return &EstudioRepo{pool: pool}
}

func (r *EstudioRepo) UpsertByClerkOrgID(ctx context.Context, clerkOrgID, nombre string) (*domain.Estudio, error) {
	const q = `
		INSERT INTO estudios (clerk_org_id, nombre)
		VALUES ($1, $2)
		ON CONFLICT (clerk_org_id)
		DO UPDATE SET nombre = EXCLUDED.nombre
		RETURNING id, clerk_org_id, nombre, rut, created_at`

	row := r.pool.QueryRow(ctx, q, clerkOrgID, nombre)
	return scanEstudio(row)
}

func (r *EstudioRepo) GetByClerkOrgID(ctx context.Context, clerkOrgID string) (*domain.Estudio, error) {
	const q = `SELECT id, clerk_org_id, nombre, rut, created_at FROM estudios WHERE clerk_org_id = $1`
	row := r.pool.QueryRow(ctx, q, clerkOrgID)
	return scanEstudio(row)
}

type estudioScanner interface {
	Scan(dest ...any) error
}

func scanEstudio(row estudioScanner) (*domain.Estudio, error) {
	var e domain.Estudio
	var createdAt time.Time
	err := row.Scan(&e.ID, &e.ClerkOrgID, &e.Nombre, &e.RUT, &createdAt)
	if err != nil {
		return nil, err
	}
	e.CreatedAt = createdAt
	return &e, nil
}
