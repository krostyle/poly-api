-- name: CrearPlazo :one
INSERT INTO plazos (
    id, caso_id, tipo, fecha_inicio, dias_habiles, fecha_limite, cumplido, created_at
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, false, now()
)
RETURNING *;

-- name: ListarPlazosPorCaso :many
SELECT * FROM plazos
WHERE caso_id = $1
ORDER BY fecha_limite ASC;

-- name: MarcarPlazoCumplido :one
UPDATE plazos
SET cumplido = true, fecha_cumplido = $2
WHERE id = $1
RETURNING *;

-- name: PlazosActivosVencidos :many
SELECT p.*, c.estudio_id, c.banco_id
FROM plazos p
JOIN casos c ON c.id = p.caso_id
WHERE p.cumplido = false
  AND p.fecha_limite < CURRENT_DATE
  AND c.estado NOT IN ('CIERRE', 'TERMINADO');
