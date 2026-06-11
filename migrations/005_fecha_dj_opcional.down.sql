-- Revert: set a default date for existing NULLs before re-adding NOT NULL
UPDATE casos SET fecha_dj = created_at::date WHERE fecha_dj IS NULL;
ALTER TABLE casos ALTER COLUMN fecha_dj SET NOT NULL;
