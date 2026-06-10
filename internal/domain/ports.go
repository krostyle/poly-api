package domain

import (
	"context"
	"time"

	"poly.app/api/internal/domain/caso"
	"poly.app/api/internal/domain/estado"
	"poly.app/api/internal/domain/plazo"
)

// ── Auth domain types ────────────────────────────────────────────────────────

type Estudio struct {
	ID         string
	ClerkOrgID string
	Nombre     string
	RUT        *string
	CreatedAt  time.Time
}

type Usuario struct {
	ID          string
	ClerkUserID string
	EstudioID   string
	Nombre      string
	Email       string
	Rol         string
	CreatedAt   time.Time
}

type Banco struct {
	ID        string
	EstudioID string
	Nombre    string
	CreatedAt time.Time
}

type UpsertUsuarioInput struct {
	ClerkUserID string
	EstudioID   string
	Nombre      string
	Email       string
	Rol         string
}

// ── Auth repositories ────────────────────────────────────────────────────────

type EstudioRepository interface {
	UpsertByClerkOrgID(ctx context.Context, clerkOrgID, nombre string) (*Estudio, error)
	GetByClerkOrgID(ctx context.Context, clerkOrgID string) (*Estudio, error)
}

type UsuarioRepository interface {
	UpsertByClerkUserID(ctx context.Context, in UpsertUsuarioInput) (*Usuario, error)
	GetByClerkUserID(ctx context.Context, clerkUserID string) (*Usuario, error)
	GetBancoIDs(ctx context.Context, usuarioID string) ([]string, error)
	ListByEstudio(ctx context.Context, estudioID string) ([]*Usuario, error)
	GetByEstudioAndID(ctx context.Context, estudioID, id string) (*Usuario, error)
}

type UsuarioBanco struct {
	ID     string
	Nombre string
	Email  string
	Rol    string
}

type BancoRepository interface {
	Create(ctx context.Context, estudioID, nombre string) (*Banco, error)
	List(ctx context.Context, estudioID string) ([]*Banco, error)
	GetByID(ctx context.Context, estudioID, id string) (*Banco, error)
	Update(ctx context.Context, estudioID, id, nombre string) (*Banco, error)
	Delete(ctx context.Context, estudioID, id string) error
	HasCasos(ctx context.Context, id string) (bool, error)
	ListUsuarios(ctx context.Context, bancoID string) ([]*UsuarioBanco, error)
	AsignarUsuario(ctx context.Context, bancoID, usuarioID string) error
	DesasignarUsuario(ctx context.Context, bancoID, usuarioID string) error
}

// ── Caso list / detail types ─────────────────────────────────────────────────

type CasoListItem struct {
	ID             string
	BancoID        string
	BancoNombre    string
	ClienteID      string
	ClienteRUT     string
	ClienteNombre  string
	AbogadoID      *string
	NumeroOT       *string
	Estado         estado.Estado
	FechaDJ        time.Time
	DenunciaValida bool
	CreatedAt      time.Time
}

type Cliente struct {
	ID        string
	EstudioID string
	BancoID   string
	RUT       string
	Nombre    string
	Contacto  *string
	CreatedAt time.Time
}

type Operacion struct {
	ID        string
	CasoID    string
	MedioPago string
	Relacion  string
	MontoCLP  int64
	MontoUF   *float64
	FechaOp   time.Time
	CreatedAt time.Time
}

type CasoDetalle struct {
	Caso        *caso.Caso
	Cliente     *Cliente
	Operaciones []*Operacion
}

// ── Caso repositories ────────────────────────────────────────────────────────

// CasoRepository defines persistence operations for casos.
type CasoRepository interface {
	Create(ctx context.Context, c *caso.Caso) error
	GetByID(ctx context.Context, estudioID, id string) (*caso.Caso, error)
	Update(ctx context.Context, c *caso.Caso) error
	List(ctx context.Context, estudioID string, filters CaseFilters) ([]*caso.Caso, error)
	ListRich(ctx context.Context, estudioID string, filters CaseFilters) ([]*CasoListItem, int, error)
	UpdateState(ctx context.Context, id string, newState estado.Estado) error
	GetDetalle(ctx context.Context, estudioID, id string) (*CasoDetalle, error)
}

// CaseFilters parameterizes list queries.
type CaseFilters struct {
	BancoIDs  []string
	Estado    *estado.Estado
	AbogadoID *string
	Limit     int
	Offset    int
}

// UpsertClienteInput holds data for create-or-update of a client.
type UpsertClienteInput struct {
	EstudioID string
	BancoID   string
	RUT       string
	Nombre    string
	Contacto  *string
}

// ClienteRepository manages client persistence.
type ClienteRepository interface {
	Upsert(ctx context.Context, in UpsertClienteInput) (*Cliente, error)
	GetByID(ctx context.Context, estudioID, id string) (*Cliente, error)
}

// NewOperacionInput holds data for creating an operation.
type NewOperacionInput struct {
	CasoID    string
	MedioPago string
	Relacion  string
	MontoCLP  int64
	MontoUF   *float64
	FechaOp   time.Time
}

// OperacionRepository manages operations persistence.
type OperacionRepository interface {
	Create(ctx context.Context, in NewOperacionInput) (*Operacion, error)
	ListByCaso(ctx context.Context, casoID string) ([]*Operacion, error)
}

// PlazoRepository manages the legal deadlines associated with a caso.
type PlazoRepository interface {
	CreateBatch(ctx context.Context, plazos []NewPlazoInput) error
	ListByCase(ctx context.Context, casoID string) ([]StoredPlazo, error)
	MarkCompleted(ctx context.Context, plazoID string, date time.Time) error
}

type NewPlazoInput struct {
	CasoID      string
	Tipo        plazo.TipoPlazo
	FechaInicio time.Time
	DiasHabiles int
	FechaLimite time.Time
}

type StoredPlazo struct {
	ID            string
	CasoID        string
	Tipo          plazo.TipoPlazo
	FechaInicio   time.Time
	DiasHabiles   int
	FechaLimite   time.Time
	Completed     bool
	FechaCumplido *time.Time
}

// ── Supporting ports ─────────────────────────────────────────────────────────

// FeriadoProvider supplies the Chilean public holiday calendar.
type FeriadoProvider interface {
	GetHolidays(ctx context.Context, from, to time.Time) ([]time.Time, error)
}

// DocumentStorage manages file storage (Vercel Blob).
type DocumentStorage interface {
	Upload(ctx context.Context, name string, content []byte, contentType string) (url string, err error)
	Delete(ctx context.Context, url string) error
}

// AuditLogger records every caso mutation. The underlying table is append-only.
type AuditLogger interface {
	Log(ctx context.Context, entry AuditEntry) error
}

type AuditEntry struct {
	EstudioID string
	UsuarioID *string
	CasoID    *string
	Accion    string
	Detalle   map[string]any
}
