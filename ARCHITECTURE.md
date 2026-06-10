# poly-api — Arquitectura

## Sistema completo

```
┌─────────────────────────────────────────────────────┐
│  poly-web (Next.js 16 · Vercel)                     │
│  Browser → App Router → Clerk (auth) → lib/api      │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP/JSON  Authorization: Bearer <Clerk JWT>
                       ▼
┌─────────────────────────────────────────────────────┐
│  poly-api (Go · Railway)                            │
│  chi router → middleware (Clerk JWT) → use cases    │
│            → sqlc / Postgres (Neon)                 │
└─────────────────────────────────────────────────────┘
```

### Contrato entre repos

- Auth: poly-web pide JWT a Clerk → lo envía como `Authorization: Bearer`. poly-api verifica con Clerk SDK y extrae `org_id` (= `estudio_id`) + `user_id`.
- Base URL: `NEXT_PUBLIC_API_URL` en poly-web apunta a poly-api.
- Formato: JSON, fechas ISO 8601 (`2025-06-08`), IDs UUID v4.
- Errores: `{ "error": "mensaje" }` con HTTP code apropiado.

### Servicios externos

| Servicio | Rol | Env var clave |
| --- | --- | --- |
| Clerk | Auth + Organizations (estudios) | `CLERK_SECRET_KEY`, `CLERK_ISSUER_URL` |
| Neon | Postgres serverless | `DATABASE_URL` |
| Railway | Deploy de poly-api | `PORT` |
| Vercel | Deploy de poly-web | — |
| Vercel Blob | Documentos (Fase 5) | `BLOB_READ_WRITE_TOKEN` |

### Flujo de una request

```
1. Usuario en poly-web hace acción
2. useAuth().getToken() → JWT de Clerk
3. lib/api/client.ts → fetch a poly-api con Bearer token
4. middleware auth.go → verifica JWT, extrae estudio_id + banco_ids
5. Handler → use case (con scope inyectado)
6. Use case → port (interface) → sqlc adapter
7. Query con WHERE estudio_id = $1 AND banco_id = ANY($2)
8. Resultado → JSON response
```

---

## Clean Architecture

## Regla de dependencia
```
adapters  →  application  →  domain
   ↑               ↑            ↑
frameworks      orquesta      reglas puras
   IO           use cases     sin imports
                              de infra
```
Las flechas apuntan hacia adentro. Nunca al revés.

## Estructura de directorios

```
cmd/
  api/main.go       HTTP server (arranca, inyecta dependencias, escucha)
  worker/main.go    Cron de recálculo de plazos

internal/
  domain/           ← NÚCLEO. Cero imports de infra.
    caso/           Entidad Caso + invariantes
    plazo/          CalcularFechaLimite · DiasHabilesRestantes · EstadoSemaforo
    estado/         Máquina de estados (TRANSICIONES map + Transicionar)
    operacion/      Entidad Operacion
    ports.go        Interfaces: CasoRepository, PlazoRepository,
                    FeriadoProvider, DocumentStorage, AuditLogger

  application/      ← CASOS DE USO. Solo imports de domain/ y puertos.
    casos/          CrearCaso · TransicionarEstado · AsignarAbogado
    plazos/         RecalcularPlazos · EvaluarSemaforo
    dashboard/      Consultas de agregación

  adapters/         ← INFRAESTRUCTURA. Implementa puertos.
    http/
      router.go     chi router + grupos de rutas
      middleware/   auth.go (Clerk JWT) · tenant.go (scope guard)
      handlers/     casos.go · plazos.go · health.go
    persistence/
      sqlc/         Código GENERADO por sqlc (no editar a mano)
      caso_repo.go  Implementa CasoRepository
    feriados/       Implementa FeriadoProvider
    storage/        Implementa DocumentStorage (Vercel Blob)

migrations/         SQL puro para golang-migrate
queries/            SQL puro para sqlc (casos.sql · plazos.sql · ...)
```

## Inyección de dependencias
El `main.go` es el único lugar donde se instancian adaptadores concretos y se
pasan a los use cases como interfaces. Ningún use case conoce la implementación.

## Multi-tenancy
El middleware `adapters/http/middleware/auth.go` extrae `estudio_id` y `banco_ids`
del JWT de Clerk e los inyecta en el contexto. El guard `tenant.go` rechaza
requests sin scope. Todo `CasoRepository` método aplica ese scope.

## Máquina de estados
```
LLAMADA → REVISION → SUSPENSION → PRE_JUDICIALIZACION → JUDICIALIZACION → CIERRE
                                         ↓                    ↑
                                    RESTITUCION ──────────────┘
Cualquier estado (menos CIERRE) → TERMINADO (con motivo obligatorio)
```
Implementada en `internal/domain/estado/machine.go`.
