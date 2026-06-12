package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/caso"
	"poly.app/api/internal/domain/estado"
)

type CasoRepo struct {
	pool *pgxpool.Pool
}

func NewCasoRepo(pool *pgxpool.Pool) *CasoRepo {
	return &CasoRepo{pool: pool}
}

func (r *CasoRepo) Create(ctx context.Context, c *caso.Caso) error {
	const q = `INSERT INTO casos
		(id, estudio_id, banco_id, cliente_id, abogado_id, numero_ot, estado,
		 fecha_dj, fecha_denuncia, estado_denuncia, motivo_termino, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, q,
		c.ID, c.EstudioID, c.BancoID, c.ClienteID, c.AbogadoID, c.NumeroOT,
		string(c.Estado), c.FechaDJ, c.FechaDenuncia, string(c.EstadoDenuncia),
		c.MotivoTermino, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *CasoRepo) GetByID(ctx context.Context, estudioID, id string) (*caso.Caso, error) {
	const q = `SELECT id, estudio_id, banco_id, cliente_id, abogado_id, numero_ot, estado,
		fecha_dj, fecha_denuncia, estado_denuncia, motivo_termino, created_at, updated_at
		FROM casos WHERE id = $1 AND estudio_id = $2`
	row := r.pool.QueryRow(ctx, q, id, estudioID)
	return scanCaso(row)
}

func (r *CasoRepo) Update(ctx context.Context, c *caso.Caso) error {
	const q = `UPDATE casos SET abogado_id=$2, numero_ot=$3, fecha_denuncia=$4,
		estado_denuncia=$5, motivo_termino=$6, fecha_dj=$7, updated_at=now() WHERE id=$1`
	_, err := r.pool.Exec(ctx, q,
		c.ID, c.AbogadoID, c.NumeroOT, c.FechaDenuncia, string(c.EstadoDenuncia), c.MotivoTermino, c.FechaDJ,
	)
	return err
}

func (r *CasoRepo) List(ctx context.Context, estudioID string, filters domain.CaseFilters) ([]*caso.Caso, error) {
	if len(filters.BancoIDs) == 0 {
		return nil, nil
	}
	limit := filters.Limit
	if limit == 0 {
		limit = 50
	}
	const q = `SELECT id, estudio_id, banco_id, cliente_id, abogado_id, numero_ot, estado,
		fecha_dj, fecha_denuncia, estado_denuncia, motivo_termino, created_at, updated_at
		FROM casos WHERE estudio_id = $1 AND banco_id = ANY($2::uuid[])
		ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	rows, err := r.pool.Query(ctx, q, estudioID, filters.BancoIDs, limit, filters.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*caso.Caso
	for rows.Next() {
		c, err := scanCasoRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *CasoRepo) ListRich(ctx context.Context, estudioID string, filters domain.CaseFilters) ([]*domain.CasoListItem, int, error) {
	if len(filters.BancoIDs) == 0 {
		return nil, 0, nil
	}
	limit := filters.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	bancoScope := filters.BancoIDs
	if filters.BancoIDFilter != "" {
		bancoScope = []string{filters.BancoIDFilter}
	}

	countQ := `SELECT COUNT(*)
		FROM casos c
		JOIN clientes cl ON cl.id = c.cliente_id
		WHERE c.estudio_id = $1
		  AND c.banco_id   = ANY($2::uuid[])
		  AND ($3 = ''   OR c.estado::text = $3)
		  AND ($4::uuid IS NULL OR c.abogado_id = $4::uuid)
		  AND ($5 = ''   OR cl.nombre ILIKE '%' || $5 || '%' OR cl.rut ILIKE '%' || $5 || '%')
		  AND (NOT $6 OR c.estado::text != 'CIERRE')`

	var totalRaw int64
	estFilter := ""
	if filters.Estado != nil {
		estFilter = string(*filters.Estado)
	}
	var abogadoFilter interface{} = nil
	if filters.AbogadoID != nil {
		abogadoFilter = *filters.AbogadoID
	}
	if err := r.pool.QueryRow(ctx, countQ, estudioID, bancoScope, estFilter, abogadoFilter, filters.Query, filters.ExcluirCierre).Scan(&totalRaw); err != nil {
		return nil, 0, err
	}
	total := int(totalRaw)

	q := `SELECT c.id, c.banco_id, b.nombre, c.cliente_id, cl.rut, cl.nombre,
		c.abogado_id, c.numero_ot, c.estado, c.fecha_dj, c.estado_denuncia, c.created_at
		FROM casos c
		JOIN bancos   b  ON b.id  = c.banco_id
		JOIN clientes cl ON cl.id = c.cliente_id
		WHERE c.estudio_id = $1
		  AND c.banco_id   = ANY($2::uuid[])
		  AND ($3 = ''   OR c.estado::text = $3)
		  AND ($4::uuid IS NULL OR c.abogado_id = $4::uuid)
		  AND ($5 = ''   OR cl.nombre ILIKE '%' || $5 || '%' OR cl.rut ILIKE '%' || $5 || '%')
		  AND (NOT $6 OR c.estado::text != 'CIERRE')
		ORDER BY c.created_at DESC LIMIT $7 OFFSET $8`
	rows, err := r.pool.Query(ctx, q, estudioID, bancoScope, estFilter, abogadoFilter, filters.Query, filters.ExcluirCierre, limit, filters.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*domain.CasoListItem
	for rows.Next() {
		var item domain.CasoListItem
		var est, estDenuncia string
		var fechaDJ *time.Time
		if err := rows.Scan(
			&item.ID, &item.BancoID, &item.BancoNombre,
			&item.ClienteID, &item.ClienteRUT, &item.ClienteNombre,
			&item.AbogadoID, &item.NumeroOT, &est,
			&fechaDJ, &estDenuncia, &item.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		item.Estado = estado.Estado(est)
		item.FechaDJ = fechaDJ
		item.EstadoDenuncia = caso.EstadoDenuncia(estDenuncia)
		items = append(items, &item)
	}
	return items, total, rows.Err()
}

func (r *CasoRepo) UpdateState(ctx context.Context, id string, newState estado.Estado) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE casos SET estado=$2, updated_at=now() WHERE id=$1`,
		id, string(newState),
	)
	return err
}

func (r *CasoRepo) GetDetalle(ctx context.Context, estudioID, id string) (*domain.CasoDetalle, error) {
	c, err := r.GetByID(ctx, estudioID, id)
	if err != nil {
		return nil, err
	}

	var cliente domain.Cliente
	const clienteQ = `SELECT id, estudio_id, banco_id, rut, nombre, contacto, created_at
		FROM clientes WHERE id = $1`
	row := r.pool.QueryRow(ctx, clienteQ, c.ClienteID)
	if err := row.Scan(&cliente.ID, &cliente.EstudioID, &cliente.BancoID,
		&cliente.RUT, &cliente.Nombre, &cliente.Contacto, &cliente.CreatedAt); err != nil {
		return nil, err
	}

	const opQ = `SELECT id, caso_id, medio_pago, relacion, monto_clp, monto_uf, fecha_op, created_at
		FROM operaciones WHERE caso_id = $1 ORDER BY fecha_op DESC`
	rows, err := r.pool.Query(ctx, opQ, id)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &domain.CasoDetalle{Caso: c, Cliente: &cliente, Operaciones: ops}, nil
}

func (r *CasoRepo) Delete(ctx context.Context, estudioID, casoID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	for _, table := range []string{"plazos", "documentos", "auditoria", "operaciones"} {
		if _, err := tx.Exec(ctx, "DELETE FROM "+table+" WHERE caso_id = $1", casoID); err != nil {
			return fmt.Errorf("delete from %s: %w", table, err)
		}
	}

	tag, err := tx.Exec(ctx, "DELETE FROM casos WHERE id = $1 AND estudio_id = $2", casoID, estudioID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("not found")
	}
	return tx.Commit(ctx)
}

type casoScanner interface {
	Scan(dest ...any) error
}

func scanCaso(row casoScanner) (*caso.Caso, error) {
	var c caso.Caso
	var est, estDenuncia string
	var fechaDJ *time.Time
	var createdAt, updatedAt time.Time
	err := row.Scan(
		&c.ID, &c.EstudioID, &c.BancoID, &c.ClienteID, &c.AbogadoID, &c.NumeroOT,
		&est, &fechaDJ, &c.FechaDenuncia, &estDenuncia, &c.MotivoTermino,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	c.Estado = estado.Estado(est)
	c.FechaDJ = fechaDJ
	c.EstadoDenuncia = caso.EstadoDenuncia(estDenuncia)
	c.CreatedAt = createdAt
	c.UpdatedAt = updatedAt
	return &c, nil
}

type pgxRows interface {
	Scan(dest ...any) error
}

func scanCasoRow(row pgxRows) (*caso.Caso, error) {
	return scanCaso(row)
}
