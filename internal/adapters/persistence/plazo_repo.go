package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/plazo"
)

type PlazoRepo struct {
	pool *pgxpool.Pool
}

func NewPlazoRepo(pool *pgxpool.Pool) *PlazoRepo {
	return &PlazoRepo{pool: pool}
}

func (r *PlazoRepo) CreateBatch(ctx context.Context, plazos []domain.NewPlazoInput) error {
	for _, p := range plazos {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO plazos (caso_id, tipo, fecha_inicio, dias_habiles, fecha_limite)
			 VALUES ($1, $2, $3, $4, $5)`,
			p.CasoID, string(p.Tipo), p.FechaInicio, p.DiasHabiles, p.FechaLimite,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PlazoRepo) ListByCase(ctx context.Context, casoID string) ([]domain.StoredPlazo, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, caso_id, tipo, fecha_inicio, dias_habiles, fecha_limite, cumplido, fecha_cumplido
		 FROM plazos WHERE caso_id = $1 ORDER BY fecha_limite ASC`,
		casoID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.StoredPlazo
	for rows.Next() {
		var p domain.StoredPlazo
		var tipo string
		var fechaCumplido *time.Time
		if err := rows.Scan(
			&p.ID, &p.CasoID, &tipo, &p.FechaInicio, &p.DiasHabiles,
			&p.FechaLimite, &p.Completed, &fechaCumplido,
		); err != nil {
			return nil, err
		}
		p.Tipo = plazo.TipoPlazo(tipo)
		p.FechaCumplido = fechaCumplido
		result = append(result, p)
	}
	return result, rows.Err()
}

func (r *PlazoRepo) MarkCompleted(ctx context.Context, plazoID string, date time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE plazos SET cumplido = true, fecha_cumplido = $2 WHERE id = $1`,
		plazoID, date,
	)
	return err
}
