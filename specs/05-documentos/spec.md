# SPEC-05 Documentos

## Propósito
Subir y gestionar los documentos de cada caso (DJ, denuncia penal, cartola, demanda, etc.) almacenados en Vercel Blob.

## Tipos de documento
`CARTOLA` · `EVIDENCIA` · `DJ` · `DENUNCIA` · `CARTA_BANCO` · `DEMANDA` · `RESOLUCION` · `OTRO`

## Acceptance criteria
- `POST /v1/casos/:id/documentos` (multipart/form-data) → sube a Vercel Blob, guarda metadata en DB
- `GET /v1/casos/:id/documentos` → lista de documentos con `blob_url`, `tipo`, `nombre`, `subido_por`, `created_at`
- El `blob_url` es una URL pública o pre-signed retornada por Vercel Blob
- Solo usuarios del mismo estudio pueden ver/subir documentos del caso
- Tamaño máximo de archivo: 20 MB (rechazar con `413` si excede)
- Tipos MIME permitidos: PDF, imágenes (JPEG, PNG), documentos Office

## Dependencias
- SPEC-01 (auth), SPEC-02 (casos)

## Referencias
- `internal/adapters/storage/vercel_blob.go` — stub existente
- `internal/domain/ports.go` — interfaz `DocumentStorage`
- Tabla `documentos`
