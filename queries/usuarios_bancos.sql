-- name: AssignBancoToUsuario :exec
INSERT INTO usuarios_bancos (usuario_id, banco_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: ListUsuariosBanco :many
SELECT usuario_id FROM usuarios_bancos WHERE banco_id = $1;
