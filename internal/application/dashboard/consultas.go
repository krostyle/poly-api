package dashboard

import (
	"context"

	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/caso"
)

type DashboardUseCase struct {
	casos domain.CasoRepository
}

func NewDashboardUseCase(casos domain.CasoRepository) *DashboardUseCase {
	return &DashboardUseCase{casos: casos}
}

// UpcomingDeadlines returns casos with plazos about to expire (delegates to an optimized query).
func (uc *DashboardUseCase) UpcomingDeadlines(ctx context.Context, estudioID string, bancoIDs []string) ([]*caso.Caso, error) {
	return uc.casos.List(ctx, estudioID, domain.CaseFilters{
		BancoIDs: bancoIDs,
		Limit:    50,
	})
}
