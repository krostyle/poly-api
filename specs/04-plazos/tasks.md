# SPEC-04 Motor de Plazos — Tasks

## Estado: 🔲 Pendiente (requiere SPEC-02, SPEC-03)

## Tareas

### Feriados
- [ ] Script de seed de feriados chilenos 2025–2026 (migration o seed SQL)
- [ ] Implementar `adapters/feriados/provider.go` (query a tabla feriados)

### Application layer
- [ ] `application/plazos/crear_plazos_iniciales.go` (llamado al crear caso)
- [ ] Completar `application/plazos/recalcular_plazos.go` (conectar FeriadoProvider real)
- [ ] Integrar creación de plazos en `transicionar_estado.go` (SUSPENSION, JUDICIALIZACION, RESTITUCION)

### Adapters
- [ ] Implementar `PlazoRepository` en `adapters/persistence/plazo_repo.go`

### Handlers + rutas
- [ ] `GET /v1/casos/:id/plazos` — lista con diasRestantes + semaforo calculados
- [ ] `POST /v1/plazos/:id/cumplir` — marcar cumplido

### Worker (cron)
- [ ] Completar `cmd/worker/main.go` con loop diario
- [ ] `application/plazos/recalcular_todos.go` — recorre casos activos

### Verificación
- [ ] Crear caso → 3 plazos creados automáticamente
- [ ] `GET /v1/casos/:id/plazos` → `diasRestantes` y `semaforo` correctos
- [ ] Feriado en el rango → fecha_limite se desplaza correctamente
- [ ] Cron corre sin errores (log)
