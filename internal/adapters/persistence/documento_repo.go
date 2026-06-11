package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
)

type DocumentoRepo struct {
	pool *pgxpool.Pool
}

func NewDocumentoRepo(pool *pgxpool.Pool) *DocumentoRepo {
	return &DocumentoRepo{pool: pool}
}

func (r *DocumentoRepo) Create(ctx context.Context, in domain.NewDocumentoInput) (*domain.Documento, error) {
	var d domain.Documento
	err := r.pool.QueryRow(ctx,
		`INSERT INTO documentos (caso_id, tipo, blob_url, nombre, subido_por)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, caso_id, tipo, blob_url, nombre, subido_por, created_at`,
		in.CasoID, in.Tipo, in.BlobURL, in.Nombre, in.SubidoPor,
	).Scan(&d.ID, &d.CasoID, &d.Tipo, &d.BlobURL, &d.Nombre, &d.SubidoPor, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DocumentoRepo) ListByCaso(ctx context.Context, casoID string) ([]*domain.Documento, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT d.id, d.caso_id, d.tipo, d.blob_url, d.nombre, d.subido_por, d.created_at
		 FROM documentos d
		 WHERE d.caso_id = $1
		 ORDER BY d.created_at DESC`,
		casoID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Documento
	for rows.Next() {
		var d domain.Documento
		var subidoPor *string
		var createdAt time.Time
		if err := rows.Scan(&d.ID, &d.CasoID, &d.Tipo, &d.BlobURL, &d.Nombre, &subidoPor, &createdAt); err != nil {
			return nil, err
		}
		d.SubidoPor = subidoPor
		d.CreatedAt = createdAt
		result = append(result, &d)
	}
	return result, rows.Err()
}
