# SPEC-13 — JPL, Dominio Ley 20.009 y Correcciones de Plazos

## Contexto

Implementación retroactiva documentada tras las sesiones de trabajo del 12-13 junio 2026.
Cubre correcciones de dominio legal, nuevas funcionalidades de seguimiento JPL, configuración
de plazos por estudio, visualización del flujo de cada caso, y correcciones de bugs detectados
en producción.

---

## 1. Correcciones de dominio Ley 20.009

### 1.1 EstadoDenuncia — valores correctos

Los valores anteriores (`PENDIENTE`, `ACOGIDA`, `RECHAZADA`) eran incorrectos respecto al
proceso real de Ley 20.009. Valores correctos:

| Valor | Significado |
|---|---|
| `SOLICITADA` | Denuncia presentada al banco, esperando respuesta (valor por defecto) |
| `VALIDA` | El banco reconoció la denuncia como válida |
| `INVALIDA` | El banco encontró la denuncia inválida |
| `SIN_DENUNCIA` | No se presentó denuncia |

**Migración:** `012_denuncia_valores` convierte datos existentes y actualiza el CHECK constraint.

### 1.2 Flujo PREJUDICIAL — la medida precautoria es obligatoria

La medida precautoria ante el JPL es un **deber**, no una opción, cuando el caso está en el estudio.
El estado PREJUDICIAL siempre se inicia al recibir el caso.

### 1.3 Quién determina el camino JUDICIAL vs PAGO_NORMATIVO

**Antes (incorrecto):** el `estado_denuncia` determinaba el camino.
**Correcto:** es la **resolución del JPL** sobre la suspensión del monto reclamado:

| Resolución JPL | Significado | Camino |
|---|---|---|
| `RECHAZA_SUSPENSION` | El JPL rechaza mantener el monto suspendido → banco debe devolver | `PAGO_NORMATIVO` |
| `ACEPTA_SUSPENSION` | El JPL acepta la suspensión → banco debe demandar formalmente en ~10 días | `JUDICIAL` |

Los guards basados en `estado_denuncia` fueron eliminados del backend.

### 1.4 Caso excepcional: tribunales que exigen demanda + medida precautoria simultáneas

Algunos tribunales requieren presentar la demanda al mismo tiempo que la medida precautoria
para cumplir con los plazos. Para esto se habilitó la transición `INGRESO → JUDICIAL` en
la máquina de estados.

### 1.5 Plazos INGRESO y la fecha_dj

Los plazos `ANALISIS_INTERNO`, `RESTITUCION`, `ASIGNACION`, `RESPUESTA_DENUNCIA` se calculan
desde `fecha_dj` (no desde `now()`). Si `fecha_dj` cambia, estos plazos se **recalculan
automáticamente** (solo los no cumplidos).

---

## 2. Seguimiento de resolución JPL

### 2.1 Campos nuevos en `casos`

```sql
ALTER TABLE casos ADD COLUMN resultado_jpl TEXT;
ALTER TABLE casos ADD COLUMN fecha_resolucion_jpl DATE;
```

### 2.2 Tipo ResultadoJPL

```go
type ResultadoJPL string

const (
    JPLAceptaSuspension  ResultadoJPL = "ACEPTA_SUSPENSION"
    JPLRechazaSuspension ResultadoJPL = "RECHAZA_SUSPENSION"
    JPLFalloFavorable    ResultadoJPL = "FALLO_FAVORABLE"
    JPLFalloDesfavorable ResultadoJPL = "FALLO_DESFAVORABLE"
)
```

Los primeros dos son la resolución de la medida precautoria. Los últimos dos son el fallo
definitivo del JPL (post-audiencia/sentencia).

### 2.3 Plazo RESOLUCION_JPL

Tipo nuevo `RESOLUCION_JPL` (3 días hábiles). Se crea automáticamente al entrar a PREJUDICIAL,
junto con `PRECAUTELAR` (13 días).

---

## 3. Plazos configurables por estudio

### 3.1 Tabla nueva `configuracion_plazos`

```sql
CREATE TABLE configuracion_plazos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  estudio_id UUID NOT NULL REFERENCES estudios(id) ON DELETE CASCADE,
  tipo_plazo TEXT NOT NULL,
  dias_habiles INT NOT NULL CHECK (dias_habiles > 0),
  UNIQUE (estudio_id, tipo_plazo)
);
```

### 3.2 Plazos configurables

Solo dos plazos son configurables (no tienen base legal fija):

| Tipo | Default |
|---|---|
| `ANALISIS_INTERNO` | 5 días hábiles |
| `ASIGNACION` | 7 días hábiles |

Los demás tienen base legal fija (Ley 20.009, Art. 5 mod. Ley 21.673/2024) y no son configurables.

### 3.3 Endpoints

- `GET /v1/configuracion/plazos` — lista los configurables con valor actual del estudio (ADMIN)
- `PUT /v1/configuracion/plazos/{tipo}` — actualiza un plazo configurable (ADMIN, 1-90 días)

### 3.4 Plazos legales fijos (referencia)

| Tipo | Días | Base |
|---|---|---|
| `RESTITUCION` | 10 hábiles | Ley 20.009 Art. 5, mod. Ley 21.673/2024 |
| `PRECAUTELAR` | 13 hábiles | Plazo para resolución de medida precautoria |
| `RESOLUCION_JPL` | 3 hábiles | Plazo para que el JPL resuelva |
| `DEMANDA` | 10 hábiles | Plazo del banco para demandar tras aceptación de suspensión |
| `RESTITUCION_RECHAZO` | 3 hábiles | Plazo tras rechazo de suspensión |
| `RESPUESTA_DENUNCIA` | 30 hábiles | Plazo del banco para responder la denuncia |

---

## 4. Visualización del flujo del caso (CasoFlowView)

Componente `CasoFlowView` visible en la vista de detalle de cada caso.
Muestra un timeline vertical con los estados del caso, sus plazos y el punto de decisión JPL.

### Comportamiento

- **Completado** (verde): estados ya visitados, con fecha de entrada extraída del historial de auditoría
- **Actual** (navy): estado corriente con badge "Actual"
- **Futuro** (gris): estados aún no alcanzados
- `onlyIfVisited`: PREJUDICIAL, PAGO_NORMATIVO, APELACION, SENTENCIA_SEGUNDA, CIERRE solo
  aparecen cuando han sido visitados o son el estado actual
- **BranchIndicator**: bloque de decisión que aparece tras PREJUDICIAL cuando el caso está
  en ese estado o ya lo pasó. Muestra las dos opciones (acepta/rechaza suspensión) y resalta
  la tomada cuando `resultado_jpl` está registrado

### Fuentes de datos

- `useCaso(casoId)` → estado actual, resultado_jpl
- `usePlazosCaso(casoId)` → chips de plazo con semáforo
- `useHistorialCaso(casoId)` → fechas de entrada a cada estado (entradas `ESTADO_CAMBIADO`)

React Query deduplica las requests (el mismo dato ya se carga en el resto de la vista).

---

## 5. Página de configuración de plazos (admin)

Ruta: `/configuracion/plazos`

- Solo visible en el sidebar para rol ADMIN
- Lista los dos plazos configurables con input numérico inline
- Indica si el valor es el default o si fue personalizado por el estudio
- Muestra card informativa con los plazos legales fijos (no editables)

---

## 6. Correcciones de bugs

### 6.1 Transición INGRESO → PREJUDICIAL retornaba 500

Los errores `ErrFechaDenunciaRequerida` y `ErrDenunciaPendienteInvalida` no estaban en
`isBadRequest`, por lo que retornaban 500 en vez de 422. Corregido al eliminar esos guards
(ya que la medida precautoria no depende del estado de la denuncia).

### 6.2 fecha_dj no recalculaba plazos

Al actualizar `fecha_dj`, los plazos existentes de tipo INGRESO no se recalculaban.
`UpdateCaseUseCase.Execute` ahora detecta el cambio y llama a `recalcularPlazosDJ`.

---

## 7. Correcciones de UI

### 7.1 Región bloqueada cuando hay tribunal seleccionado

Cuando el campo `tribunal` está lleno, el select de `región` se reemplaza por texto de solo
lectura. La región queda determinada por el tribunal y no puede desincronizarse.

### 7.2 Resultado JPL full-width

El select de `resultado_jpl` cambió de layout horizontal (con max-w-52 insuficiente) a layout
vertical (label arriba, select full-width abajo) para que los labels no se corten.

---

## Archivos modificados / creados

### poly-api

| Archivo | Cambio |
|---|---|
| `internal/domain/caso/caso.go` | EstadoDenuncia + ResultadoJPL types |
| `internal/domain/estado/machine.go` | INGRESO→JUDICIAL; RESTITUCION 5→10 |
| `internal/domain/plazo/calculator.go` | TipoResolucionJPL |
| `internal/domain/ports.go` | ConfiguracionPlazo + repo; recalc support |
| `internal/application/casos/crear_caso.go` | createInitialPlazos desde fecha_dj; configurable |
| `internal/application/casos/actualizar_caso.go` | recalcularPlazosDJ cuando cambia fecha_dj |
| `internal/application/casos/transicionar_estado.go` | PREJUDICIAL crea PRECAUTELAR+RESOLUCION_JPL; elimina validateDenunciaGuard |
| `internal/adapters/persistence/caso_repo.go` | SELECT/UPDATE resultado_jpl, fecha_resolucion_jpl |
| `internal/adapters/persistence/configuracion_plazo_repo.go` | NEW — GetByEstudio, Upsert |
| `internal/adapters/http/handlers/casos.go` | ResultadoJPL en request/response; isBadRequest limpio |
| `internal/adapters/http/handlers/configuracion.go` | NEW — Listar, Actualizar |
| `internal/adapters/http/router.go` | rutas /v1/configuracion/plazos |
| `migrations/011_jpl_configuracion.up.sql` | NEW — resultado_jpl, fecha_resolucion_jpl, configuracion_plazos |
| `migrations/012_denuncia_valores.up.sql` | NEW — migra EstadoDenuncia a valores correctos |

### poly-web

| Archivo | Cambio |
|---|---|
| `lib/api/types.ts` | EstadoDenuncia, ResultadoJPL, TipoPlazo, ConfiguracionPlazo types |
| `lib/api/casos.ts` | mapeo resultado_jpl, fecha_resolucion_jpl |
| `lib/api/configuracion.ts` | NEW — listarConfiguracionPlazos, actualizarConfiguracionPlazo |
| `features/configuracion/hooks/useConfiguracionPlazos.ts` | NEW |
| `features/configuracion/hooks/useActualizarConfiguracionPlazo.ts` | NEW |
| `app/(app)/configuracion/plazos/page.tsx` | NEW — página admin |
| `features/casos/components/CasoFlowView.tsx` | NEW — timeline vertical del flujo |
| `features/casos/components/CasoDetalleView.tsx` | ResultadoJPL fields; CasoFlowView integrado; región bloqueada por tribunal; JPL full-width |
| `features/casos/components/TransicionarEstadoDialog.tsx` | Guards corregidos (JPL-based); INGRESO→JUDICIAL |
| `features/plazos/components/PlazosCard.tsx` | RESOLUCION_JPL label |
| `features/plazos/components/PlazosTable.tsx` | RESOLUCION_JPL label |
| `components/layout/Sidebar.tsx` | Link /configuracion/plazos en adminNav |
