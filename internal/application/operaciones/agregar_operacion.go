package operaciones

import (
	"context"
	"errors"
	"time"

	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/plazo"
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
	casos       domain.CasoRepository
	operaciones domain.OperacionRepository
	plazos      domain.PlazoRepository
	feriados    domain.FeriadoProvider
	auditor     domain.AuditLogger
}

func NewAgregarOperacionUseCase(
	casos domain.CasoRepository,
	ops domain.OperacionRepository,
	plazos domain.PlazoRepository,
	feriados domain.FeriadoProvider,
	auditor domain.AuditLogger,
) *AgregarOperacionUseCase {
	return &AgregarOperacionUseCase{casos: casos, operaciones: ops, plazos: plazos, feriados: feriados, auditor: auditor}
}

func (uc *AgregarOperacionUseCase) Execute(ctx context.Context, in AgregarOperacionInput) (*domain.Operacion, error) {
	if in.MontoCLP <= 0 {
		return nil, errors.New("monto_clp must be greater than zero")
	}

	c, err := uc.casos.GetByID(ctx, in.EstudioID, in.CasoID)
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

	// Ley 20.009 Art. 5: retiros en cajero automático tienen plazo de 15 días hábiles
	// (vs. 5 días para otras operaciones). Si se agrega una operación de cajero,
	// extendemos el plazo de RESTITUCION si aún no fue cumplido y sigue en 5 días.
	if in.MedioPago == "CAJERO" && c.FechaDJ != nil {
		uc.ajustarRestitucionCajero(ctx, in.CasoID, *c.FechaDJ)
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

func (uc *AgregarOperacionUseCase) ajustarRestitucionCajero(ctx context.Context, casoID string, fechaDJ time.Time) {
	plazos, err := uc.plazos.ListByCase(ctx, casoID)
	if err != nil {
		return
	}
	for _, p := range plazos {
		if p.Tipo == plazo.TipoRestitucion && !p.Completed && p.DiasHabiles < 15 {
			horizon := fechaDJ.AddDate(0, 1, 0)
			holidays, _ := uc.feriados.GetHolidays(ctx, fechaDJ, horizon)
			nuevaFechaLimite := plazo.CalculateDeadline(fechaDJ, 15, holidays)
			_ = uc.plazos.UpdateDiasHabiles(ctx, p.ID, 15, nuevaFechaLimite)
			return
		}
	}
}
