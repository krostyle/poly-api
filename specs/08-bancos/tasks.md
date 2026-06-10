# SPEC-08 Gestión de Bancos — Tasks (poly-api)

## Estado: 🔲 Pendiente

## Tareas

### Domain + ports

- [ ] `internal/domain/ports.go` — agregar `BancoRepository`:
  ```go
  type BancoRepository interface {
      Create(ctx, estudioID, nombre string) (*Banco, error)
      List(ctx, estudioID string) ([]*Banco, error)
      GetByID(ctx, estudioID, id string) (*Banco, error)
      Update(ctx, estudioID, id, nombre string) (*Banco, error)
      Delete(ctx, estudioID, id string) error
      HasCasos(ctx, id string) (bool, error)
      // asignaciones
      ListUsuarios(ctx, bancoID string) ([]*UsuarioBanco, error)
      AsignarUsuario(ctx, bancoID, usuarioID string) error
      DesasignarUsuario(ctx, bancoID, usuarioID string) error
      UsuarioAsignado(ctx, bancoID, usuarioID string) (bool, error)
  }
  ```
- [ ] Agregar tipos `Banco` y `UsuarioBanco` al dominio

### SQL

- [ ] `queries/bancos.sql`:
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

- [ ] `internal/adapters/persistence/banco_repo.go` — implementa `BancoRepository`

### Application

- [ ] `internal/application/bancos/crear_banco.go` — valida rol ADMIN, crea banco
- [ ] `internal/application/bancos/listar_bancos.go` — lista bancos del estudio
- [ ] `internal/application/bancos/actualizar_banco.go` — valida rol ADMIN, renombra
- [ ] `internal/application/bancos/eliminar_banco.go` — valida rol ADMIN, verifica sin casos
- [ ] `internal/application/bancos/asignar_usuario.go` — valida rol ADMIN, valida mismo estudio, asigna
- [ ] `internal/application/bancos/desasignar_usuario.go` — valida rol ADMIN, desasigna

### Handlers + rutas

- [ ] `internal/adapters/http/handlers/bancos.go` — handlers para todos los endpoints
- [ ] Montar rutas en `router.go`

### Verificación

- [ ] `POST /v1/bancos` con ADMIN → 201
- [ ] `POST /v1/bancos` con ABOGADO → 403
- [ ] `GET /v1/bancos` → lista del estudio
- [ ] `DELETE /v1/bancos/:id` con casos → 409
- [ ] `DELETE /v1/bancos/:id` sin casos → 204
- [ ] `POST /v1/bancos/:id/usuarios` → usuario aparece en GET
- [ ] `DELETE /v1/bancos/:id/usuarios/:uid` → usuario desaparece de GET
- [ ] `/v1/me` tras asignación refleja el nuevo banco en `bancos[]`
