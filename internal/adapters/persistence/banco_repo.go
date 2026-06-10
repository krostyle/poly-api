package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
)

type BancoRepo struct {
	pool *pgxpool.Pool
}

func NewBancoRepo(pool *pgxpool.Pool) *BancoRepo {
	return &BancoRepo{pool: pool}
}

func (r *BancoRepo) Create(ctx context.Context, estudioID, nombre string) (*domain.Banco, error) {
	const q = `INSERT INTO bancos (estudio_id, nombre) VALUES ($1, $2) RETURNING id, estudio_id, nombre, created_at`
	row := r.pool.QueryRow(ctx, q, estudioID, nombre)
	return scanBanco(row)
}

func (r *BancoRepo) List(ctx context.Context, estudioID string) ([]*domain.Banco, error) {
	const q = `SELECT id, estudio_id, nombre, created_at FROM bancos WHERE estudio_id = $1 ORDER BY nombre`
	rows, err := r.pool.Query(ctx, q, estudioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bancos []*domain.Banco
	for rows.Next() {
		var b domain.Banco
		var createdAt time.Time
		if err := rows.Scan(&b.ID, &b.EstudioID, &b.Nombre, &createdAt); err != nil {
			return nil, err
		}
		b.CreatedAt = createdAt
		bancos = append(bancos, &b)
	}
	return bancos, rows.Err()
}

func (r *BancoRepo) GetByID(ctx context.Context, estudioID, id string) (*domain.Banco, error) {
	const q = `SELECT id, estudio_id, nombre, created_at FROM bancos WHERE id = $1 AND estudio_id = $2`
	row := r.pool.QueryRow(ctx, q, id, estudioID)
	return scanBanco(row)
}

func (r *BancoRepo) AssignToUsuario(ctx context.Context, usuarioID, bancoID string) error {
	const q = `INSERT INTO usuarios_bancos (usuario_id, banco_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, q, usuarioID, bancoID)
	return err
}

type bancoScanner interface {
	Scan(dest ...any) error
}

func scanBanco(row bancoScanner) (*domain.Banco, error) {
	var b domain.Banco
	var createdAt time.Time
	err := row.Scan(&b.ID, &b.EstudioID, &b.Nombre, &createdAt)
	if err != nil {
		return nil, err
	}
	b.CreatedAt = createdAt
	return &b, nil
}
