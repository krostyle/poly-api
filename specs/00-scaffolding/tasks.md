# SPEC-00 Scaffolding — Tasks

> Trabajo inicial de infraestructura. No tiene spec.md ni plan.md porque no es una feature de producto.

## Estado: ✅ Completado

## Tareas realizadas

### Repo raíz
- [x] `git init`
- [x] `.gitignore` (Go, Node, `.env*`)
- [x] `CLAUDE.md`, `ARCHITECTURE.md`, `ROADMAP.md` en raíz
- [x] `CLAUDE.md`, `ARCHITECTURE.md` por repo

### poly-api
- [x] `go mod init poly.app/api`
- [x] Estructura Clean Architecture completa (`cmd/`, `internal/domain/`, `application/`, `adapters/`)
- [x] `internal/domain/estado/machine.go` — máquina de estados + 2 tests
- [x] `internal/domain/plazo/calculator.go` — días hábiles + semáforo + 6 tests
- [x] `internal/domain/caso/caso.go` — entidad + invariantes
- [x] `internal/domain/operacion/operacion.go`
- [x] `internal/domain/ports.go` — todas las interfaces
- [x] `internal/application/casos/` — CrearCaso, TransicionarEstado, AsignarAbogado
- [x] `internal/application/plazos/recalcular_plazos.go`
- [x] `internal/application/dashboard/consultas.go`
- [x] `internal/adapters/http/router.go` + middleware (auth, tenant) + handlers stub
- [x] `cmd/api/main.go` — servidor chi en :8080, `GET /health`
- [x] `cmd/worker/main.go` — stub del cron
- [x] `migrations/001_initial_schema.up/down.sql` — esquema completo
- [x] `queries/casos.sql`, `plazos.sql`, `operaciones.sql`, `feriados.sql`
- [x] `sqlc.yaml`
- [x] `Makefile` (run, test, build, sqlc, migrate-up/down)
- [x] `.env.example`
- [x] Dependencias: chi, pgx, uuid, godotenv → `go build ./...` limpio
- [x] `go test ./internal/domain/...` → 8 tests passing
