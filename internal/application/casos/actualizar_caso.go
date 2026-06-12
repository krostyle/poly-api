package casos

import (
	"context"
	"time"

	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/caso"
	"poly.app/api/internal/domain/plazo"
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
	ClearFechaDJ       bool
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

	oldFechaDJ := c.FechaDJ

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
	if in.ClearFechaDJ {
		c.FechaDJ = nil
	} else if in.FechaDJ != nil {
		c.FechaDJ = in.FechaDJ
	}

	if err := uc.casos.Update(ctx, c); err != nil {
		return nil, err
	}

	// Create the 30-business-day response deadline the first time FechaDJ is recorded.
	if !in.ClearFechaDJ && in.FechaDJ != nil && oldFechaDJ == nil {
		uc.createRespuestaDenunciaPlazos(ctx, c.ID, *in.FechaDJ)
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

func (uc *UpdateCaseUseCase) createRespuestaDenunciaPlazos(ctx context.Context, casoID string, fechaDJ time.Time) {
	horizon := fechaDJ.AddDate(0, 3, 0)
	holidays, _ := uc.feriados.GetHolidays(ctx, fechaDJ, horizon)
	_ = uc.plazos.CreateBatch(ctx, []domain.NewPlazoInput{{
		CasoID:      casoID,
		Tipo:        plazo.TipoRespuestaDenuncia,
		FechaInicio: fechaDJ,
		DiasHabiles: 30,
		FechaLimite: plazo.CalculateDeadline(fechaDJ, 30, holidays),
	}})
}
