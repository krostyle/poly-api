# poly-api — Harness

## Stack
Go 1.25 · chi v5 · pgx v5 · sqlc · golang-migrate · clerk-sdk-go v2 · godotenv

## Lee antes de codear
1. `ARCHITECTURE.md` — reglas de capas y dirección de dependencias
2. `internal/domain/ports.go` — todas las interfaces
3. `specs/{feature}/spec.md` — requisitos de la feature activa
4. `specs/{feature}/tasks.md` — qué falta hacer

## Reglas de Clean Architecture (innegociables)
- `internal/domain/` no importa NADA de `adapters/`, `net/http`, `database/sql`, SDKs
- `internal/application/` no importa nada de `adapters/`; recibe puertos por inyección
- Los handlers son finos: validar input → llamar use case → serializar respuesta
- Toda la lógica de negocio vive en `domain/` o `application/`, nunca en handlers

## Reglas de datos
- SQL solo en `queries/*.sql` — nunca strings SQL en código Go
- Todo acceso a DB pasa por el scope guard: `WHERE estudio_id = $1 AND banco_id = ANY($2)`
- La tabla `auditoria` es append-only: ningún UPDATE ni DELETE jamás
- Toda mutación de caso llama a `AuditLogger.Registrar()` — sin excepciones

## Reglas de dominio
- Transiciones de estado solo vía `estado.Transicionar()` — nunca directas
- Plazos solo vía `plazo.CalcularFechaLimite()` — nunca aritmética de fechas ad-hoc
- Los umbrales (35 UF, días de plazos) son configurables — no hardcodear

## Flujo de trabajo
- Antes de implementar una feature: leer su `specs/{n}-{feature}/spec.md`
- Al terminar cada tarea: marcarla como `[x]` en `tasks.md`
- Al terminar la feature: actualizar `ROADMAP.md` en la raíz

## Comandos útiles
```
make run          # levanta el servidor en :8080
make test         # go test ./internal/domain/...
make sqlc         # regenera código desde queries/*.sql
make migrate-up   # aplica migraciones pendientes
```
