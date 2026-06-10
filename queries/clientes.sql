-- name: BuscarClientePorRutBanco :one
SELECT id, estudio_id, banco_id, rut, nombre, contacto, created_at
FROM clientes
WHERE banco_id = $1 AND rut = $2 AND estudio_id = $3
LIMIT 1;

-- name: CrearCliente :one
INSERT INTO clientes (estudio_id, banco_id, rut, nombre, contacto)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, estudio_id, banco_id, rut, nombre, contacto, created_at;

-- name: ObtenerClientePorID :one
SELECT id, estudio_id, banco_id, rut, nombre, contacto, created_at
FROM clientes
WHERE id = $1 AND estudio_id = $2
LIMIT 1;
