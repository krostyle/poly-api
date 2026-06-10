package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
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
	return scanBanco(r.pool.QueryRow(ctx, q, estudioID, nombre))
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
	return scanBanco(r.pool.QueryRow(ctx, q, id, estudioID))
}

func (r *BancoRepo) Update(ctx context.Context, estudioID, id, nombre string) (*domain.Banco, error) {
	const q = `UPDATE bancos SET nombre=$3 WHERE id=$1 AND estudio_id=$2 RETURNING id, estudio_id, nombre, created_at`
	return scanBanco(r.pool.QueryRow(ctx, q, id, estudioID, nombre))
}

func (r *BancoRepo) Delete(ctx context.Context, estudioID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM bancos WHERE id=$1 AND estudio_id=$2`, id, estudioID)
	return err
}

func (r *BancoRepo) HasCasos(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM casos WHERE banco_id=$1)`, id).Scan(&exists)
	return exists, err
}

func (r *BancoRepo) ListUsuarios(ctx context.Context, bancoID string) ([]*domain.UsuarioBanco, error) {
	const q = `SELECT u.id, u.nombre, u.email, u.rol
		FROM usuarios u
		JOIN usuarios_bancos ub ON ub.usuario_id = u.id
		WHERE ub.banco_id = $1
		ORDER BY u.nombre`
	rows, err := r.pool.Query(ctx, q, bancoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.UsuarioBanco
	for rows.Next() {
		var u domain.UsuarioBanco
		if err := rows.Scan(&u.ID, &u.Nombre, &u.Email, &u.Rol); err != nil {
			return nil, err
		}
		result = append(result, &u)
	}
	return result, rows.Err()
}

func (r *BancoRepo) AsignarUsuario(ctx context.Context, bancoID, usuarioID string) error {
	const q = `INSERT INTO usuarios_bancos (usuario_id, banco_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, q, usuarioID, bancoID)
	return err
}

func (r *BancoRepo) DesasignarUsuario(ctx context.Context, bancoID, usuarioID string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM usuarios_bancos WHERE banco_id=$1 AND usuario_id=$2`, bancoID, usuarioID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

type bancoScanner interface {
	Scan(dest ...any) error
}

func scanBanco(row bancoScanner) (*domain.Banco, error) {
	var b domain.Banco
	var createdAt time.Time
	if err := row.Scan(&b.ID, &b.EstudioID, &b.Nombre, &createdAt); err != nil {
		return nil, err
	}
	b.CreatedAt = createdAt
	return &b, nil
}
