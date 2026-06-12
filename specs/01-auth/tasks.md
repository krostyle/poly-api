# SPEC-01 Auth — Tasks

## Estado: ✅ Completado

## Tareas

### SQL + sqlc
- [x] `queries/estudios.sql` (GetByClerkOrgID, Upsert)
- [x] `queries/usuarios.sql` (GetByClerkUserID, Upsert, GetConBancos)
- [x] `queries/bancos.sql` (Create, List, GetByID)
- [x] `queries/usuarios_bancos.sql` (Assign, List)
- [x] `make sqlc` → código generado sin errores

### Application layer
- [x] `internal/application/auth/bootstrap.go` — BootstrapEstudio, BootstrapUsuario

### Adapters de persistencia
- [x] `internal/adapters/persistence/estudio_repo.go`
- [x] `internal/adapters/persistence/usuario_repo.go`
- [x] `internal/adapters/persistence/banco_repo.go`

### Middleware
- [x] `internal/adapters/http/middleware/auth.go` — verificación JWT real con Clerk SDK v2
- [x] Contexto inyecta: `estudio_id`, `usuario_id`, `banco_ids[]`

### Handlers + rutas
- [x] `POST /v1/bootstrap` — handler + ruta
- [x] `GET /v1/me` — handler + ruta
- [x] `POST /v1/admin/bancos` — handler + ruta (guard rol ADMIN)
- [x] `POST /v1/admin/usuarios/:id/bancos` — handler + ruta

### Verificación
- [x] `curl -X GET http://localhost:8080/v1/casos` → `401` sin token
- [x] `curl` con JWT válido → `200`
- [x] Usuario sin bancos → listado vacío, no error
