# SPEC-07 Auditoría — Tasks

## Estado: 🔲 Pendiente (parcialmente adelantado en SPEC-03)

## Tareas
- [ ] `queries/auditoria.sql` (Insert, ListByCaso)
- [ ] Implementar `adapters/persistence/audit_logger.go` (AuditLogger Postgres)
- [ ] Conectar `AuditLogger` en `cmd/api/main.go` (inyección)
- [ ] `GET /v1/casos/:id/historial` — handler completo con join a usuarios
- [ ] Verificación: crear caso → entrada en auditoria; `GET /v1/casos/:id/historial` → retorna la entrada
