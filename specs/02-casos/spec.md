# SPEC-02 CRUD Casos, Operaciones y Clientes

## Propósito
Permitir crear y gestionar casos de fraude con sus operaciones impugnadas y los datos del cliente afectado.

## User stories
- Como tramitador, puedo crear un caso con la fecha de DJ y datos del cliente.
- Como tramitador, puedo agregar operaciones impugnadas a un caso (una o varias).
- Como abogado, puedo listar los casos de mi banco filtrados por estado.
- Como abogado, puedo ver el detalle completo de un caso con sus operaciones.

## Acceptance criteria

### Casos
- `POST /v1/casos` con `{ bancoId, clienteId, fechaDj }` → crea caso en estado `LLAMADA`, retorna el caso creado
- `GET /v1/casos` → lista paginada, filtrada por `estudio_id + banco_ids` del token
- `GET /v1/casos/:id` → detalle completo (caso + operaciones + cliente)
- `fechaDj` en el pasado reciente es válida; fechas futuras → `400`
- El `clienteId` debe pertenecer al mismo estudio+banco → `404` si no

### Clientes
- `POST /v1/clientes` con `{ rut, nombre, contacto?, bancoId }` → crea cliente
- `GET /v1/clientes/:id` → detalle del cliente
- Si ya existe un cliente con el mismo `rut` en el mismo banco, retorna el existente (upsert-like)

### Operaciones
- `POST /v1/casos/:id/operaciones` con `{ medioPago, relacion, montoCLP, fechaOp, montoUF? }` → agrega operación
- `GET /v1/casos/:id/operaciones` → lista de operaciones del caso
- `fechaOp` debe estar dentro de los 120 días anteriores a la `fechaDenuncia` del caso (si ya existe)
- `montoCLP` > 0

## Reglas de negocio
- Un caso puede tener N operaciones impugnadas
- El banco es quien marca la `relacion` (CUENTA_PROPIA, FAMILIAR, TERCERO)
- La suma de `montoCLP` de las operaciones determina si el caso supera 35 UF (relevante para restitución)
- Al crear el caso se dispara la creación de los primeros plazos (ANALISIS_INTERNO: 5 días, RESTITUCION: 13 días, ASIGNACION: 7 días) — coordinado con SPEC-04

## API

```
POST   /v1/clientes
GET    /v1/clientes/:id
POST   /v1/casos
GET    /v1/casos
GET    /v1/casos/:id
PATCH  /v1/casos/:id            (actualizar abogadoId, numerOT, denunciaValida, fechaDenuncia)
POST   /v1/casos/:id/operaciones
GET    /v1/casos/:id/operaciones
```

## Dependencias
- SPEC-01 (auth + scope) debe estar completo

## Referencias
- `internal/application/casos/crear_caso.go` — stub existente
- Tablas: `casos`, `clientes`, `operaciones`
