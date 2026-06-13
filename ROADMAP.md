# Poly — Roadmap MVP

## Estado actual

**Última fase completada: 13 — JPL, Dominio Ley 20.009 y Correcciones de Plazos**

## Fases

### Fase 0 — Scaffolding ✅

- [x] Git + estructura de repos
- [x] poly-api: Clean Architecture, dominio completo, HTTP server
- [x] poly-web: Next.js 16, Clerk, TanStack Query, design system
- [x] Migraciones SQL, sqlc config, Makefile
- [x] Tests de dominio (estado, plazos): passing

### Fase 1 — Auth + Multi-tenancy ✅

**Objetivo:** Un abogado puede hacer login y ver solo los casos de su estudio/banco.

- [x] poly-api: middleware Clerk JWT + guard tenant scope
- [x] poly-api: endpoints de bootstrapping (estudio, banco, usuario)
- [x] poly-web: flujo sign-in → dashboard con datos reales de Clerk
- **Spec:** `poly-api/specs/01-auth/`

### Fase 2 — CRUD de casos, operaciones y clientes ✅

**Objetivo:** Crear casos, registrar operaciones impugnadas y datos del cliente.

- [x] poly-api: CRUD completo con validaciones de dominio
- [x] poly-web: formulario de creación + listado de casos
- **Spec:** `poly-api/specs/02-casos/`

### Fase 2.5 — Gestión de bancos y asignaciones ✅

**Objetivo:** Que un admin pueda crear bancos y asignar usuarios a ellos, desbloqueando la creación de casos.

- [x] poly-api: CRUD bancos + asignación usuario↔banco con guard de rol ADMIN
- [x] poly-web: página de configuración de bancos + dialogs de asignación
- **Spec:** `poly-api/specs/08-bancos/`

### Fase 3 — Máquina de estados ✅

**Objetivo:** Transicionar casos entre los 12 estados del flujo Ley 20.009 con validaciones de negocio.

- [x] poly-api: endpoint de transición + validación (denuncia para JUDICIAL, motivo para TERMINADO)
- [x] poly-api: RBAC — guard de rol por operación
- [x] poly-web: controles de transición en la vista del caso
- **Spec:** `poly-api/specs/03-estados/` · `poly-api/specs/11-estados-rediseno/`

### Fase 4 — Motor de plazos + semáforo ✅

**Objetivo:** Calcular fechas límite en días hábiles y mostrar semáforo en tiempo real.

- [x] poly-api: creación automática de plazos al entrar a cada estado
- [x] poly-api: cron diario de recálculo
- [x] poly-web: semáforo visible en dashboard y vista de caso
- **Spec:** `poly-api/specs/04-plazos/`

### Fase 5 — Documentos ✅

**Objetivo:** Subir y ver documentos (DJ, denuncia, demanda, etc.) via Vercel Blob.

- [x] poly-api: endpoint de upload + metadata en DB
- [x] poly-web: componente de upload + lista de documentos por caso
- **Spec:** `poly-api/specs/05-documentos/` · `poly-api/specs/09-documentos-storage/`

### Fase 6 — Dashboard ✅

**Objetivo:** Vista de control con casos por vencer, nuevos y estancados.

- [x] poly-api: endpoints de dashboard (por vencer, nuevos, estancados, por abogado)
- [x] poly-web: dashboard con tablas y semáforos
- **Spec:** `poly-api/specs/06-dashboard/`

### Fase 7 — Auditoría ✅

**Objetivo:** Registro inmutable de toda mutación de casos.

- [x] poly-api: AuditLogger adapter (Postgres append-only)
- [x] poly-web: vista de historial de un caso
- **Spec:** `poly-api/specs/07-auditoria/`

### Fase 8 — Filtros y búsqueda de casos ✅

**Objetivo:** Buscar y filtrar casos por nombre, RUT, estado, abogado, banco.

- [x] poly-api: filtros `q`, `estado`, `banco_id`, `abogado_id`, paginación con `total`
- [x] poly-web: controles de filtro en listado de casos
- **Spec:** `poly-api/specs/10-casos-filtros/`

### Fase 9 — Eliminar caso ✅ (backend) / ⏳ (frontend pendiente)

**Objetivo:** ADMIN puede eliminar un caso (con borrado en cascada) y filtrar casos cerrados.

- [x] poly-api: `DELETE /v1/casos/:id` con transacción en cascada
- [x] poly-api: filtro `excluir_cierre` en listado
- [ ] poly-web: botón eliminar + toggle "Incluir cerrados"
- **Spec:** `poly-api/specs/12-eliminar-caso/`

### Fase 10 — JPL, Dominio Ley 20.009 y Correcciones de Plazos ✅

**Objetivo:** Modelar correctamente el flujo Ley 20.009, rastrear resoluciones JPL, configurar plazos por estudio y visualizar el flujo de cada caso.

- [x] poly-api: `EstadoDenuncia` corregido — `SOLICITADA/VALIDA/INVALIDA/SIN_DENUNCIA`
- [x] poly-api: `ResultadoJPL` — `ACEPTA_SUSPENSION/RECHAZA_SUSPENSION/FALLO_FAVORABLE/FALLO_DESFAVORABLE`
- [x] poly-api: plazo `RESOLUCION_JPL` (3 días) creado al entrar a PREJUDICIAL
- [x] poly-api: `INGRESO → JUDICIAL` habilitado (tribunales que exigen demanda + medida precautoria simultáneas)
- [x] poly-api: guards de denuncia eliminados — el camino JUDICIAL/PAGO_NORMATIVO lo determina el JPL
- [x] poly-api: recálculo automático de plazos de ingreso cuando cambia `fecha_dj`
- [x] poly-api: plazos configurables por estudio (`ANALISIS_INTERNO`, `ASIGNACION`) con tabla `configuracion_plazos`
- [x] poly-api: migraciones 011 (JPL + configuración) y 012 (valores denuncia)
- [x] poly-web: `CasoFlowView` — timeline vertical del flujo con estados, plazos y decisión JPL
- [x] poly-web: página `/configuracion/plazos` para admin
- [x] poly-web: región bloqueada cuando hay tribunal seleccionado
- [x] poly-web: resultado JPL full-width (labels no se cortan)
- **Spec:** `poly-api/specs/13-jpl-dominio-plazos/`

---

## Post-MVP

Generación automática de documentos · Clasificación asistida de indicios · Reportería CMF · Métricas por abogado · Notificaciones email/WhatsApp
