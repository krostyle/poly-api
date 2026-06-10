# SPEC-06 Dashboard (endpoints)

## Propósito
Proveer los endpoints de agregación que alimentan el dashboard del abogado/tramitador.

## Endpoints

### `GET /v1/dashboard/por-vencer`
Casos con al menos un plazo no cumplido venciendo en los próximos N días (default: 7).
Retorna: `[{ caso, plazo_critico: { tipo, fechaLimite, diasRestantes, semaforo } }]` ordenado por `diasRestantes ASC`.

### `GET /v1/dashboard/nuevos`
Casos en estado `LLAMADA` sin abogado asignado o sin actividad en las últimas 24 h.

### `GET /v1/dashboard/estancados`
Casos activos (no CIERRE/TERMINADO) sin transición de estado en más de X días (configurable, default 5).

### `GET /v1/dashboard/por-abogado`
Carga de casos agrupada por abogado: `[{ abogado: { id, nombre }, total, porVencer, vencidos }]`.
Solo accesible para roles `TRAMITADOR` y `ADMIN`.

## Acceptance criteria
- Todos los endpoints respetan el scope `estudio_id + banco_ids`
- Todos los endpoints soportan query param `bancoId` para filtrar por banco específico
- Respuestas paginadas con `?limit=` y `?offset=`

## Dependencias
- SPEC-01, SPEC-02, SPEC-04
