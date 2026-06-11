package persistence

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/domain"
)

type ClienteRepo struct {
	pool *pgxpool.Pool
}

func NewClienteRepo(pool *pgxpool.Pool) *ClienteRepo {
	return &ClienteRepo{pool: pool}
}

func (r *ClienteRepo) Upsert(ctx context.Context, in domain.UpsertClienteInput) (*domain.Cliente, error) {
	existing, err := r.findByRutBanco(ctx, in.EstudioID, in.BancoID, in.RUT)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	const q = `INSERT INTO clientes (estudio_id, banco_id, rut, nombre, contacto)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, estudio_id, banco_id, rut, nombre, contacto, created_at`
	row := r.pool.QueryRow(ctx, q, in.EstudioID, in.BancoID, in.RUT, in.Nombre, in.Contacto)
	return scanCliente(row)
}

func (r *ClienteRepo) Update(ctx context.Context, estudioID, id, nombre string, contacto *string) (*domain.Cliente, error) {
	const q = `UPDATE clientes SET nombre = $1, contacto = $2
		WHERE id = $3 AND estudio_id = $4
		RETURNING id, estudio_id, banco_id, rut, nombre, contacto, created_at`
	row := r.pool.QueryRow(ctx, q, nombre, contacto, id, estudioID)
	return scanCliente(row)
}

func (r *ClienteRepo) GetByID(ctx context.Context, estudioID, id string) (*domain.Cliente, error) {
	const q = `SELECT id, estudio_id, banco_id, rut, nombre, contacto, created_at
		FROM clientes WHERE id = $1 AND estudio_id = $2`
	row := r.pool.QueryRow(ctx, q, id, estudioID)
	return scanCliente(row)
}

func (r *ClienteRepo) findByRutBanco(ctx context.Context, estudioID, bancoID, rut string) (*domain.Cliente, error) {
	const q = `SELECT id, estudio_id, banco_id, rut, nombre, contacto, created_at
		FROM clientes WHERE banco_id = $1 AND rut = $2 AND estudio_id = $3 LIMIT 1`
	row := r.pool.QueryRow(ctx, q, bancoID, rut, estudioID)
	return scanCliente(row)
}

type clienteScanner interface {
	Scan(dest ...any) error
}

func scanCliente(row clienteScanner) (*domain.Cliente, error) {
	var c domain.Cliente
	err := row.Scan(&c.ID, &c.EstudioID, &c.BancoID, &c.RUT, &c.Nombre, &c.Contacto, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
