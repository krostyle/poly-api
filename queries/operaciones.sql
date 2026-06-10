-- name: CrearOperacion :one
INSERT INTO operaciones (
    id, caso_id, medio_pago, relacion, monto_clp, monto_uf, fecha_op, created_at
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, now()
)
RETURNING *;

-- name: ListarOperacionesPorCaso :many
SELECT * FROM operaciones
WHERE caso_id = $1
ORDER BY fecha_op DESC;
