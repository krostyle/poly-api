# SPEC-03 Máquina de Estados

## Propósito
Permitir transicionar casos entre los 8 estados del flujo con validación de reglas de negocio y registro de auditoría.

## Acceptance criteria

### Endpoint de transición
- `POST /v1/casos/:id/transicion` con `{ estado, motivoTermino? }` → transiciona el caso
- Transición inválida (no en la tabla) → `422 Unprocessable Entity` con mensaje descriptivo
- Pasar a `JUDICIALIZACION` sin `denuncia_valida = true` → `422`
- Pasar a `TERMINADO` sin `motivoTermino` → `422`
- Cada transición se registra en `auditoria` con estado anterior y nuevo

### Reglas de transición
```
LLAMADA             → REVISION, TERMINADO
REVISION            → SUSPENSION, TERMINADO
SUSPENSION          → PRE_JUDICIALIZACION, TERMINADO
PRE_JUDICIALIZACION → JUDICIALIZACION*, RESTITUCION, TERMINADO
RESTITUCION         → JUDICIALIZACION*, CIERRE
JUDICIALIZACION     → CIERRE, TERMINADO
CIERRE              → (terminal)
TERMINADO           → (terminal)
* requiere denuncia_valida = true
```

### Efectos secundarios por transición
- Entrada a `SUSPENSION` → crear plazo PRECAUTELAR (13 días desde asignación)
- Tribunal ACOGE (→ `JUDICIALIZACION`) → crear plazo DEMANDA (10 días desde resolución)
- Tribunal RECHAZA (→ `RESTITUCION`) → crear plazo RESTITUCION_RECHAZO (3 días) + DEMANDA (10 días)
- Entrada a `TERMINADO` → marcar todos los plazos pendientes como cumplidos

### Visibilidad del historial
- `GET /v1/casos/:id/historial` → lista de entradas de auditoría del caso (estado anterior/nuevo, usuario, timestamp)

## Dependencias
- SPEC-01, SPEC-02

## Referencias
- `internal/domain/estado/machine.go` — lógica implementada y testeada
- `internal/domain/caso/caso.go` — `ValidarTransicion()`
- `internal/application/casos/transicionar_estado.go` — stub existente
