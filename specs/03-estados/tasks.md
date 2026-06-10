# SPEC-03 Máquina de Estados — Tasks

## Estado: 🔲 Pendiente (requiere SPEC-02)

## Tareas

### Application layer
- [ ] Completar `application/casos/transicionar_estado.go` (efectos secundarios por estado)
- [ ] Integrar con `application/plazos/` para crear plazos al entrar a cada estado

### Handlers + rutas
- [ ] `POST /v1/casos/:id/transicion` — handler
- [ ] `GET /v1/casos/:id/historial` — handler (desde tabla auditoria)
- [ ] Montar rutas

### Auditoría
- [ ] Implementar `AuditLogger` adapter (Postgres) — adelanto de SPEC-07
- [ ] Conectar en `transicionar_estado.go`

### Verificación
- [ ] Transición válida → caso actualizado
- [ ] Transición inválida → `422`
- [ ] `JUDICIALIZACION` sin denuncia → `422`
- [ ] `TERMINADO` sin motivo → `422`
- [ ] `GET /v1/casos/:id/historial` → entries de auditoría
