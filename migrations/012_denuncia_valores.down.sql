ALTER TABLE casos DROP CONSTRAINT IF EXISTS casos_estado_denuncia_check;

UPDATE casos SET estado_denuncia = 'PENDIENTE' WHERE estado_denuncia = 'SOLICITADA';
UPDATE casos SET estado_denuncia = 'ACOGIDA'   WHERE estado_denuncia = 'VALIDA';
UPDATE casos SET estado_denuncia = 'RECHAZADA' WHERE estado_denuncia = 'INVALIDA';
UPDATE casos SET estado_denuncia = 'PENDIENTE' WHERE estado_denuncia = 'SIN_DENUNCIA';

ALTER TABLE casos ALTER COLUMN estado_denuncia SET DEFAULT 'PENDIENTE';
ALTER TABLE casos ADD CONSTRAINT casos_estado_denuncia_check
    CHECK (estado_denuncia IN ('PENDIENTE', 'ACOGIDA', 'RECHAZADA'));
