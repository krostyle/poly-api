-- Rediseño del modelo de estados según flujo real Ley 20.009
-- Nuevo orden: INGRESO → REVISION → PREJUDICIAL → PAGO_NORMATIVO → JUDICIAL
--              → AUDIENCIA → SENTENCIA → APELACION → SENTENCIA_SEGUNDA
--              → CUMPLIMIENTO → TERMINADO → CIERRE
--
-- TERMINADO ahora precede a CIERRE (corregido respecto al modelo original).
-- motivo_termino pasa a ser un enum con 12 valores definidos por la experta.

-- 1. Migrar valores existentes al nuevo esquema de nombres
UPDATE casos SET estado = 'INGRESO'        WHERE estado = 'LLAMADA';
UPDATE casos SET estado = 'PREJUDICIAL'    WHERE estado IN ('SUSPENSION', 'PRE_JUDICIALIZACION');
UPDATE casos SET estado = 'PAGO_NORMATIVO' WHERE estado = 'RESTITUCION';
UPDATE casos SET estado = 'JUDICIAL'       WHERE estado = 'JUDICIALIZACION';
-- Antiguo CIERRE equivale al nuevo TERMINADO (fin con motivo específico)
UPDATE casos SET estado = 'TERMINADO'      WHERE estado = 'CIERRE';
-- REVISION y TERMINADO no cambian de nombre

-- 2. Limpiar motivo_termino que no encajen en el nuevo enum
--    (texto libre de versiones anteriores del sistema)
UPDATE casos
SET motivo_termino = NULL
WHERE motivo_termino IS NOT NULL
  AND motivo_termino NOT IN (
    'IMPROCEDENTE', 'EXTEMPORANEO', 'BUSQUEDAS_NEGATIVAS', 'DEUDOR_FALLECIDO',
    'DESISTIMIENTO_CLIENTE', 'DESISTIMIENTO_BANCO',
    'DESISTIMIENTO_DENUNCIA_INVALIDA', 'DESISTIMIENTO_SIN_DENUNCIA',
    'SENTENCIA_FAVORABLE_BANCO', 'SENTENCIA_DESFAVORABLE_BANCO',
    'AVENIMIENTO', 'ABANDONO_PROCEDIMIENTO'
  );

-- 3. Reemplazar constraint de estado
ALTER TABLE casos DROP CONSTRAINT IF EXISTS casos_estado_check;
ALTER TABLE casos ADD CONSTRAINT casos_estado_check
    CHECK (estado IN (
        'INGRESO', 'REVISION', 'PREJUDICIAL', 'PAGO_NORMATIVO',
        'JUDICIAL', 'AUDIENCIA', 'SENTENCIA', 'APELACION',
        'SENTENCIA_SEGUNDA', 'CUMPLIMIENTO', 'TERMINADO', 'CIERRE'
    ));

-- 4. Actualizar valor por defecto para casos nuevos
ALTER TABLE casos ALTER COLUMN estado SET DEFAULT 'INGRESO';

-- 5. Agregar constraint enum para motivo_termino
ALTER TABLE casos ADD CONSTRAINT casos_motivo_termino_check
    CHECK (motivo_termino IS NULL OR motivo_termino IN (
        'IMPROCEDENTE', 'EXTEMPORANEO', 'BUSQUEDAS_NEGATIVAS', 'DEUDOR_FALLECIDO',
        'DESISTIMIENTO_CLIENTE', 'DESISTIMIENTO_BANCO',
        'DESISTIMIENTO_DENUNCIA_INVALIDA', 'DESISTIMIENTO_SIN_DENUNCIA',
        'SENTENCIA_FAVORABLE_BANCO', 'SENTENCIA_DESFAVORABLE_BANCO',
        'AVENIMIENTO', 'ABANDONO_PROCEDIMIENTO'
    ));
