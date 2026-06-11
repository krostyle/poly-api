package persistence

import (
	"context"
	"encoding/json"
	"time"

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

func (r *AuditRepo) ListByCaso(ctx context.Context, estudioID, casoID string) ([]*domain.FullAuditEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT a.id::text,
		       a.accion,
		       a.detalle,
		       COALESCE(u.nombre, 'Sistema'),
		       a.created_at
		FROM auditoria a
		LEFT JOIN usuarios u ON u.id = a.usuario_id
		WHERE a.estudio_id = $1
		  AND a.caso_id = $2
		ORDER BY a.created_at DESC
	`, estudioID, casoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.FullAuditEntry
	for rows.Next() {
		var e domain.FullAuditEntry
		var createdAt time.Time
		var detalle []byte
		if err := rows.Scan(&e.ID, &e.Accion, &detalle, &e.UsuarioNombre, &createdAt); err != nil {
			return nil, err
		}
		if len(detalle) > 0 {
			_ = json.Unmarshal(detalle, &e.Detalle)
		}
		e.CreatedAt = createdAt
		entries = append(entries, &e)
	}
	return entries, rows.Err()
}
