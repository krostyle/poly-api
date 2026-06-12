# SPEC-10 Filtros Casos — Tasks (poly-api)

## Estado: ✅ Completado

## Tareas

### Domain
- [x] `internal/domain/ports.go` — `Query string` y `AbogadoIDFilter *string` en `CaseFilters`

### Persistence
- [x] `internal/adapters/persistence/caso_repo.go` — `ListRich` con filtros opcionales:
  - `q` → ILIKE en `clientes.nombre` y `clientes.rut`
  - `estado` → comparación exacta
  - `banco_id` → UUID filter
  - `abogado_id` → UUID filter
  - LIMIT / OFFSET
  - `total` via `COUNT(*) OVER()` en el SELECT

### Handler
- [x] `internal/adapters/http/handlers/casos.go` — query params `q`, `estado`, `banco_id`, `abogado_id`, `limit`, `offset`
  - Si `abogado_id == "me"` → resuelve a `middleware.UsuarioIDFromCtx(ctx)`

### Verificación
- [x] `go build ./...` sin errores
- [x] `GET /v1/casos?q=perez` → filtra por nombre/RUT
- [x] `GET /v1/casos?estado=JUDICIAL` → filtra por estado
- [x] `GET /v1/casos?abogado_id=me` → filtra por usuario autenticado
- [x] `GET /v1/casos?limit=5&offset=0` → max 5 resultados, `total` correcto
