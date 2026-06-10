-- ESTUDIOS (tenant principal, espejo de Clerk Organization)
CREATE TABLE estudios (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_org_id TEXT UNIQUE NOT NULL,
    nombre       TEXT NOT NULL,
    rut          TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- BANCOS (un estudio trabaja con varios; aislamiento de datos entre ellos)
CREATE TABLE bancos (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    estudio_id UUID NOT NULL REFERENCES estudios(id),
    nombre     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- USUARIOS (abogados / tramitadores, espejo de Clerk User)
CREATE TABLE usuarios (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_user_id TEXT UNIQUE NOT NULL,
    estudio_id    UUID NOT NULL REFERENCES estudios(id),
    nombre        TEXT NOT NULL,
    email         TEXT NOT NULL,
    rol           TEXT NOT NULL CHECK (rol IN ('ABOGADO', 'TRAMITADOR', 'ADMIN')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Habilitación abogado ↔ banco
CREATE TABLE usuarios_bancos (
    usuario_id UUID NOT NULL REFERENCES usuarios(id),
    banco_id   UUID NOT NULL REFERENCES bancos(id),
    PRIMARY KEY (usuario_id, banco_id)
);

-- CLIENTES AFECTADOS
CREATE TABLE clientes (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    estudio_id UUID NOT NULL REFERENCES estudios(id),
    banco_id   UUID NOT NULL REFERENCES bancos(id),
    rut        TEXT NOT NULL,
    nombre     TEXT NOT NULL,
    contacto   TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- CASOS (entidad central)
CREATE TABLE casos (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    estudio_id      UUID NOT NULL REFERENCES estudios(id),
    banco_id        UUID NOT NULL REFERENCES bancos(id),
    cliente_id      UUID NOT NULL REFERENCES clientes(id),
    abogado_id      UUID REFERENCES usuarios(id),
    numero_ot       TEXT,
    estado          TEXT NOT NULL DEFAULT 'LLAMADA'
                    CHECK (estado IN (
                        'LLAMADA', 'REVISION', 'SUSPENSION',
                        'PRE_JUDICIALIZACION', 'RESTITUCION',
                        'JUDICIALIZACION', 'CIERRE', 'TERMINADO'
                    )),
    fecha_dj        DATE NOT NULL,
    fecha_denuncia  DATE,
    denuncia_valida BOOLEAN NOT NULL DEFAULT false,
    motivo_termino  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- OPERACIONES IMPUGNADAS
CREATE TABLE operaciones (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caso_id    UUID NOT NULL REFERENCES casos(id),
    medio_pago TEXT NOT NULL CHECK (medio_pago IN (
                   'TARJETA_CREDITO', 'TARJETA_DEBITO', 'TRANSFERENCIA', 'CAJERO'
               )),
    relacion   TEXT NOT NULL CHECK (relacion IN (
                   'CUENTA_PROPIA', 'FAMILIAR', 'TERCERO'
               )),
    monto_clp  BIGINT NOT NULL,
    monto_uf   NUMERIC(10, 2),
    fecha_op   DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ÓRDENES DE TRABAJO
CREATE TABLE ordenes_trabajo (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caso_id          UUID NOT NULL REFERENCES casos(id),
    numero           TEXT NOT NULL,
    fecha_asignacion DATE NOT NULL,
    asignado_por     UUID REFERENCES usuarios(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- DOCUMENTOS
CREATE TABLE documentos (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caso_id    UUID NOT NULL REFERENCES casos(id),
    tipo       TEXT NOT NULL CHECK (tipo IN (
                   'CARTOLA', 'EVIDENCIA', 'DJ', 'DENUNCIA',
                   'CARTA_BANCO', 'DEMANDA', 'RESOLUCION', 'OTRO'
               )),
    blob_url   TEXT NOT NULL,
    nombre     TEXT NOT NULL,
    subido_por UUID REFERENCES usuarios(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- PLAZOS
CREATE TABLE plazos (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caso_id        UUID NOT NULL REFERENCES casos(id),
    tipo           TEXT NOT NULL CHECK (tipo IN (
                       'ANALISIS_INTERNO', 'RESTITUCION', 'ASIGNACION',
                       'PRECAUTELAR', 'DEMANDA', 'RESTITUCION_RECHAZO'
                   )),
    fecha_inicio   DATE NOT NULL,
    dias_habiles   INT NOT NULL,
    fecha_limite   DATE NOT NULL,
    cumplido       BOOLEAN NOT NULL DEFAULT false,
    fecha_cumplido DATE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- FERIADOS (calendario chileno)
CREATE TABLE feriados (
    fecha  DATE PRIMARY KEY,
    nombre TEXT NOT NULL
);

-- AUDITORÍA (append-only)
CREATE TABLE auditoria (
    id         BIGSERIAL PRIMARY KEY,
    estudio_id UUID NOT NULL,
    usuario_id UUID,
    caso_id    UUID,
    accion     TEXT NOT NULL,
    detalle    JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Índices
CREATE INDEX idx_casos_scope    ON casos(estudio_id, banco_id);
CREATE INDEX idx_casos_estado   ON casos(estado);
CREATE INDEX idx_casos_abogado  ON casos(abogado_id);
CREATE INDEX idx_plazos_caso    ON plazos(caso_id);
CREATE INDEX idx_plazos_limite  ON plazos(fecha_limite) WHERE cumplido = false;
CREATE INDEX idx_operaciones_caso ON operaciones(caso_id);
CREATE INDEX idx_documentos_caso  ON documentos(caso_id);
