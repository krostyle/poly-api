ALTER TABLE casos ADD COLUMN resultado_jpl TEXT;
ALTER TABLE casos ADD COLUMN fecha_resolucion_jpl DATE;

CREATE TABLE configuracion_plazos (
  id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  estudio_id     UUID        NOT NULL REFERENCES estudios(id) ON DELETE CASCADE,
  tipo_plazo     TEXT        NOT NULL,
  dias_habiles   INT         NOT NULL CHECK (dias_habiles > 0),
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (estudio_id, tipo_plazo)
);

CREATE INDEX idx_configuracion_plazos_estudio ON configuracion_plazos (estudio_id);
