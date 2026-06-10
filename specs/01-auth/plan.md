# SPEC-01 Auth — Plan de implementación

## Enfoque
1. Completar el middleware `auth.go` con verificación JWT real de Clerk SDK v2
2. Crear el adaptador de persistencia para estudios/usuarios (upsert al login)
3. Añadir queries sqlc para las tablas de auth
4. Implementar handlers de bootstrapping y `/me`

## Archivos a crear / modificar

### Nuevas queries SQL (`queries/`)
- `queries/estudios.sql` — GetByClerkOrgID, Upsert
- `queries/usuarios.sql` — GetByClerkUserID, Upsert, GetConBancos
- `queries/bancos.sql` — Create, List, GetByID
- `queries/usuarios_bancos.sql` — Assign, List

### Adapter de auth (`internal/adapters/http/middleware/auth.go`)
- Usar `github.com/clerk/clerk-sdk-go/v2/http` para `WithHeaderAuthorization`
- Extraer `claims.Subject` (user_id) y `claims.ActiveOrganizationID` (org_id)
- Hacer lookup en DB: `usuarios` por `clerk_user_id`, `estudios` por `clerk_org_id`
- Inyectar `estudio_id`, `usuario_id`, `banco_ids[]` en el contexto

### Application layer (`internal/application/auth/`)
- `bootstrap.go` — BootstrapEstudio(clerkOrgID, nombre) + BootstrapUsuario(clerkUserID, estudioID, email, nombre)

### Handlers (`internal/adapters/http/handlers/`)
- `auth.go` — Bootstrap, Me
- `admin.go` — CrearBanco, AsignarBancoAUsuario

### Router (`internal/adapters/http/router.go`)
- Montar las nuevas rutas

## Orden de implementación
1. Queries SQL → `make sqlc` → código generado
2. Adapters de persistencia (estudios, usuarios, bancos)
3. Application layer (bootstrap)
4. Middleware auth completo (Clerk SDK)
5. Handlers + rutas
6. Tests de integración (opcional en esta fase)

## Riesgos
- Clerk SDK v2 tiene una API ligeramente diferente a v1 — verificar ejemplos en la doc oficial
- El `ActiveOrganizationID` solo está presente si el usuario seleccionó una org activa en el JWT
