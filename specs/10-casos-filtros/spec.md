# SPEC-10 Filtros en GET /v1/casos

## Propósito

Extender el endpoint `GET /v1/casos` para aceptar parámetros de búsqueda y filtrado, necesario para la UI de SPEC-06 (poly-web).

## Cambios

### Query params nuevos

```
GET /v1/casos?q=&estado=&abogado_id=&limit=50&offset=0
```

| Param | Tipo | Descripción |
|---|---|---|
| `q` | string | Búsqueda ILIKE en `clientes.nombre` y `clientes.rut` |
| `estado` | string | Filtro exacto por estado del caso |
| `abogado_id` | string | UUID del abogado; si se pasa `"me"`, el backend lo resuelve al `usuario_id` del token |
| `limit` | int | Máximo de resultados (default 50, max 200) |
| `offset` | int | Para paginación |

El scope guard de `banco_ids` ya aplica antes de los filtros adicionales.

### Cambios en domain/ports.go

```go
type CaseFilters struct {
    BancoIDs        []string
    Estado          *estado.Estado
    AbogadoID       *string
    Query           string      // búsqueda por nombre/RUT cliente
    Limit           int
    Offset          int
}
```

### Cambios en persistence/caso_repo.go

Extender la query de `ListRich` con condiciones opcionales:

```sql
AND ($n = '' OR (cl.nombre ILIKE '%' || $n || '%' OR cl.rut ILIKE '%' || $n || '%'))
AND ($n = '' OR c.estado::text = $n)
AND ($n::uuid IS NULL OR c.abogado_id = $n::uuid)
LIMIT $n OFFSET $n
```

### Respuesta

```json
{ "casos": [...], "total": 142 }
```

`total` refleja el conteo sin paginar (usando `COUNT(*) OVER()` o una query separada) para que el cliente pueda mostrar "Mostrando 1-50 de 142".

## Acceptance criteria

- `?q=perez` → solo casos con cliente cuyo nombre o RUT contenga "perez" (case-insensitive)
- `?estado=JUDICIALIZACION` → solo esos casos
- `?abogado_id=me` → solo casos del usuario autenticado
- `?limit=10&offset=10` → segunda página de 10 resultados
- Sin params → comportamiento actual (todos los casos del scope)
- `total` en la respuesta refleja el count real tras filtros

## Dependencias

- SPEC-02 (CRUD Casos) ✅
- poly-web SPEC-06 (usa este endpoint)
