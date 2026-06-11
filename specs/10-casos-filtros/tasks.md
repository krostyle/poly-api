# SPEC-10 Filtros Casos — Tasks (poly-api)

## Estado: 🔲 Pendiente

## Tareas

### Domain

- [ ] `internal/domain/ports.go` — agregar `Query string` y `AbogadoIDFilter *string` a `CaseFilters`

### Persistence

- [ ] `internal/adapters/persistence/caso_repo.go` — extender `ListRich` con filtros opcionales:
  - `q` → ILIKE en `clientes.nombre` y `clientes.rut`
  - `estado` → comparación exacta
  - `abogado_id` → UUID filter
  - LIMIT / OFFSET
  - `total` via `COUNT(*) OVER()` en el SELECT

### Handler

- [ ] `internal/adapters/http/handlers/casos.go` — leer query params `q`, `estado`, `abogado_id`, `limit`, `offset`
  - Si `abogado_id == "me"` → resolver a `middleware.UsuarioIDFromCtx(ctx)`

### Verificación

- [ ] `go build ./...` sin errores
- [ ] `GET /v1/casos?q=perez` → filtra por nombre/RUT
- [ ] `GET /v1/casos?estado=JUDICIALIZACION` → filtra por estado
- [ ] `GET /v1/casos?abogado_id=me` → filtra por usuario autenticado
- [ ] `GET /v1/casos?limit=5&offset=0` → max 5 resultados, `total` correcto
