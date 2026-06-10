-- name: UpsertEstudio :one
INSERT INTO estudios (clerk_org_id, nombre)
VALUES ($1, $2)
ON CONFLICT (clerk_org_id)
DO UPDATE SET nombre = EXCLUDED.nombre
RETURNING *;

-- name: GetEstudioByClerkOrgID :one
SELECT * FROM estudios WHERE clerk_org_id = $1;
