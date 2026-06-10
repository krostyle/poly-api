# SPEC-01 Auth + Multi-tenancy (backend)

## Propósito
Verificar la identidad del usuario (JWT de Clerk) y garantizar que toda request opera
dentro del scope correcto: un usuario solo accede a los datos de su estudio y los bancos
a los que está habilitado.

## User stories
- Como abogado, cuando llamo a la API con mi token de Clerk, el sistema me identifica y
  me da acceso solo a los casos de mi estudio.
- Como sistema, cuando alguien llama sin token o con token inválido, recibo 401.
- Como tramitador habilitado para el Banco X pero no el Banco Y, no puedo ver casos del Banco Y.

## Acceptance criteria

### Middleware de auth
- `GET /v1/casos` sin `Authorization` → `401 Unauthorized`
- `GET /v1/casos` con JWT inválido → `401 Unauthorized`
- `GET /v1/casos` con JWT válido de Clerk → `200` (lista vacía si no hay casos)
- El JWT lleva el `org_id` de Clerk → mapeado a `estudio_id` en la DB
- El `user_id` de Clerk → mapeado a `usuario_id` en la DB

### Guard de tenant scope
- Toda query a `casos`, `plazos`, `operaciones`, `documentos` lleva `WHERE estudio_id = $estudio_id AND banco_id = ANY($banco_ids)`
- Si `estudio_id` no está en el contexto → `401`
- `banco_ids` vacío (usuario sin bancos asignados) → queries retornan listas vacías (no error)

### Bootstrapping de estudio/usuario
- Al primer login de una Org de Clerk, se crea el registro en `estudios` (si no existe)
- Al primer login de un User de Clerk, se crea el registro en `usuarios` (si no existe)
- Un `POST /v1/admin/bancos` crea un banco asociado al estudio del token

## Reglas de negocio
- Un usuario pertenece a exactamente un estudio (su Clerk Organization)
- Un usuario puede estar habilitado para 0..N bancos dentro de su estudio (tabla `usuarios_bancos`)
- El rol (`ABOGADO`, `TRAMITADOR`, `ADMIN`) viene de la DB, no del JWT

## API endpoints de esta fase

```
POST /v1/bootstrap          Crea estudio+usuario si no existen (llamado al primer login)
GET  /v1/me                 Retorna usuario+estudio+bancos del token actual
POST /v1/admin/bancos       Crea un banco (requiere rol ADMIN)
POST /v1/admin/usuarios/:id/bancos  Asigna banco a usuario
```

## Dependencias
- Ninguna (primera feature a implementar)

## Referencias
- `internal/adapters/http/middleware/auth.go` — stub existente
- `internal/adapters/http/middleware/tenant.go` — stub existente
- `internal/domain/ports.go` — interfaces
- Tablas: `estudios`, `usuarios`, `bancos`, `usuarios_bancos`
