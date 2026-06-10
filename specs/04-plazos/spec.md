# SPEC-04 Motor de Plazos + Semáforo

## Propósito
Calcular automáticamente las fechas límite legales en días hábiles bancarios y mantener el semáforo actualizado para que ningún plazo se venza por descuido.

## Plazos legales a gestionar

| Tipo | Días hábiles | Cuenta desde |
|---|---|---|
| `ANALISIS_INTERNO` | 5 | fecha_dj (DJ) |
| `RESTITUCION` | 13 | fecha_dj |
| `ASIGNACION` | 7 | fecha_dj |
| `PRECAUTELAR` | 13 (configurable) | fecha de asignación |
| `DEMANDA` | 10 | resolución del tribunal |
| `RESTITUCION_RECHAZO` | 3 | resolución del tribunal |

## Acceptance criteria

### Creación automática de plazos
- Al crear un caso → plazos `ANALISIS_INTERNO`, `RESTITUCION`, `ASIGNACION` con `fecha_limite` calculada
- Al entrar a `SUSPENSION` → plazo `PRECAUTELAR`
- Al entrar a `JUDICIALIZACION` → plazo `DEMANDA`
- Al entrar a `RESTITUCION` → plazos `RESTITUCION_RECHAZO` y `DEMANDA`

### Cálculo de días hábiles
- `CalcularFechaLimite` usa días L–V excluyendo feriados de la tabla `feriados`
- Feriados cargados desde la tabla `feriados` (Postgres), proveídos por `FeriadoProvider`
- Si la tabla `feriados` está vacía, el cálculo usa solo L–V

### Semáforo
- `GET /v1/casos/:id/plazos` retorna cada plazo con `diasRestantes` y `semaforo`
- Umbrales: VENCIDO ≤ 0 · ROJO ≤ 2 · AMARILLO ≤ 5 · VERDE > 5

### Cron diario
- El worker recorre todos los casos activos (no CIERRE/TERMINADO)
- Recalcula `diasRestantes` y `semaforo` para cada plazo no cumplido
- Registra alertas para plazos ROJO y VENCIDO (log por ahora, notificaciones en Fase 2)

### Marcar cumplido
- `POST /v1/plazos/:id/cumplir` marca el plazo como cumplido con fecha actual
- Solo el abogado/tramitador del caso puede marcar plazos

## Dependencias
- SPEC-02 (creación de casos), SPEC-03 (transiciones que crean plazos)
- Tabla `feriados` debe tener datos (al menos año en curso)

## Referencias
- `internal/domain/plazo/calculator.go` — implementado y testeado (6 tests)
- `internal/application/plazos/recalcular_plazos.go` — stub existente
- `internal/adapters/feriados/provider.go` — stub existente
