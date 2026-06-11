# SPEC-09 Documentos Storage — Tasks (poly-api)

## Estado: 🔲 Pendiente

## Tareas

### Verificar adaptador

- [ ] Leer `internal/adapters/storage/vercel_blob.go` — confirmar que usa la API `https://blob.vercel-storage.com/` con `Authorization: Bearer $BLOB_READ_WRITE_TOKEN`
- [ ] Confirmar que `cmd/api/main.go` inyecta `storage.NewVercelBlobStorage(token)` en el handler de documentos
- [ ] Confirmar que el handler multipart extrae `file` y `tipo` del form y llama a `DocumentStorage.Upload`

### Variables de entorno

- [ ] Agregar `BLOB_READ_WRITE_TOKEN=` a `poly-api/.env.example`
- [ ] Documentar en Railway cómo obtener el token desde el dashboard de Vercel

### Prueba end-to-end

- [ ] Con `BLOB_READ_WRITE_TOKEN` configurado en local: `POST /v1/casos/:id/documentos` con un PDF real → 201
- [ ] `GET /v1/casos/:id/documentos` → retorna el documento con `blob_url` accesible
- [ ] Subir archivo > 20 MB → 413 o mensaje de error claro

### Ajustes si el adaptador es stub

- [ ] Si `vercel_blob.go` es un stub (retorna `"https://placeholder.url"`), implementar la llamada HTTP real a Vercel Blob
- [ ] Manejar errores de la API de Vercel (token inválido, cuota superada) con mensajes descriptivos
