package casos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/caso"
	"poly.app/api/internal/domain/estado"
)

type CreateCaseInput struct {
	EstudioID string
	BancoID   string
	ClienteID string
	FechaDJ   time.Time
}

type CreateCaseUseCase struct {
	casos   domain.CasoRepository
	auditor domain.AuditLogger
}

func NewCreateCaseUseCase(casos domain.CasoRepository, auditor domain.AuditLogger) *CreateCaseUseCase {
	return &CreateCaseUseCase{casos: casos, auditor: auditor}
}

func (uc *CreateCaseUseCase) Execute(ctx context.Context, in CreateCaseInput) (*caso.Caso, error) {
	c := &caso.Caso{
		ID:        uuid.New().String(),
		EstudioID: in.EstudioID,
		BancoID:   in.BancoID,
		ClienteID: in.ClienteID,
		Estado:    estado.Llamada,
		FechaDJ:   in.FechaDJ,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := uc.casos.Create(ctx, c); err != nil {
		return nil, err
	}
	_ = uc.auditor.Log(ctx, domain.AuditEntry{
		EstudioID: in.EstudioID,
		CasoID:    &c.ID,
		Accion:    "CASO_CREADO",
		Detalle:   map[string]any{"estado": string(c.Estado)},
	})
	return c, nil
}
