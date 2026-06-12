# SPEC-02 CRUD Casos — Tasks

## Estado: ✅ Completado

## Tareas

### SQL + sqlc
- [x] `queries/clientes.sql` (Upsert, GetByID, GetByRutBanco)
- [x] Completar `queries/casos.sql` (PATCH, GetDetalle con joins)
- [x] Completar `queries/operaciones.sql`
- [x] `make sqlc`

### Application layer
- [x] Completar `application/casos/crear_caso.go` (dispara creación de plazos iniciales)
- [x] `application/casos/actualizar_caso.go`
- [x] `application/clientes/crear_cliente.go`
- [x] `application/operaciones/agregar_operacion.go`

### Adapters
- [x] `adapters/persistence/cliente_repo.go`
- [x] `adapters/persistence/operacion_repo.go`
- [x] Completar `adapters/persistence/caso_repo.go`

### Handlers + rutas
- [x] `handlers/clientes.go`
- [x] Completar `handlers/casos.go` (Listar, Crear, Obtener, Actualizar)
- [x] `handlers/operaciones.go`
- [x] Montar rutas en `router.go`

### Verificación
- [x] `POST /v1/casos` → caso en estado INGRESO
- [x] `GET /v1/casos` → solo casos del scope del token
- [x] `POST /v1/casos/:id/operaciones` → operación creada
- [x] Caso ajeno al estudio → `404`
