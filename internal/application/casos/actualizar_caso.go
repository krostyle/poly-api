package casos

import (
	"context"
	"time"

	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/caso"
)

type UpdateCaseInput struct {
	EstudioID          string
	CasoID             string
	UsuarioID          string
	AbogadoID          *string
	NumeroOT           *string
	EstadoDenuncia     *caso.EstadoDenuncia
	FechaDenuncia      *time.Time
	FechaDJ            *time.Time
	ClearFechaDenuncia bool
}

type UpdateCaseUseCase struct {
	casos    domain.CasoRepository
	plazos   domain.PlazoRepository
	feriados domain.FeriadoProvider
	auditor  domain.AuditLogger
}

func NewUpdateCaseUseCase(
	casos domain.CasoRepository,
	plazos domain.PlazoRepository,
	feriados domain.FeriadoProvider,
	auditor domain.AuditLogger,
) *UpdateCaseUseCase {
	return &UpdateCaseUseCase{casos: casos, plazos: plazos, feriados: feriados, auditor: auditor}
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
	if in.EstadoDenuncia != nil {
		c.EstadoDenuncia = *in.EstadoDenuncia
	}
	if in.ClearFechaDenuncia {
		c.FechaDenuncia = nil
	} else if in.FechaDenuncia != nil {
		c.FechaDenuncia = in.FechaDenuncia
	}
	if in.FechaDJ != nil {
		c.FechaDJ = in.FechaDJ
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

