-- name: CrearCaso :one
INSERT INTO casos (
    id, estudio_id, banco_id, cliente_id, abogado_id,
    numero_ot, estado, fecha_dj, fecha_denuncia, denuncia_valida,
    motivo_termino, created_at, updated_at
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4,
    $5, 'LLAMADA', $6, $7, $8,
    NULL, now(), now()
)
RETURNING *;

-- name: ObtenerCasoPorID :one
SELECT * FROM casos
WHERE id = $1 AND estudio_id = $2
LIMIT 1;

-- name: ListarCasos :many
SELECT * FROM casos
WHERE estudio_id = $1
  AND banco_id = ANY($2::uuid[])
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ActualizarEstadoCaso :one
UPDATE casos
SET estado = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: ActualizarCaso :one
UPDATE casos
SET abogado_id     = $2,
    numero_ot      = $3,
    fecha_denuncia = $4,
    denuncia_valida = $5,
    motivo_termino = $6,
    updated_at     = now()
WHERE id = $1
RETURNING *;

-- name: CasosPorVencer :many
SELECT c.*, p.fecha_limite, p.tipo AS plazo_tipo
FROM casos c
JOIN plazos p ON p.caso_id = c.id
WHERE c.estudio_id = $1
  AND c.banco_id = ANY($2::uuid[])
  AND p.cumplido = false
  AND p.fecha_limite <= (CURRENT_DATE + ($3 || ' days')::interval)
  AND c.estado NOT IN ('CIERRE', 'TERMINADO')
ORDER BY p.fecha_limite ASC;
