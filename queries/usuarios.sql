-- name: UpsertUsuario :one
INSERT INTO usuarios (clerk_user_id, estudio_id, nombre, email, rol)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (clerk_user_id)
DO UPDATE SET nombre = EXCLUDED.nombre, email = EXCLUDED.email
RETURNING *;

-- name: GetUsuarioByClerkUserID :one
SELECT * FROM usuarios WHERE clerk_user_id = $1;

-- name: GetBancoIDsByUsuarioID :many
SELECT banco_id FROM usuarios_bancos WHERE usuario_id = $1;
