package casos

import (
	"context"
	"time"

	"poly.app/api/internal/domain"
)

type UpdateCaseInput struct {
	EstudioID      string
	CasoID         string
	UsuarioID      string
	AbogadoID      *string
	NumeroOT       *string
	DenunciaValida *bool
	FechaDenuncia  *time.Time
}

type UpdateCaseUseCase struct {
	casos   domain.CasoRepository
	auditor domain.AuditLogger
}

func NewUpdateCaseUseCase(casos domain.CasoRepository, auditor domain.AuditLogger) *UpdateCaseUseCase {
	return &UpdateCaseUseCase{casos: casos, auditor: auditor}
}

func (uc *UpdateCaseUseCase) Execute(ctx context.Context, in UpdateCaseInput) (*domain.CasoDetalle, error) {
	c, err := uc.casos.GetByID(ctx, in.EstudioID, in.CasoID)
	if err != nil {
		return nil, err
	}

	if in.AbogadoID != nil {
		c.AbogadoID = in.AbogadoID
	}
	if in.NumeroOT != nil {
		c.NumeroOT = in.NumeroOT
	}
	if in.DenunciaValida != nil {
		c.DenunciaValida = *in.DenunciaValida
	}
	if in.FechaDenuncia != nil {
		c.FechaDenuncia = in.FechaDenuncia
	}

	if err := uc.casos.Update(ctx, c); err != nil {
		return nil, err
	}

	uid := in.UsuarioID
	_ = uc.auditor.Log(ctx, domain.AuditEntry{
		EstudioID: in.EstudioID,
		UsuarioID: &uid,
		CasoID:    &in.CasoID,
		Accion:    "CASO_ACTUALIZADO",
		Detalle:   map[string]any{"caso_id": in.CasoID},
	})

	return uc.casos.GetDetalle(ctx, in.EstudioID, in.CasoID)
}
