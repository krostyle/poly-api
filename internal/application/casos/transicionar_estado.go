package casos

import (
	"context"
	"errors"
	"time"

	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/estado"
	"poly.app/api/internal/domain/plazo"
)

type TransitionStateInput struct {
	EstudioID         string
	CasoID            string
	UsuarioID         string
	NewState          estado.Estado
	TerminationReason *string
	// Forzar bypasses the state machine for administrative corrections.
	// The audit entry records forzado:true so it is traceable.
	Forzar bool
}

type TransitionStateUseCase struct {
	casos    domain.CasoRepository
	plazos   domain.PlazoRepository
	feriados domain.FeriadoProvider
	auditor  domain.AuditLogger
}

func NewTransitionStateUseCase(
	casos domain.CasoRepository,
	plazos domain.PlazoRepository,
	feriados domain.FeriadoProvider,
	auditor domain.AuditLogger,
) *TransitionStateUseCase {
	return &TransitionStateUseCase{casos: casos, plazos: plazos, feriados: feriados, auditor: auditor}
}

func (uc *TransitionStateUseCase) Execute(ctx context.Context, in TransitionStateInput) error {
	c, err := uc.casos.GetByID(ctx, in.EstudioID, in.CasoID)
	if err != nil {
		return err
	}
	if !in.Forzar {
		if err := c.ValidateTransition(in.NewState); err != nil {
			return err
		}
	}
	if in.NewState == estado.Terminado && in.TerminationReason == nil {
		return errors.New("termination reason is required when closing a caso as TERMINADO")
	}

	previousState := c.Estado
	if err := uc.casos.UpdateState(ctx, in.CasoID, in.NewState); err != nil {
		return err
	}

	uc.createTransitionPlazos(ctx, in.CasoID, in.NewState)

	uid := in.UsuarioID
	detalle := map[string]any{
		"anterior": string(previousState),
		"nuevo":    string(in.NewState),
	}
	if in.Forzar {
		detalle["forzado"] = true
	}
	_ = uc.auditor.Log(ctx, domain.AuditEntry{
		EstudioID: in.EstudioID,
		UsuarioID: &uid,
		CasoID:    &in.CasoID,
		Accion:    "ESTADO_CAMBIADO",
		Detalle:   detalle,
	})
	return nil
}

func (uc *TransitionStateUseCase) createTransitionPlazos(ctx context.Context, casoID string, newState estado.Estado) {
	now := time.Now()
	horizon := now.AddDate(0, 3, 0)
	holidays, _ := uc.feriados.GetHolidays(ctx, now, horizon)

	type spec struct {
		tipo plazo.TipoPlazo
		dias int
	}

	var specs []spec
	switch newState {
	case estado.Suspension:
		specs = []spec{{plazo.TipoPrecautelar, 13}}
	case estado.Judicializacion:
		specs = []spec{{plazo.TipoDemanda, 10}}
	case estado.Restitucion:
		specs = []spec{
			{plazo.TipoRestitucionRechazo, 3},
			{plazo.TipoDemanda, 10},
		}
	default:
		return
	}

	var inputs []domain.NewPlazoInput
	for _, s := range specs {
		inputs = append(inputs, domain.NewPlazoInput{
			CasoID:      casoID,
			Tipo:        s.tipo,
			FechaInicio: now,
			DiasHabiles: s.dias,
			FechaLimite: plazo.CalculateDeadline(now, s.dias, holidays),
		})
	}
	_ = uc.plazos.CreateBatch(ctx, inputs)
}
