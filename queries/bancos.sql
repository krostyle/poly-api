-- name: CreateBanco :one
INSERT INTO bancos (estudio_id, nombre)
VALUES ($1, $2)
RETURNING *;

-- name: ListBancosByEstudio :many
SELECT * FROM bancos WHERE estudio_id = $1 ORDER BY nombre;

-- name: GetBancoByID :one
SELECT * FROM bancos WHERE id = $1 AND estudio_id = $2;
