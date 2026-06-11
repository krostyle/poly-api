-- fecha_dj puede ser registrada después del ingreso del reclamo
ALTER TABLE casos ALTER COLUMN fecha_dj DROP NOT NULL;
