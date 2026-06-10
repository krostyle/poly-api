package casos

import (
	"context"

	"poly.app/api/internal/domain"
)

type AssignLawyerInput struct {
	EstudioID string
	CasoID    string
	AbogadoID string
	NumeroOT  string
}

type AssignLawyerUseCase struct {
	casos   domain.CasoRepository
	auditor domain.AuditLogger
}

func NewAssignLawyerUseCase(casos domain.CasoRepository, auditor domain.AuditLogger) *AssignLawyerUseCase {
	return &AssignLawyerUseCase{casos: casos, auditor: auditor}
}

func (uc *AssignLawyerUseCase) Execute(ctx context.Context, in AssignLawyerInput) error {
	c, err := uc.casos.GetByID(ctx, in.EstudioID, in.CasoID)
	if err != nil {
		return err
	}
	c.AbogadoID = &in.AbogadoID
	c.NumeroOT = &in.NumeroOT
	if err := uc.casos.Update(ctx, c); err != nil {
		return err
	}
	_ = uc.auditor.Log(ctx, domain.AuditEntry{
		EstudioID: in.EstudioID,
		CasoID:    &in.CasoID,
		Accion:    "ABOGADO_ASIGNADO",
		Detalle:   map[string]any{"abogado_id": in.AbogadoID, "numero_ot": in.NumeroOT},
	})
	return nil
}
