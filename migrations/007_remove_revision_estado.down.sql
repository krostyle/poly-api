-- Revertir: restaurar REVISION al catálogo de estados y hacer fecha_dj opcional
ALTER TABLE casos DROP CONSTRAINT IF EXISTS casos_estado_check;
ALTER TABLE casos ADD CONSTRAINT casos_estado_check CHECK (estado IN (
    'INGRESO', 'REVISION', 'PREJUDICIAL', 'PAGO_NORMATIVO', 'JUDICIAL',
    'AUDIENCIA', 'SENTENCIA', 'APELACION', 'SENTENCIA_SEGUNDA',
    'CUMPLIMIENTO', 'TERMINADO', 'CIERRE'
));
ALTER TABLE casos ALTER COLUMN fecha_dj DROP NOT NULL;
