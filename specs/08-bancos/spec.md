# SPEC-08 Gestión de Bancos y Asignaciones

## Propósito

Permitir a un administrador del estudio crear y nombrar los bancos con los que trabaja, y asignar qué usuarios pueden operar en cada banco. Sin esta feature, ningún usuario puede crear casos.

## Contexto del schema

```
bancos           (id, estudio_id, nombre, created_at)
usuarios_bancos  (usuario_id, banco_id)  ← junction de acceso
```

Un banco pertenece a un estudio específico (no es un catálogo global). Distintos estudios pueden tener su propio "Banco de Chile". El acceso de un usuario a un banco se controla vía `usuarios_bancos`; si un usuario no está en esa tabla para ningún banco, `/v1/me` retorna `bancos: []` y no puede crear casos.

## User stories

- Como admin del estudio, puedo crear un banco (nombre) para que mis usuarios puedan registrar casos en él.
- Como admin, puedo renombrar un banco.
- Como admin, puedo asignar usuarios a un banco para que puedan ver y crear sus casos.
- Como admin, puedo desasignar un usuario de un banco.
- Como admin, puedo ver qué usuarios están asignados a cada banco.
- Como cualquier usuario autenticado, puedo listar los bancos de mi estudio (para selectores en formularios).

## Acceptance criteria

### Bancos

- `POST /v1/bancos` con `{ nombre }` → crea banco bajo el estudio del token, retorna banco creado con 201
  - Solo rol `ADMIN` puede llamar este endpoint → `403` si es ABOGADO o TRAMITADOR
  - `nombre` no vacío → `400` si falta
- `GET /v1/bancos` → lista todos los bancos del estudio del token (todos los roles)
- `PATCH /v1/bancos/:id` con `{ nombre }` → renombra, solo ADMIN → `403` si no
  - Banco debe pertenecer al estudio del token → `404` si no existe
- `DELETE /v1/bancos/:id` → elimina banco, solo ADMIN
  - Si el banco tiene casos → `409 Conflict` con mensaje "El banco tiene casos asociados y no puede eliminarse"
  - Si no tiene casos → elimina también sus entradas en `usuarios_bancos`

### Asignaciones usuario ↔ banco

- `GET /v1/bancos/:id/usuarios` → lista usuarios asignados al banco, solo ADMIN
- `POST /v1/bancos/:id/usuarios` con `{ usuarioId }` → asigna usuario al banco
  - El usuario debe pertenecer al mismo estudio → `404` si no
  - Si ya está asignado → `409` idempotente (o ignorar silenciosamente — ver nota)
  - Solo ADMIN → `403` si no
- `DELETE /v1/bancos/:id/usuarios/:usuarioId` → desasigna usuario del banco, solo ADMIN
  - Si el usuario no estaba asignado → `404`

> **Nota sobre idempotencia**: `POST /v1/bancos/:id/usuarios` puede ser idempotente (devolver `200` si ya existía) para simplificar el cliente. Preferible a `409`.

## Reglas de negocio

- Un banco solo puede ser creado/modificado/eliminado por un usuario con `rol = ADMIN` en el mismo estudio
- Un banco sin casos puede eliminarse; uno con casos no (los casos quedarían huérfanos)
- Un usuario puede estar asignado a cero, uno o varios bancos del mismo estudio
- El scope guard del middleware (`banco_id = ANY(...)`) se construye a partir de `usuarios_bancos` — si un usuario pierde acceso a un banco, sus consultas ya no devuelven esos casos

## API

```
GET    /v1/bancos
POST   /v1/bancos
PATCH  /v1/bancos/:id
DELETE /v1/bancos/:id

GET    /v1/bancos/:id/usuarios
POST   /v1/bancos/:id/usuarios
DELETE /v1/bancos/:id/usuarios/:usuarioId
```

## Tipos de respuesta

```json
// Banco
{ "id": "uuid", "nombre": "Banco de Chile", "created_at": "ISO8601" }

// Usuario asignado (para GET /v1/bancos/:id/usuarios)
{ "id": "uuid", "nombre": "Ana García", "email": "ana@estudio.cl", "rol": "ABOGADO" }
```

## Dependencias

- SPEC-01 (auth + scope) completado ✅
- Tabla `bancos` y `usuarios_bancos` ya existen en la migración 001

## Referencias

- `internal/domain/ports.go` — agregar `BancoRepository`, `UsuarioRepository` (si no existen)
- Tablas: `bancos`, `usuarios_bancos`, `usuarios`
