-- SPEC: Eliminar estado REVISION del flujo
-- El estudio siempre recibe casos con DJ, por lo que la revisión del banco
-- ya ocurrió antes de ingresar al sistema. INGRESO pasa directo a PREJUDICIAL.

-- 1. Migrar casos existentes en REVISION → INGRESO
UPDATE casos SET estado = 'INGRESO' WHERE estado = 'REVISION';

-- 2. Reemplazar el CHECK constraint de estado (sin REVISION)
ALTER TABLE casos DROP CONSTRAINT IF EXISTS casos_estado_check;
ALTER TABLE casos ADD CONSTRAINT casos_estado_check CHECK (estado IN (
    'INGRESO', 'PREJUDICIAL', 'PAGO_NORMATIVO', 'JUDICIAL',
    'AUDIENCIA', 'SENTENCIA', 'APELACION', 'SENTENCIA_SEGUNDA',
    'CUMPLIMIENTO', 'TERMINADO', 'CIERRE'
));

-- 3. fecha_dj pasa a ser obligatoria
ALTER TABLE casos ALTER COLUMN fecha_dj SET NOT NULL;
