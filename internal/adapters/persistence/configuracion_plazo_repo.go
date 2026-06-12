package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/plazo"
)

type ConfiguracionPlazoRepo struct {
	pool *pgxpool.Pool
}

func NewConfiguracionPlazoRepo(pool *pgxpool.Pool) *ConfiguracionPlazoRepo {
	return &ConfiguracionPlazoRepo{pool: pool}
}

func (r *ConfiguracionPlazoRepo) GetByEstudio(ctx context.Context, estudioID string) ([]domain.ConfiguracionPlazo, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, estudio_id, tipo_plazo, dias_habiles FROM configuracion_plazos WHERE estudio_id = $1`,
		estudioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ConfiguracionPlazo
	for rows.Next() {
		var c domain.ConfiguracionPlazo
		var tipo string
		if err := rows.Scan(&c.ID, &c.EstudioID, &tipo, &c.DiasHabiles); err != nil {
			return nil, err
		}
		c.TipoPlazo = plazo.TipoPlazo(tipo)
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *ConfiguracionPlazoRepo) Upsert(ctx context.Context, estudioID string, tipoPlazo plazo.TipoPlazo, diasHabiles int) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO configuracion_plazos (id, estudio_id, tipo_plazo, dias_habiles)
		VALUES (gen_random_uuid(), $1, $2, $3)
		ON CONFLICT (estudio_id, tipo_plazo)
		DO UPDATE SET dias_habiles = EXCLUDED.dias_habiles, updated_at = NOW()
	`, estudioID, string(tipoPlazo), diasHabiles)
	return err
}
