-- name: ObtenerFeriados :many
SELECT fecha FROM feriados
WHERE fecha BETWEEN $1 AND $2
ORDER BY fecha ASC;

-- name: InsertarFeriado :exec
INSERT INTO feriados (fecha, nombre)
VALUES ($1, $2)
ON CONFLICT (fecha) DO UPDATE SET nombre = EXCLUDED.nombre;
