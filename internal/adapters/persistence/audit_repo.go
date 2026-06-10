package persistence

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
)

type AuditRepo struct {
	pool *pgxpool.Pool
}

func NewAuditRepo(pool *pgxpool.Pool) *AuditRepo {
	return &AuditRepo{pool: pool}
}

func (r *AuditRepo) Log(ctx context.Context, entry domain.AuditEntry) error {
	det, _ := json.Marshal(entry.Detalle)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO auditoria (estudio_id, usuario_id, caso_id, accion, detalle) VALUES ($1,$2,$3,$4,$5)`,
		entry.EstudioID, entry.UsuarioID, entry.CasoID, entry.Accion, det,
	)
	return err
}
