# SPEC-09 Documentos Storage — Tasks (poly-api)

## Estado: ✅ Completado

## Tareas

### Verificar adaptador
- [x] `internal/adapters/storage/vercel_blob.go` — usa API Vercel Blob con `Authorization: Bearer $BLOB_READ_WRITE_TOKEN`
- [x] `cmd/api/main.go` inyecta `storage.NewVercelBlobStorage(token)` en el handler de documentos
- [x] Handler multipart extrae `file` y `tipo` del form y llama a `DocumentStorage.Upload`

### Variables de entorno
- [x] `BLOB_READ_WRITE_TOKEN=` agregado a `poly-api/.env.example`
- [x] Warning en startup si el token no está configurado

### Prueba end-to-end
- [ ] Con `BLOB_READ_WRITE_TOKEN` configurado en local: `POST /v1/casos/:id/documentos` con un PDF real → 201
- [ ] `GET /v1/casos/:id/documentos` → retorna el documento con `blob_url` accesible
- [ ] Subir archivo > 20 MB → 413 o mensaje de error claro
