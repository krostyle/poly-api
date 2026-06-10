# SPEC-07 Auditoría

## Propósito
Registro inmutable (append-only) de toda mutación de casos para cumplimiento legal y trazabilidad.

## Reglas
- La tabla `auditoria` nunca recibe UPDATE ni DELETE
- Todo `AuditLogger.Registrar()` es una inserción atómica
- Si falla el insert de auditoría → la operación principal no se revierte (log de error, no falla silenciosa)
- El `detalle` (JSONB) contiene antes/después del estado relevante

## Acceptance criteria
- Cada llamada a `transicionar_estado`, `crear_caso`, `asignar_abogado` → entrada en `auditoria`
- `GET /v1/casos/:id/historial` → entries de auditoría del caso, ordenadas `created_at DESC`
- La respuesta incluye: `accion`, `detalle`, `created_at`, `usuario` (nombre)
- Solo usuarios del mismo estudio pueden ver el historial

## Dependencias
- SPEC-01 a SPEC-03 (para que haya entradas que auditar)

## Referencias
- `internal/domain/ports.go` — interfaz `AuditLogger`
- Tabla `auditoria`
