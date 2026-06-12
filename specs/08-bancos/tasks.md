# SPEC-08 Gestión de Bancos — Tasks (poly-api)

## Estado: ✅ Completado

## Tareas

### Domain + ports
- [x] `internal/domain/ports.go` — `BancoRepository` con CRUD + asignaciones
- [x] Tipos `Banco` y `UsuarioBanco` en el dominio

### SQL
- [x] `queries/bancos.sql`:
  - `CrearBanco`
  - `ListarBancosPorEstudio`
  - `ObtenerBancoPorID`
  - `ActualizarBanco`
  - `EliminarBanco`
  - `BancoTieneCasos` (SELECT EXISTS)
  - `ListarUsuariosDeBanco`
  - `AsignarUsuarioABanco` (INSERT ... ON CONFLICT DO NOTHING)
  - `DesasignarUsuarioDeBanco`
  - `UsuarioEstaAsignado`

### Persistence
- [x] `internal/adapters/persistence/banco_repo.go` — implementa `BancoRepository`

### Application
- [x] `internal/application/bancos/crear_banco.go`
- [x] `internal/application/bancos/listar_bancos.go`
- [x] `internal/application/bancos/actualizar_banco.go`
- [x] `internal/application/bancos/eliminar_banco.go`
- [x] `internal/application/bancos/asignar_usuario.go`
- [x] `internal/application/bancos/desasignar_usuario.go`

### Handlers + rutas
- [x] `internal/adapters/http/handlers/bancos.go`
- [x] Montar rutas en `router.go`

### Verificación
- [x] `POST /v1/bancos` con ADMIN → 201
- [x] `POST /v1/bancos` con ABOGADO → 403
- [x] `GET /v1/bancos` → lista del estudio
- [x] `DELETE /v1/bancos/:id` con casos → 409
- [x] `DELETE /v1/bancos/:id` sin casos → 204
- [x] `POST /v1/bancos/:id/usuarios` → usuario aparece en GET
- [x] `DELETE /v1/bancos/:id/usuarios/:uid` → usuario desaparece de GET
- [x] `/v1/me` tras asignación refleja el nuevo banco en `bancos[]`
