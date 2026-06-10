package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
)

type UsuarioRepo struct {
	pool *pgxpool.Pool
}

func NewUsuarioRepo(pool *pgxpool.Pool) *UsuarioRepo {
	return &UsuarioRepo{pool: pool}
}

func (r *UsuarioRepo) UpsertByClerkUserID(ctx context.Context, in domain.UpsertUsuarioInput) (*domain.Usuario, error) {
	const q = `
		INSERT INTO usuarios (clerk_user_id, estudio_id, nombre, email, rol)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (clerk_user_id)
		DO UPDATE SET nombre = EXCLUDED.nombre, email = EXCLUDED.email
		RETURNING id, clerk_user_id, estudio_id, nombre, email, rol, created_at`

	row := r.pool.QueryRow(ctx, q, in.ClerkUserID, in.EstudioID, in.Nombre, in.Email, in.Rol)
	return scanUsuario(row)
}

func (r *UsuarioRepo) GetByClerkUserID(ctx context.Context, clerkUserID string) (*domain.Usuario, error) {
	const q = `SELECT id, clerk_user_id, estudio_id, nombre, email, rol, created_at FROM usuarios WHERE clerk_user_id = $1`
	row := r.pool.QueryRow(ctx, q, clerkUserID)
	return scanUsuario(row)
}

func (r *UsuarioRepo) GetBancoIDs(ctx context.Context, usuarioID string) ([]string, error) {
	const q = `SELECT banco_id FROM usuarios_bancos WHERE usuario_id = $1`
	rows, err := r.pool.Query(ctx, q, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

type usuarioScanner interface {
	Scan(dest ...any) error
}

func scanUsuario(row usuarioScanner) (*domain.Usuario, error) {
	var u domain.Usuario
	var createdAt time.Time
	err := row.Scan(&u.ID, &u.ClerkUserID, &u.EstudioID, &u.Nombre, &u.Email, &u.Rol, &createdAt)
	if err != nil {
		return nil, err
	}
	u.CreatedAt = createdAt
	return &u, nil
}
