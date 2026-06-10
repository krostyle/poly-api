package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
)

type OperacionRepo struct {
	pool *pgxpool.Pool
}

func NewOperacionRepo(pool *pgxpool.Pool) *OperacionRepo {
	return &OperacionRepo{pool: pool}
}

func (r *OperacionRepo) Create(ctx context.Context, in domain.NewOperacionInput) (*domain.Operacion, error) {
	const q = `INSERT INTO operaciones (caso_id, medio_pago, relacion, monto_clp, monto_uf, fecha_op)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, caso_id, medio_pago, relacion, monto_clp, monto_uf, fecha_op, created_at`
	row := r.pool.QueryRow(ctx, q,
		in.CasoID, in.MedioPago, in.Relacion, in.MontoCLP, in.MontoUF, in.FechaOp,
	)
	var op domain.Operacion
	err := row.Scan(&op.ID, &op.CasoID, &op.MedioPago, &op.Relacion,
		&op.MontoCLP, &op.MontoUF, &op.FechaOp, &op.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *OperacionRepo) ListByCaso(ctx context.Context, casoID string) ([]*domain.Operacion, error) {
	const q = `SELECT id, caso_id, medio_pago, relacion, monto_clp, monto_uf, fecha_op, created_at
		FROM operaciones WHERE caso_id = $1 ORDER BY fecha_op DESC`
	rows, err := r.pool.Query(ctx, q, casoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ops []*domain.Operacion
	for rows.Next() {
		var op domain.Operacion
		if err := rows.Scan(&op.ID, &op.CasoID, &op.MedioPago, &op.Relacion,
			&op.MontoCLP, &op.MontoUF, &op.FechaOp, &op.CreatedAt); err != nil {
			return nil, err
		}
		ops = append(ops, &op)
	}
	return ops, rows.Err()
}
