-- Rollback SPEC-13: restaurar denuncia_valida y quitar RESPUESTA_DENUNCIA

ALTER TABLE casos ADD COLUMN denuncia_valida BOOLEAN NOT NULL DEFAULT false;
UPDATE casos SET denuncia_valida = TRUE WHERE estado_denuncia = 'ACOGIDA';
ALTER TABLE casos DROP CONSTRAINT IF EXISTS casos_estado_denuncia_check;
ALTER TABLE casos DROP COLUMN estado_denuncia;

ALTER TABLE plazos DROP CONSTRAINT IF EXISTS plazos_tipo_check;
ALTER TABLE plazos ADD CONSTRAINT plazos_tipo_check
    CHECK (tipo IN (
        'ANALISIS_INTERNO', 'RESTITUCION', 'ASIGNACION',
        'PRECAUTELAR', 'DEMANDA', 'RESTITUCION_RECHAZO'
    ));
