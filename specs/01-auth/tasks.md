# SPEC-01 Auth — Tasks

## Estado: 🔲 Pendiente

## Tareas

### SQL + sqlc
- [ ] `queries/estudios.sql` (GetByClerkOrgID, Upsert)
- [ ] `queries/usuarios.sql` (GetByClerkUserID, Upsert, GetConBancos)
- [ ] `queries/bancos.sql` (Create, List, GetByID)
- [ ] `queries/usuarios_bancos.sql` (Assign, List)
- [ ] `make sqlc` → código generado sin errores

### Application layer
- [ ] `internal/application/auth/bootstrap.go` — BootstrapEstudio, BootstrapUsuario

### Adapters de persistencia
- [ ] `internal/adapters/persistence/estudio_repo.go`
- [ ] `internal/adapters/persistence/usuario_repo.go`
- [ ] `internal/adapters/persistence/banco_repo.go`

### Middleware
- [ ] `internal/adapters/http/middleware/auth.go` — verificación JWT real con Clerk SDK v2
- [ ] Contexto inyecta: `estudio_id`, `usuario_id`, `banco_ids[]`

### Handlers + rutas
- [ ] `POST /v1/bootstrap` — handler + ruta
- [ ] `GET /v1/me` — handler + ruta
- [ ] `POST /v1/admin/bancos` — handler + ruta (guard rol ADMIN)
- [ ] `POST /v1/admin/usuarios/:id/bancos` — handler + ruta

### Verificación
- [ ] `curl -X GET http://localhost:8080/v1/casos` → `401` sin token
- [ ] `curl` con JWT válido → `200`
- [ ] Usuario sin bancos → listado vacío, no error
