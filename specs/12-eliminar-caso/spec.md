# SPEC-12 — Eliminar caso y filtro "Incluir cerrados"

## Contexto

Los casos en estado INGRESO corresponden a registros recién creados que aún no han iniciado
flujo judicial. Si fueron creados por error, no existe forma de eliminarlos hoy.
Adicionalmente, los casos en estado CIERRE (administrativamente finalizados) saturan la lista
principal porque nunca desaparecen.

## Funcionalidades

### 1. Eliminar caso (hard delete)

**Reglas de negocio:**
- Solo disponible para usuarios con `rol = 'ADMIN'`.
- Solo si `estado = 'INGRESO'` (el caso no ha avanzado en el flujo).
- Elimina permanentemente el caso y todos sus registros dependientes:
  `plazos → documentos → auditoria → operaciones → caso`
- No hay papelera de reciclaje — la acción es irreversible.

**Endpoint:** `DELETE /v1/casos/:id`
- 204 No Content — eliminado correctamente.
- 403 Forbidden — usuario no es ADMIN.
- 409 Conflict — caso no está en INGRESO.
- 404 Not Found — caso no existe o pertenece a otro estudio.

**UI:**
- Botón "Eliminar caso" visible en la cabecera de `CasoDetalleView` solo cuando:
  `rol === 'ADMIN'` Y `estado === 'INGRESO'`
- Clic abre diálogo de confirmación con mensaje de advertencia.
- Al confirmar: llama al endpoint, luego redirige a `/casos`.

### 2. Filtro "Incluir cerrados" (ExcluirCierre)

**Regla de negocio:**
- Por defecto los casos `CIERRE` se ocultan de la lista (estado administrativo final).
- Un toggle "Incluir cerrados" en `CasosFilters` permite verlos.
- El parámetro se pasa como `excluir_cierre=true|false` en la query string.
- La selección de estado específico (ej: "Cierre") siempre muestra esos casos independiente del toggle.

**Endpoint:** `GET /v1/casos?excluir_cierre=true`
- `CaseFilters.ExcluirCierre bool` nuevo campo.
- Condición SQL: `AND (NOT $N OR c.estado::text != 'CIERRE')`.

## Decisiones de diseño

- Hard delete (no soft delete / archivado) para simplificar el modelo. Solo INGRESO porque
  en ese estado no hay plazos legales activos ni historial judicial relevante.
- ExcluirCierre por defecto = `true` (comportamiento limpio por defecto para casos activos).
- No se elimina el cliente — puede tener otros casos.
