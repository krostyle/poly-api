package dashboard

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain/plazo"
)

type CasoPorVencer struct {
	CasoID        string
	BancoID       string
	BancoNombre   string
	NumeroOT      *string
	ClienteRUT    string
	ClienteNombre string
	Estado        string
	PlazoCritico  PlazoCritico
}

type PlazoCritico struct {
	ID            string
	Tipo          string
	FechaLimite   time.Time
	DiasRestantes int
	Semaforo      string
}

type CasoNuevo struct {
	CasoID        string
	BancoID       string
	BancoNombre   string
	ClienteRUT    string
	ClienteNombre string
	AbogadoID     *string
	CreatedAt     time.Time
}

type CasoEstancado struct {
	CasoID           string
	BancoID          string
	BancoNombre      string
	NumeroOT         *string
	ClienteRUT       string
	ClienteNombre    string
	Estado           string
	UltimoMovimiento time.Time
	DiasEstancado    int
}

type CargaAbogado struct {
	AbogadoID string
	Nombre    string
	Total     int
	PorVencer int
	Vencidos  int
}

type DashboardUseCase struct {
	pool *pgxpool.Pool
}

func NewDashboardUseCase(pool *pgxpool.Pool) *DashboardUseCase {
	return &DashboardUseCase{pool: pool}
}

func (uc *DashboardUseCase) PorVencer(ctx context.Context, estudioID string, bancoIDs []string, dias int) ([]CasoPorVencer, error) {
	if len(bancoIDs) == 0 {
		return nil, nil
	}
	horizon := time.Now().AddDate(0, 0, dias)

	rows, err := uc.pool.Query(ctx, `
		SELECT sub.caso_id, sub.banco_id, sub.banco_nombre, sub.numero_ot,
		       sub.cliente_rut, sub.cliente_nombre, sub.estado,
		       sub.plazo_id, sub.plazo_tipo, sub.fecha_limite
		FROM (
			SELECT DISTINCT ON (c.id)
				c.id        AS caso_id,
				c.banco_id,
				b.nombre    AS banco_nombre,
				c.numero_ot,
				cl.rut      AS cliente_rut,
				cl.nombre   AS cliente_nombre,
				c.estado,
				p.id        AS plazo_id,
				p.tipo      AS plazo_tipo,
				p.fecha_limite
			FROM casos c
			JOIN clientes cl ON cl.id = c.cliente_id
			JOIN bancos   b  ON b.id  = c.banco_id
			JOIN plazos   p  ON p.caso_id = c.id
			WHERE c.estudio_id = $1
			  AND c.banco_id = ANY($2)
			  AND c.estado NOT IN ('TERMINADO','CIERRE')
			  AND p.cumplido = false
			  AND p.fecha_limite <= $3
			ORDER BY c.id, p.fecha_limite ASC
		) sub
		ORDER BY sub.fecha_limite ASC
		LIMIT 50`,
		estudioID, bancoIDs, horizon,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	today := time.Now()
	var result []CasoPorVencer
	for rows.Next() {
		var r CasoPorVencer
		if err := rows.Scan(
			&r.CasoID, &r.BancoID, &r.BancoNombre, &r.NumeroOT,
			&r.ClienteRUT, &r.ClienteNombre, &r.Estado,
			&r.PlazoCritico.ID, &r.PlazoCritico.Tipo, &r.PlazoCritico.FechaLimite,
		); err != nil {
			return nil, err
		}
		remaining := plazo.RemainingBusinessDays(today, r.PlazoCritico.FechaLimite, nil)
		r.PlazoCritico.DiasRestantes = remaining
		r.PlazoCritico.Semaforo = string(plazo.EvaluateSemaforo(remaining, plazo.DefaultThresholds))
		result = append(result, r)
	}
	return result, rows.Err()
}

func (uc *DashboardUseCase) Nuevos(ctx context.Context, estudioID string, bancoIDs []string) ([]CasoNuevo, error) {
	if len(bancoIDs) == 0 {
		return nil, nil
	}
	rows, err := uc.pool.Query(ctx, `
		SELECT c.id, c.banco_id, b.nombre, cl.rut, cl.nombre, c.abogado_id, c.created_at
		FROM casos c
		JOIN clientes cl ON cl.id = c.cliente_id
		JOIN bancos   b  ON b.id  = c.banco_id
		WHERE c.estudio_id = $1
		  AND c.banco_id = ANY($2)
		  AND c.estado = 'INGRESO'
		ORDER BY c.created_at DESC
		LIMIT 50`,
		estudioID, bancoIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CasoNuevo
	for rows.Next() {
		var r CasoNuevo
		if err := rows.Scan(&r.CasoID, &r.BancoID, &r.BancoNombre, &r.ClienteRUT, &r.ClienteNombre, &r.AbogadoID, &r.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

func (uc *DashboardUseCase) Estancados(ctx context.Context, estudioID string, bancoIDs []string, diasSinMovimiento int) ([]CasoEstancado, error) {
	if len(bancoIDs) == 0 {
		return nil, nil
	}
	umbral := time.Now().AddDate(0, 0, -diasSinMovimiento)

	rows, err := uc.pool.Query(ctx, `
		SELECT c.id, c.banco_id, b.nombre, c.numero_ot, cl.rut, cl.nombre, c.estado, c.updated_at
		FROM casos c
		JOIN clientes cl ON cl.id = c.cliente_id
		JOIN bancos   b  ON b.id  = c.banco_id
		WHERE c.estudio_id = $1
		  AND c.banco_id = ANY($2)
		  AND c.estado NOT IN ('TERMINADO','CIERRE')
		  AND c.updated_at < $3
		ORDER BY c.updated_at ASC
		LIMIT 50`,
		estudioID, bancoIDs, umbral,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	now := time.Now()
	var result []CasoEstancado
	for rows.Next() {
		var r CasoEstancado
		if err := rows.Scan(&r.CasoID, &r.BancoID, &r.BancoNombre, &r.NumeroOT, &r.ClienteRUT, &r.ClienteNombre, &r.Estado, &r.UltimoMovimiento); err != nil {
			return nil, err
		}
		r.DiasEstancado = int(now.Sub(r.UltimoMovimiento).Hours() / 24)
		result = append(result, r)
	}
	return result, rows.Err()
}

func (uc *DashboardUseCase) PorAbogado(ctx context.Context, estudioID string, bancoIDs []string) ([]CargaAbogado, error) {
	if len(bancoIDs) == 0 {
		return nil, nil
	}
	rows, err := uc.pool.Query(ctx, `
		SELECT
			u.id,
			u.nombre,
			COUNT(c.id)::int AS total,
			COUNT(CASE WHEN EXISTS (
				SELECT 1 FROM plazos p
				WHERE p.caso_id = c.id AND p.cumplido = false
				  AND p.fecha_limite <= CURRENT_DATE + 7
			) THEN 1 END)::int AS por_vencer,
			COUNT(CASE WHEN EXISTS (
				SELECT 1 FROM plazos p
				WHERE p.caso_id = c.id AND p.cumplido = false
				  AND p.fecha_limite < CURRENT_DATE
			) THEN 1 END)::int AS vencidos
		FROM usuarios u
		JOIN casos c ON c.abogado_id = u.id
		WHERE u.estudio_id = $1
		  AND c.banco_id = ANY($2)
		  AND c.estado NOT IN ('TERMINADO','CIERRE')
		GROUP BY u.id, u.nombre
		ORDER BY total DESC`,
		estudioID, bancoIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CargaAbogado
	for rows.Next() {
		var r CargaAbogado
		if err := rows.Scan(&r.AbogadoID, &r.Nombre, &r.Total, &r.PorVencer, &r.Vencidos); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}
