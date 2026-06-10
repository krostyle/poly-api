package plazos

import (
	"context"
	"time"

	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/plazo"
)

type RecalcDeadlinesUseCase struct {
	plazoRepo domain.PlazoRepository
	holidays  domain.FeriadoProvider
}

func NewRecalcDeadlinesUseCase(plazoRepo domain.PlazoRepository, holidays domain.FeriadoProvider) *RecalcDeadlinesUseCase {
	return &RecalcDeadlinesUseCase{plazoRepo: plazoRepo, holidays: holidays}
}

type PlazoEvaluation struct {
	PlazoID       string
	Tipo          plazo.TipoPlazo
	FechaLimite   time.Time
	RemainingDays int
	Semaforo      plazo.Semaforo
}

func (uc *RecalcDeadlinesUseCase) EvaluateCase(ctx context.Context, casoID string) ([]PlazoEvaluation, error) {
	stored, err := uc.plazoRepo.ListByCase(ctx, casoID)
	if err != nil {
		return nil, err
	}

	today := time.Now()
	var from, to time.Time
	for _, p := range stored {
		if !p.Completed {
			if from.IsZero() || p.FechaInicio.Before(from) {
				from = p.FechaInicio
			}
			if to.IsZero() || p.FechaLimite.After(to) {
				to = p.FechaLimite.AddDate(0, 0, 30)
			}
		}
	}

	var holidayList []time.Time
	if !from.IsZero() {
		holidayList, _ = uc.holidays.GetHolidays(ctx, from, to)
	}

	var results []PlazoEvaluation
	for _, p := range stored {
		if p.Completed {
			continue
		}
		remaining := plazo.RemainingBusinessDays(today, p.FechaLimite, holidayList)
		results = append(results, PlazoEvaluation{
			PlazoID:       p.ID,
			Tipo:          p.Tipo,
			FechaLimite:   p.FechaLimite,
			RemainingDays: remaining,
			Semaforo:      plazo.EvaluateSemaforo(remaining, plazo.DefaultThresholds),
		})
	}
	return results, nil
}
