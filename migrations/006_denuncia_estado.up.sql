-- SPEC-13: Estado de la denuncia + plazo de respuesta
-- Reemplaza denuncia_valida (bool) por estado_denuncia (PENDIENTE/ACOGIDA/RECHAZADA)
-- Agrega RESPUESTA_DENUNCIA al catálogo de tipos de plazo

-- 1. Agregar columna con el enum de estado
ALTER TABLE casos ADD COLUMN estado_denuncia VARCHAR(20) NOT NULL DEFAULT 'PENDIENTE';
ALTER TABLE casos ADD CONSTRAINT casos_estado_denuncia_check
    CHECK (estado_denuncia IN ('PENDIENTE', 'ACOGIDA', 'RECHAZADA'));

-- 2. Migrar datos: denuncia_valida=true → ACOGIDA, el resto queda PENDIENTE
UPDATE casos SET estado_denuncia = 'ACOGIDA' WHERE denuncia_valida = TRUE;

-- 3. Eliminar columna vieja
ALTER TABLE casos DROP COLUMN denuncia_valida;

-- 4. Ampliar el catálogo de plazos con RESPUESTA_DENUNCIA (30 días hábiles desde DJ)
ALTER TABLE plazos DROP CONSTRAINT IF EXISTS plazos_tipo_check;
ALTER TABLE plazos ADD CONSTRAINT plazos_tipo_check
    CHECK (tipo IN (
        'ANALISIS_INTERNO', 'RESTITUCION', 'ASIGNACION',
        'PRECAUTELAR', 'DEMANDA', 'RESTITUCION_RECHAZO',
        'RESPUESTA_DENUNCIA'
    ));
