package casos

import (
	"context"
	"errors"

	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/estado"
)

type TransitionStateInput struct {
	EstudioID         string
	CasoID            string
	UsuarioID         string
	NewState          estado.Estado
	TerminationReason *string
}

type TransitionStateUseCase struct {
	casos   domain.CasoRepository
	auditor domain.AuditLogger
}

func NewTransitionStateUseCase(casos domain.CasoRepository, auditor domain.AuditLogger) *TransitionStateUseCase {
	return &TransitionStateUseCase{casos: casos, auditor: auditor}
}

func (uc *TransitionStateUseCase) Execute(ctx context.Context, in TransitionStateInput) error {
	c, err := uc.casos.GetByID(ctx, in.EstudioID, in.CasoID)
	if err != nil {
		return err
	}
	if err := c.ValidateTransition(in.NewState); err != nil {
		return err
	}
	if in.NewState == estado.Terminado && in.TerminationReason == nil {
		return errors.New("termination reason is required when closing a caso as TERMINADO")
	}

	previousState := c.Estado
	if err := uc.casos.UpdateState(ctx, in.CasoID, in.NewState); err != nil {
		return err
	}

	uid := in.UsuarioID
	_ = uc.auditor.Log(ctx, domain.AuditEntry{
		EstudioID: in.EstudioID,
		UsuarioID: &uid,
		CasoID:    &in.CasoID,
		Accion:    "ESTADO_CAMBIADO",
		Detalle: map[string]any{
			"anterior": string(previousState),
			"nuevo":    string(in.NewState),
		},
	})
	return nil
}
