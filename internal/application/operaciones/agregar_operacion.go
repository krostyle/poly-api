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

	// Recalculate RESTITUCION after every operation: transaction type and total UF
	// both affect the deadline under Ley 20.009 Art. 5.
	if c.FechaDJ != nil {
		uc.ajustarRestitucion(ctx, in.CasoID, *c.FechaDJ)
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

// ajustarRestitucion recalculates the RESTITUCION plazo according to Ley 20.009 Art. 5:
//   - Base: 5 días hábiles (cards/transfers) or 15 días (ATM/cajero)
//   - +7 días adicionales if total disputed amount exceeds 35 UF
//
// Only extends — never shortens — the deadline.
func (uc *AgregarOperacionUseCase) ajustarRestitucion(ctx context.Context, casoID string, fechaDJ time.Time) {
	ops, err := uc.operaciones.ListByCaso(ctx, casoID)
	if err != nil {
		return
	}

	base := 5
	var totalUF float64
	for _, op := range ops {
		if op.MedioPago == "CAJERO" {
			base = 15
		}
		if op.MontoUF != nil {
			totalUF += *op.MontoUF
		}
	}

	targetDias := base
	if totalUF > 35 {
		targetDias = base + 7
	}

	plazos, err := uc.plazos.ListByCase(ctx, casoID)
	if err != nil {
		return
	}
	for _, p := range plazos {
		if p.Tipo == plazo.TipoRestitucion && !p.Completed && p.DiasHabiles < targetDias {
			horizon := fechaDJ.AddDate(0, 1, 0)
			holidays, _ := uc.feriados.GetHolidays(ctx, fechaDJ, horizon)
			nuevaFechaLimite := plazo.CalculateDeadline(fechaDJ, targetDias, holidays)
			_ = uc.plazos.UpdateDiasHabiles(ctx, p.ID, targetDias, nuevaFechaLimite)
			return
		}
	}
}
