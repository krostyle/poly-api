package casos

import (
	"context"
	"errors"
	"time"

	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/caso"
	"poly.app/api/internal/domain/estado"
	"poly.app/api/internal/domain/plazo"
)

var (
	ErrMotivoTerminoRequerido     = errors.New("se requiere seleccionar un motivo de término")
	ErrMotivoTerminoInvalido      = errors.New("el motivo de término seleccionado no es válido")
	ErrDenunciaRechazadaRequerida = errors.New("el banco debe haber rechazado la denuncia para ingresar a Pago Normativo")
	ErrDenunciaAcogidaRequerida   = errors.New("el banco debe haber acogido la denuncia para pasar directamente a la etapa Judicial")
)

type TransitionStateInput struct {
	EstudioID         string
	CasoID            string
	UsuarioID         string
	NewState          estado.Estado
	TerminationReason *string
	// Forzar bypasses the state machine for administrative corrections.
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
		if err := validateDenunciaGuard(c, in.NewState); err != nil {
			return err
		}
	}
	if in.NewState == estado.Terminado {
		if in.TerminationReason == nil || *in.TerminationReason == "" {
			return ErrMotivoTerminoRequerido
		}
		if !caso.IsValidMotivoTermino(*in.TerminationReason) {
			return ErrMotivoTerminoInvalido
		}
	}

	previousState := c.Estado
	if err := uc.casos.UpdateState(ctx, in.CasoID, in.NewState); err != nil {
		return err
	}

	uc.createTransitionPlazos(ctx, in.CasoID, in.NewState)

	// Cuando el caso entra a JUDICIAL, el proceso administrativo de restitución
	// de Ley 20.009 queda superado por la vía judicial. Marcamos el plazo como
	// cumplido para que no genere alertas innecesarias.
	if in.NewState == estado.Judicial {
		uc.marcarRestitucionCumplida(ctx, in.CasoID)
	}

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

// validateDenunciaGuard enforces the business rule that the bank's response to the
// denuncia determines which path is taken from Prejudicial:
//   - Rechazada → PagoNormativo (bank rejects, normative payment phase applies)
//   - Acogida   → Judicial directly (bank accepts, normative payment is skipped)
func validateDenunciaGuard(c *caso.Caso, target estado.Estado) error {
	switch target {
	case estado.PagoNormativo:
		if c.EstadoDenuncia != caso.DenunciaRechazada {
			return ErrDenunciaRechazadaRequerida
		}
	case estado.Judicial:
		// Guard only applies when coming from Prejudicial (skipping PagoNormativo).
		// From PagoNormativo → Judicial no additional check is needed.
		if c.Estado == estado.Prejudicial && c.EstadoDenuncia != caso.DenunciaAcogida {
			return ErrDenunciaAcogidaRequerida
		}
	}
	return nil
}

func (uc *TransitionStateUseCase) marcarRestitucionCumplida(ctx context.Context, casoID string) {
	plazos, err := uc.plazos.ListByCase(ctx, casoID)
	if err != nil {
		return
	}
	now := time.Now()
	for _, p := range plazos {
		if p.Tipo == plazo.TipoRestitucion && !p.Completed {
			_ = uc.plazos.MarkCompleted(ctx, p.ID, now)
			return
		}
	}
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
	case estado.Prejudicial:
		specs = []spec{
			{plazo.TipoPrecautelar, 13},
			{plazo.TipoResolucionJPL, 3},
		}
	case estado.PagoNormativo:
		specs = []spec{
			{plazo.TipoDemanda, 10},
			{plazo.TipoRestitucionRechazo, 3},
		}
	case estado.Judicial:
		specs = []spec{{plazo.TipoDemanda, 10}}
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
