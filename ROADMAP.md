# Poly — Roadmap MVP

## Estado actual

**Fase completada: 0 — Scaffolding**

## Fases

### Fase 0 — Scaffolding ✅

- [x] Git + estructura de repos
- [x] poly-api: Clean Architecture, dominio completo, HTTP server
- [x] poly-web: Next.js 16, Clerk, TanStack Query, design system
- [x] Migraciones SQL, sqlc config, Makefile
- [x] Tests de dominio (estado, plazos): 8 passing

### Fase 1 — Auth + Multi-tenancy 🔲

**Objetivo:** Un abogado puede hacer login y ver solo los casos de su estudio/banco.

- [ ] poly-api: middleware Clerk JWT + guard tenant scope
- [ ] poly-api: endpoints de bootstrapping (estudio, banco, usuario)
- [ ] poly-web: flujo sign-in → dashboard con datos reales de Clerk
- **Spec:** `poly-api/specs/01-auth/` · `poly-web/specs/01-auth-ui/`

### Fase 2 — CRUD de casos, operaciones y clientes 🔲

**Objetivo:** Crear casos, registrar operaciones impugnadas y datos del cliente.

- [ ] poly-api: CRUD completo con validaciones de dominio
- [ ] poly-web: formulario de creación + listado de casos
- **Spec:** `poly-api/specs/02-casos/` · `poly-web/specs/02-casos-ui/`

### Fase 3 — Máquina de estados 🔲

**Objetivo:** Transicionar casos entre los 8 estados con validaciones de negocio.

- [ ] poly-api: endpoint de transición + validación (denuncia válida para JUDICIALIZACION, motivo para TERMINADO)
- [ ] poly-web: controles de transición en la vista del caso
- **Spec:** `poly-api/specs/03-estados/`

### Fase 4 — Motor de plazos + semáforo 🔲

**Objetivo:** Calcular fechas límite en días hábiles y mostrar semáforo en tiempo real.

- [ ] poly-api: creación automática de plazos al entrar a cada estado
- [ ] poly-api: cron diario de recálculo
- [ ] poly-web: semáforo visible en dashboard y vista de caso
- **Spec:** `poly-api/specs/04-plazos/` · `poly-web/specs/04-plazos-ui/`

### Fase 5 — Documentos 🔲

**Objetivo:** Subir y ver documentos (DJ, denuncia, demanda, etc.) via Vercel Blob.

- [ ] poly-api: endpoint de upload + metadata en DB
- [ ] poly-web: componente de upload + lista de documentos por caso
- **Spec:** `poly-api/specs/05-documentos/`

### Fase 6 — Dashboard 🔲

**Objetivo:** Vista de control con casos por vencer, nuevos y estancados.

- [ ] poly-api: endpoints de dashboard (por vencer, nuevos, estancados, por abogado)
- [ ] poly-web: dashboard con tablas y semáforos
- **Spec:** `poly-api/specs/06-dashboard/` · `poly-web/specs/03-dashboard-ui/`

### Fase 7 — Auditoría 🔲

**Objetivo:** Registro inmutable de toda mutación de casos.

- [ ] poly-api: AuditLogger adapter (Postgres append-only)
- [ ] poly-web: vista de historial de un caso
- **Spec:** `poly-api/specs/07-auditoria/`

---

## Post-MVP

Generación automática de documentos · Clasificación asistida de indicios · Reportería CMF · Métricas por abogado · Notificaciones email/WhatsApp
