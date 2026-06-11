# SPEC-11 — Tasks

## Backend (poly-api)
- [x] `migrations/004_estados_rediseno.up.sql` — migración de datos + nuevo CHECK constraint
- [x] `migrations/004_estados_rediseno.down.sql` — rollback
- [x] `internal/domain/estado/machine.go` — 12 nuevos estados + grafo de transiciones
- [x] `internal/domain/caso/caso.go` — tipo `MotivoTermino` + `IsValidMotivoTermino()`
- [x] `internal/application/casos/transicionar_estado.go` — validación motivo enum + plazos nuevos estados
- [x] `internal/application/casos/crear_caso.go` — estado inicial INGRESO (era LLAMADA)
- [x] `internal/adapters/http/handlers/casos.go` — `isBadRequest` actualizado
- [x] `queries/casos.sql` — INSERT usa 'INGRESO', filtros NOT IN actualizados
- [x] `internal/application/dashboard/consultas.go` — filtros estado actualizados

## Frontend (poly-web)
- [x] `lib/api/types.ts` — `Estado` (12 valores) + nuevo tipo `MotivoTermino`
- [x] `features/casos/components/EstadoBadge.tsx` — labels + colores por fase
- [x] `features/casos/components/TransicionarEstadoDialog.tsx` — motivo_termino select, grupos por fase, modo corrección
- [x] `features/casos/components/CasosFilters.tsx` — lista ESTADOS actualizada

## Pendiente
- [ ] Ejecutar migración en Railway (`migrate up 4`)
- [ ] Smoke test del flujo completo en producción
