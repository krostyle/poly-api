# SPEC-02 CRUD Casos — Tasks

## Estado: 🔲 Pendiente (requiere SPEC-01)

## Tareas

### SQL + sqlc
- [ ] `queries/clientes.sql` (Upsert, GetByID, GetByRutBanco)
- [ ] Completar `queries/casos.sql` (PATCH, GetDetalle con joins)
- [ ] Completar `queries/operaciones.sql`
- [ ] `make sqlc`

### Application layer
- [ ] Completar `application/casos/crear_caso.go` (dispara creación de plazos iniciales)
- [ ] `application/casos/actualizar_caso.go`
- [ ] `application/clientes/crear_cliente.go`
- [ ] `application/operaciones/agregar_operacion.go`

### Adapters
- [ ] `adapters/persistence/cliente_repo.go`
- [ ] `adapters/persistence/operacion_repo.go`
- [ ] Completar `adapters/persistence/caso_repo.go`

### Handlers + rutas
- [ ] `handlers/clientes.go`
- [ ] Completar `handlers/casos.go` (Listar, Crear, Obtener, Actualizar)
- [ ] `handlers/operaciones.go`
- [ ] Montar rutas en `router.go`

### Verificación
- [ ] `POST /v1/casos` → caso en estado LLAMADA
- [ ] `GET /v1/casos` → solo casos del scope del token
- [ ] `POST /v1/casos/:id/operaciones` → operación creada
- [ ] Caso ajeno al estudio → `404`
