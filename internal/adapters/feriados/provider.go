package feriados

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBFeriadoProvider struct {
	pool *pgxpool.Pool
}

func NewDBFeriadoProvider(pool *pgxpool.Pool) *DBFeriadoProvider {
	return &DBFeriadoProvider{pool: pool}
}

func (p *DBFeriadoProvider) GetHolidays(ctx context.Context, from, to time.Time) ([]time.Time, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT fecha FROM feriados WHERE fecha BETWEEN $1 AND $2 ORDER BY fecha ASC`,
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []time.Time
	for rows.Next() {
		var fecha time.Time
		if err := rows.Scan(&fecha); err != nil {
			return nil, err
		}
		result = append(result, fecha)
	}
	return result, rows.Err()
}
