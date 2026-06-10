package operaciones

import (
	"context"
	"errors"
	"time"

	"poly.app/api/internal/domain"
)

type AgregarOperacionInput struct {
	EstudioID string
	CasoID    string
	UsuarioID string
	MedioPago string
	Relacion  string
	MontoCLP  int64
	MontoUF   *float64
	FechaOp   time.Time
}

type AgregarOperacionUseCase struct {
	casos      domain.CasoRepository
	operaciones domain.OperacionRepository
	auditor    domain.AuditLogger
}

func NewAgregarOperacionUseCase(
	casos domain.CasoRepository,
	ops domain.OperacionRepository,
	auditor domain.AuditLogger,
) *AgregarOperacionUseCase {
	return &AgregarOperacionUseCase{casos: casos, operaciones: ops, auditor: auditor}
}

func (uc *AgregarOperacionUseCase) Execute(ctx context.Context, in AgregarOperacionInput) (*domain.Operacion, error) {
	if in.MontoCLP <= 0 {
		return nil, errors.New("monto_clp must be greater than zero")
	}

	_, err := uc.casos.GetByID(ctx, in.EstudioID, in.CasoID)
	if err != nil {
		return nil, err
	}

	op, err := uc.operaciones.Create(ctx, domain.NewOperacionInput{
		CasoID:    in.CasoID,
		MedioPago: in.MedioPago,
		Relacion:  in.Relacion,
		MontoCLP:  in.MontoCLP,
		MontoUF:   in.MontoUF,
		FechaOp:   in.FechaOp,
	})
	if err != nil {
		return nil, err
	}

	uid := in.UsuarioID
	_ = uc.auditor.Log(ctx, domain.AuditEntry{
		EstudioID: in.EstudioID,
		UsuarioID: &uid,
		CasoID:    &in.CasoID,
		Accion:    "OPERACION_AGREGADA",
		Detalle:   map[string]any{"operacion_id": op.ID, "monto_clp": in.MontoCLP},
	})

	return op, nil
}
