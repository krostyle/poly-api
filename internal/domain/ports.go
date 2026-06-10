package domain

import (
	"context"
	"time"

	"poly.app/api/internal/domain/caso"
	"poly.app/api/internal/domain/estado"
	"poly.app/api/internal/domain/plazo"
)

// CasoRepository defines persistence operations for casos.
// The domain declares this interface; adapters implement it.
type CasoRepository interface {
	Create(ctx context.Context, c *caso.Caso) error
	GetByID(ctx context.Context, estudioID, id string) (*caso.Caso, error)
	Update(ctx context.Context, c *caso.Caso) error
	List(ctx context.Context, estudioID string, filters CaseFilters) ([]*caso.Caso, error)
	UpdateState(ctx context.Context, id string, newState estado.Estado) error
}

// CaseFilters parameterizes list queries.
type CaseFilters struct {
	BancoIDs  []string
	Estado    *estado.Estado
	AbogadoID *string
	Limit     int
	Offset    int
}

// PlazoRepository manages the legal deadlines associated with a caso.
type PlazoRepository interface {
	CreateBatch(ctx context.Context, plazos []NewPlazoInput) error
	ListByCase(ctx context.Context, casoID string) ([]StoredPlazo, error)
	MarkCompleted(ctx context.Context, plazoID string, date time.Time) error
}

// NewPlazoInput is the creation DTO for a plazo.
type NewPlazoInput struct {
	CasoID      string
	Tipo        plazo.TipoPlazo
	FechaInicio time.Time
	DiasHabiles int
	FechaLimite time.Time
}

// StoredPlazo is the plazo as it exists in the database.
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

// AuditEntry represents one line in the audit log.
type AuditEntry struct {
	EstudioID string
	UsuarioID *string
	CasoID    *string
	Accion    string
	Detalle   map[string]any
}
