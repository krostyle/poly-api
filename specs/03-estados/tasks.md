# SPEC-03 Máquina de Estados — Tasks

## Estado: ✅ Completado

## Tareas

### Application layer
- [x] Completar `application/casos/transicionar_estado.go` (efectos secundarios por estado)
- [x] Integrar con `application/plazos/` para crear plazos al entrar a cada estado

### Handlers + rutas
- [x] `POST /v1/casos/:id/transicion` — handler
- [x] `GET /v1/casos/:id/historial` — handler (desde tabla auditoria)
- [x] Montar rutas

### Auditoría
- [x] Implementar `AuditLogger` adapter (Postgres) — adelanto de SPEC-07
- [x] Conectar en `transicionar_estado.go`

### Verificación
- [x] Transición válida → caso actualizado
- [x] Transición inválida → `422`
- [x] `JUDICIAL` sin denuncia → `422`
- [x] `TERMINADO` sin motivo → `422`
- [x] `GET /v1/casos/:id/historial` → entries de auditoría
