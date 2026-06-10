# poly-api — Arquitectura (Clean Architecture)

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
