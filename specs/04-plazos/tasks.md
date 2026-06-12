# SPEC-04 Motor de Plazos — Tasks

## Estado: ✅ Completado

## Tareas

### Feriados
- [x] Script de seed de feriados chilenos 2025–2026 (migration o seed SQL)
- [x] Implementar `adapters/feriados/provider.go` (query a tabla feriados)

### Application layer
- [x] `application/plazos/crear_plazos_iniciales.go` (llamado al crear caso)
- [x] Completar `application/plazos/recalcular_plazos.go` (conectar FeriadoProvider real)
- [x] Integrar creación de plazos en `transicionar_estado.go`

### Adapters
- [x] Implementar `PlazoRepository` en `adapters/persistence/plazo_repo.go`

### Handlers + rutas
- [x] `GET /v1/casos/:id/plazos` — lista con diasRestantes + semaforo calculados
- [x] `POST /v1/plazos/:id/cumplir` — marcar cumplido

### Worker (cron)
- [x] Completar `cmd/worker/main.go` con loop diario
- [x] `application/plazos/recalcular_todos.go` — recorre casos activos

### Verificación
- [x] Crear caso → plazos creados automáticamente
- [x] `GET /v1/casos/:id/plazos` → `diasRestantes` y `semaforo` correctos
- [x] Feriado en el rango → fecha_limite se desplaza correctamente
- [x] Cron corre sin errores (log)
