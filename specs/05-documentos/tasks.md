# SPEC-05 Documentos — Tasks

## Estado: 🔲 Pendiente (requiere SPEC-02)

## Tareas
- [ ] `queries/documentos.sql` (Create, ListByCaso)
- [ ] Implementar `adapters/storage/vercel_blob.go` (SDK de Vercel Blob)
- [ ] `adapters/persistence/documento_repo.go`
- [ ] `application/documentos/subir_documento.go`
- [ ] Handler multipart `handlers/documentos.go`
- [ ] `POST /v1/casos/:id/documentos` + `GET /v1/casos/:id/documentos`
- [ ] Verificación: subir PDF → URL en respuesta, metadata en DB
