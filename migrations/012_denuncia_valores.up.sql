-- Actualiza los valores de estado_denuncia:
-- PENDIENTE → SOLICITADA, ACOGIDA → VALIDA, RECHAZADA → INVALIDA

ALTER TABLE casos DROP CONSTRAINT IF EXISTS casos_estado_denuncia_check;

UPDATE casos SET estado_denuncia = 'SOLICITADA' WHERE estado_denuncia = 'PENDIENTE';
UPDATE casos SET estado_denuncia = 'VALIDA'     WHERE estado_denuncia = 'ACOGIDA';
UPDATE casos SET estado_denuncia = 'INVALIDA'   WHERE estado_denuncia = 'RECHAZADA';

ALTER TABLE casos ALTER COLUMN estado_denuncia SET DEFAULT 'SOLICITADA';
ALTER TABLE casos ADD CONSTRAINT casos_estado_denuncia_check
    CHECK (estado_denuncia IN ('SOLICITADA', 'VALIDA', 'INVALIDA', 'SIN_DENUNCIA'));
