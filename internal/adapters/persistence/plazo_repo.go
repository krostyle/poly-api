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

func (r *PlazoRepo) ListGlobal(ctx context.Context, estudioID string, bancoIDs []string, tipoFilter string) ([]domain.StoredPlazoGlobal, error) {
	q := `
		SELECT p.id, p.caso_id, p.tipo, p.fecha_inicio, p.dias_habiles, p.fecha_limite,
		       c.numero_ot, cl.nombre, cl.rut, b.nombre, c.estado
		FROM plazos p
		JOIN casos    c  ON c.id  = p.caso_id
		JOIN clientes cl ON cl.id = c.cliente_id
		JOIN bancos   b  ON b.id  = c.banco_id
		WHERE c.estudio_id = $1
		  AND c.banco_id = ANY($2)
		  AND p.cumplido = false
		  AND c.estado NOT IN ('CIERRE','TERMINADO')
		  AND ($3 = '' OR p.tipo = $3)
		ORDER BY p.fecha_limite ASC
		LIMIT 200`

	rows, err := r.pool.Query(ctx, q, estudioID, bancoIDs, tipoFilter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.StoredPlazoGlobal
	for rows.Next() {
		var g domain.StoredPlazoGlobal
		var tipo string
		if err := rows.Scan(
			&g.ID, &g.CasoID, &tipo, &g.FechaInicio, &g.DiasHabiles, &g.FechaLimite,
			&g.NumeroOT, &g.ClienteNombre, &g.ClienteRUT, &g.BancoNombre, &g.Estado,
		); err != nil {
			return nil, err
		}
		g.Tipo = plazo.TipoPlazo(tipo)
		result = append(result, g)
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

func (r *PlazoRepo) UpdateDiasHabiles(ctx context.Context, plazoID string, diasHabiles int, fechaLimite time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE plazos SET dias_habiles = $2, fecha_limite = $3 WHERE id = $1 AND cumplido = false`,
		plazoID, diasHabiles, fechaLimite,
	)
	return err
}
