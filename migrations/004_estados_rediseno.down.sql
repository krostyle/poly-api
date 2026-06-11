-- Revertir rediseño de estados

ALTER TABLE casos DROP CONSTRAINT IF EXISTS casos_motivo_termino_check;
ALTER TABLE casos DROP CONSTRAINT IF EXISTS casos_estado_check;

-- Revertir nuevos estados sin equivalente al más cercano
UPDATE casos SET estado = 'JUDICIALIZACION'
    WHERE estado IN ('AUDIENCIA', 'SENTENCIA', 'APELACION', 'SENTENCIA_SEGUNDA', 'CUMPLIMIENTO');
UPDATE casos SET estado = 'JUDICIALIZACION' WHERE estado = 'JUDICIAL';
UPDATE casos SET estado = 'RESTITUCION'     WHERE estado = 'PAGO_NORMATIVO';
UPDATE casos SET estado = 'SUSPENSION'      WHERE estado = 'PREJUDICIAL';
UPDATE casos SET estado = 'LLAMADA'         WHERE estado = 'INGRESO';
-- CIERRE antiguo ← TERMINADO (los que originalmente llegaron desde CIERRE)
-- No podemos distinguirlos, dejarlos como TERMINADO

ALTER TABLE casos ADD CONSTRAINT casos_estado_check
    CHECK (estado IN (
        'LLAMADA', 'REVISION', 'SUSPENSION',
        'PRE_JUDICIALIZACION', 'RESTITUCION',
        'JUDICIALIZACION', 'CIERRE', 'TERMINADO'
    ));

ALTER TABLE casos ALTER COLUMN estado SET DEFAULT 'LLAMADA';
