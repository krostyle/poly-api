# SPEC-09 Documentos — Vercel Blob Storage

## Propósito

Conectar el backend de documentos a un proveedor real de almacenamiento de archivos. Actualmente el handler `POST /v1/casos/:id/documentos` existe pero el adaptador `storage/vercel_blob.go` puede ser un stub o no estar configurado correctamente para producción en Railway.

## Contexto

El frontend `DocumentosCard` y el hook `useSubirDocumento` ya están implementados. El backend tiene:
- `internal/adapters/http/handlers/documentos.go` — handler multipart
- `internal/adapters/storage/vercel_blob.go` — adaptador Vercel Blob
- `internal/domain/ports.go` — interfaz `DocumentStorage`

El objetivo es verificar que el flujo completo funciona en producción: el archivo sube a Vercel Blob, la URL queda guardada en la tabla `documentos`, y el cliente puede abrirla.

## Acceptance criteria

- `POST /v1/casos/:id/documentos` recibe un `multipart/form-data` con campos `file` y `tipo`
- El archivo se sube a Vercel Blob y retorna una URL pública
- Se persiste un registro en `documentos` con `{caso_id, tipo, blob_url, nombre, subido_por}`
- `GET /v1/casos/:id/documentos` lista los documentos del caso con `{id, tipo, nombre, blob_url, created_at}`
- Archivos hasta 20 MB aceptados
- Si `BLOB_READ_WRITE_TOKEN` no está configurado → error 500 descriptivo en logs, no silencioso
- Solo usuarios con acceso al caso (mismo estudio + banco asignado) pueden subir/listar

## Variables de entorno requeridas

```
BLOB_READ_WRITE_TOKEN=vercel_blob_rw_...   # desde Vercel Blob dashboard
```

Debe agregarse al `.env.example` y al panel de Railway.

## Verificación

1. Subir un PDF desde `DocumentosCard` en `/casos/:id`
2. El archivo aparece en la lista con su nombre y tipo
3. El ícono ExternalLink abre el archivo en una nueva pestaña
4. Intentar subir un archivo > 20 MB → mensaje de error en UI
5. `GET /v1/casos/:id/documentos` en Postman retorna el documento subido

## Dependencias

- Cuenta Vercel Blob activa (o cualquier proveedor S3-compatible si se cambia el adaptador)
- `BLOB_READ_WRITE_TOKEN` disponible en Railway
