# SPEC-05 Documentos — Tasks

## Estado: ✅ Completado

## Tareas
- [x] `queries/documentos.sql` (Create, ListByCaso) — integrado en queries existentes
- [x] Implementar `adapters/storage/vercel_blob.go`
- [x] `adapters/persistence/documento_repo.go`
- [x] `application/documentos/subir_documento.go`
- [x] Handler multipart `handlers/documentos.go`
- [x] `POST /v1/casos/:id/documentos` + `GET /v1/casos/:id/documentos`
- [x] Verificación: subir PDF → URL en respuesta, metadata en DB
