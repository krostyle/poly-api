ALTER TABLE usuarios ADD COLUMN onboarding_completado boolean NOT NULL DEFAULT false;
-- usuarios existentes ya tienen su rol definido, no necesitan onboarding
UPDATE usuarios SET onboarding_completado = true;
