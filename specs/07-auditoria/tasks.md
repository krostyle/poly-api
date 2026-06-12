# SPEC-07 Auditoría — Tasks

## Estado: ✅ Completado

## Tareas
- [x] `queries/auditoria.sql` (Insert, ListByCaso)
- [x] Implementar `adapters/persistence/audit_logger.go` (AuditLogger Postgres)
- [x] Conectar `AuditLogger` en `cmd/api/main.go` (inyección)
- [x] `GET /v1/casos/:id/historial` — handler completo con join a usuarios
- [x] Verificación: crear caso → entrada en auditoria; `GET /v1/casos/:id/historial` → retorna la entrada
