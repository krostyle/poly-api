DROP TABLE IF EXISTS configuracion_plazos;
ALTER TABLE casos DROP COLUMN IF EXISTS resultado_jpl;
ALTER TABLE casos DROP COLUMN IF EXISTS fecha_resolucion_jpl;
