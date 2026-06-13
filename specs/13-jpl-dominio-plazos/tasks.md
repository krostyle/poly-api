# SPEC-13 — Tasks

## Estado: ✅ Completado (implementado sin spec previo — documentado retroactivamente 2026-06-13)

## Backend (poly-api)

### Dominio
- [x] `internal/domain/caso/caso.go` — EstadoDenuncia: SOLICITADA/VALIDA/INVALIDA/SIN_DENUNCIA
- [x] `internal/domain/caso/caso.go` — ResultadoJPL type con 4 constantes + IsValidResultadoJPL
- [x] `internal/domain/caso/caso.go` — campos ResultadoJPL *ResultadoJPL, FechaResolucionJPL *time.Time en Caso
- [x] `internal/domain/estado/machine.go` — transición INGRESO→JUDICIAL (caso excepcional)
- [x] `internal/domain/plazo/calculator.go` — TipoResolucionJPL = "RESOLUCION_JPL" (3 días)

### Puertos
- [x] `internal/domain/ports.go` — ConfiguracionPlazo struct + ConfiguracionPlazoRepository interface

### Aplicación
- [x] `internal/application/casos/crear_caso.go` — createInitialPlazos desde fecha_dj + plazos configurables por estudio
- [x] `internal/application/casos/actualizar_caso.go` — recalcularPlazosDJ cuando cambia fecha_dj
- [x] `internal/application/casos/transicionar_estado.go` — PREJUDICIAL crea PRECAUTELAR (13d) + RESOLUCION_JPL (3d)
- [x] `internal/application/casos/transicionar_estado.go` — eliminar validateDenunciaGuard (guards basados en denuncia incorrectos)
- [x] `internal/application/casos/transicionar_estado.go` — eliminar ErrDenunciaRechazadaRequerida, ErrDenunciaAcogidaRequerida, ErrFechaDenunciaRequerida, ErrDenunciaPendienteInvalida

### Persistencia
- [x] `internal/adapters/persistence/caso_repo.go` — SELECT/UPDATE resultado_jpl, fecha_resolucion_jpl
- [x] `internal/adapters/persistence/caso_repo.go` — scanCaso incluye resultadoJPL y FechaResolucionJPL
- [x] `internal/adapters/persistence/configuracion_plazo_repo.go` — NEW: GetByEstudio, Upsert

### HTTP
- [x] `internal/adapters/http/handlers/casos.go` — ResultadoJPL en casoJSON + Actualizar + isBadRequest limpio
- [x] `internal/adapters/http/handlers/configuracion.go` — NEW: Listar, Actualizar (validación 1-90 días)
- [x] `internal/adapters/http/router.go` — GET /v1/configuracion/plazos + PUT /v1/configuracion/plazos/{tipo}

### Migraciones
- [x] `migrations/011_jpl_configuracion.up.sql` — ADD resultado_jpl, fecha_resolucion_jpl; CREATE configuracion_plazos
- [x] `migrations/012_denuncia_valores.up.sql` — migrar PENDIENTE→SOLICITADA, ACOGIDA→VALIDA, RECHAZADA→INVALIDA + nuevo CHECK

## Frontend (poly-web)

### Tipos y API
- [x] `lib/api/types.ts` — EstadoDenuncia (4 valores), ResultadoJPL (4 valores), RESOLUCION_JPL en TipoPlazo, ConfiguracionPlazo
- [x] `lib/api/casos.ts` — RawCaso + mapCasoDetalle con resultado_jpl, fecha_resolucion_jpl
- [x] `lib/api/configuracion.ts` — NEW: listarConfiguracionPlazos, actualizarConfiguracionPlazo

### Hooks
- [x] `features/configuracion/hooks/useConfiguracionPlazos.ts` — NEW
- [x] `features/configuracion/hooks/useActualizarConfiguracionPlazo.ts` — NEW

### Páginas
- [x] `app/(app)/configuracion/plazos/page.tsx` — NEW: página admin con inputs inline por plazo configurable

### Componentes
- [x] `features/casos/components/CasoFlowView.tsx` — NEW: timeline vertical del flujo del caso
- [x] `features/casos/components/CasoDetalleView.tsx` — campos ResultadoJPL + FechaResolucionJPL
- [x] `features/casos/components/CasoDetalleView.tsx` — CasoFlowView integrado entre datos judiciales y operaciones
- [x] `features/casos/components/CasoDetalleView.tsx` — EstadoDenuncia actualizado (4 valores + labels)
- [x] `features/casos/components/CasoDetalleView.tsx` — región bloqueada (read-only) cuando hay tribunal seleccionado
- [x] `features/casos/components/CasoDetalleView.tsx` — resultado JPL en layout vertical full-width
- [x] `features/casos/components/TransicionarEstadoDialog.tsx` — INGRESO→JUDICIAL en TRANSICIONES
- [x] `features/casos/components/TransicionarEstadoDialog.tsx` — guards correctos basados en resultado_jpl (informativo, no bloqueante)
- [x] `features/casos/components/TransicionarEstadoDialog.tsx` — prop estadoDenuncia → resultadoJpl
- [x] `features/casos/components/TransicionarEstadoDialog.tsx` — textos actualizados: medida precautoria es obligatoria
- [x] `features/plazos/components/PlazosCard.tsx` — RESOLUCION_JPL: "Resolución JPL"
- [x] `features/plazos/components/PlazosTable.tsx` — RESOLUCION_JPL: "Resolución JPL"
- [x] `components/layout/Sidebar.tsx` — /configuracion/plazos en adminNav
