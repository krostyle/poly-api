# SPEC-12 — Tasks

## Estado: ✅ Completado

## Backend (poly-api)
- [x] `internal/domain/ports.go` — `Delete` en `CasoRepository`, `ExcluirCierre bool` en `CaseFilters`
- [x] `internal/adapters/persistence/caso_repo.go` — impl `Delete` (tx: plazos→documentos→auditoria→operaciones→caso) + filtro `ExcluirCierre` en `ListRich`
- [x] `internal/adapters/http/handlers/casos.go` — handler `Eliminar` (ADMIN, cualquier estado) + parseo `excluir_cierre` en `Listar`
- [x] `internal/adapters/http/router.go` — ruta `DELETE /{id}` en `/v1/casos`

## Frontend (poly-web)
- [ ] `lib/api/casos.ts` — fn `eliminarCaso`, campo `excluirCierre` en `CasoFilters`, parámetro en `listarCasos`
- [ ] `features/casos/hooks/useEliminarCaso.ts` — mutation hook
- [ ] `features/casos/components/CasosFilters.tsx` — toggle "Incluir cerrados", emit `excluirCierre` (default false = excluir)
- [ ] `features/casos/components/CasoDetalleView.tsx` — botón eliminar + diálogo de confirmación (ADMIN)
