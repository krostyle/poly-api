package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	appplazos "poly.app/api/internal/application/plazos"
	"poly.app/api/internal/adapters/http/middleware"
	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/plazo"
)

type PlazosHandler struct {
	repo     domain.PlazoRepository
	feriados domain.FeriadoProvider
	recalc   *appplazos.RecalcDeadlinesUseCase
}

func NewPlazosHandler(repo domain.PlazoRepository, feriados domain.FeriadoProvider) *PlazosHandler {
	return &PlazosHandler{
		repo:     repo,
		feriados: feriados,
		recalc:   appplazos.NewRecalcDeadlinesUseCase(repo, feriados),
	}
}

type plazoResponse struct {
	ID            string  `json:"id"`
	Tipo          string  `json:"tipo"`
	FechaInicio   string  `json:"fecha_inicio"`
	FechaLimite   string  `json:"fecha_limite"`
	DiasHabiles   int     `json:"dias_habiles"`
	DiasRestantes int     `json:"dias_restantes"`
	Semaforo      string  `json:"semaforo"`
	Cumplido      bool    `json:"cumplido"`
	FechaCumplido *string `json:"fecha_cumplido,omitempty"`
}

func (h *PlazosHandler) ListarPorCaso(w http.ResponseWriter, r *http.Request) {
	casoID := chi.URLParam(r, "id")
	ctx := r.Context()

	stored, err := h.repo.ListByCase(ctx, casoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	var holidays []time.Time
	if !from.IsZero() {
		holidays, _ = h.feriados.GetHolidays(ctx, from, to)
	}

	result := make([]plazoResponse, 0, len(stored))
	for _, p := range stored {
		remaining := 0
		sem := plazo.Verde
		if !p.Completed {
			remaining = plazo.RemainingBusinessDays(today, p.FechaLimite, holidays)
			sem = plazo.EvaluateSemaforo(remaining, plazo.DefaultThresholds)
		}

		resp := plazoResponse{
			ID:            p.ID,
			Tipo:          string(p.Tipo),
			FechaInicio:   p.FechaInicio.Format("2006-01-02"),
			FechaLimite:   p.FechaLimite.Format("2006-01-02"),
			DiasHabiles:   p.DiasHabiles,
			DiasRestantes: remaining,
			Semaforo:      string(sem),
			Cumplido:      p.Completed,
		}
		if p.FechaCumplido != nil {
			s := p.FechaCumplido.Format("2006-01-02")
			resp.FechaCumplido = &s
		}
		result = append(result, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"plazos": result})
}

func (h *PlazosHandler) CumplirPlazo(w http.ResponseWriter, r *http.Request) {
	plazoID := chi.URLParam(r, "plazoID")
	if err := h.repo.MarkCompleted(r.Context(), plazoID, time.Now()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /v1/plazos — all active (non-cumplido) plazos for the estudio, enriched with caso info.
func (h *PlazosHandler) ListarGlobal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	estudioID := middleware.EstudioIDFromCtx(ctx)
	bancoIDs := middleware.BancoIDsFromCtx(ctx)
	if len(bancoIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"plazos": []any{}})
		return
	}

	tipoFilter := r.URL.Query().Get("tipo")
	semaforoFilter := r.URL.Query().Get("semaforo")

	stored, err := h.repo.ListGlobal(ctx, estudioID, bancoIDs, tipoFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	today := time.Now()
	holidays, _ := h.feriados.GetHolidays(ctx, today, today.AddDate(0, 3, 0))

	type plazoGlobalResponse struct {
		ID            string  `json:"id"`
		CasoID        string  `json:"caso_id"`
		NumeroOT      *string `json:"numero_ot,omitempty"`
		ClienteNombre string  `json:"cliente_nombre"`
		ClienteRUT    string  `json:"cliente_rut"`
		BancoNombre   string  `json:"banco_nombre"`
		Estado        string  `json:"estado"`
		Tipo          string  `json:"tipo"`
		FechaInicio   string  `json:"fecha_inicio"`
		FechaLimite   string  `json:"fecha_limite"`
		DiasHabiles   int     `json:"dias_habiles"`
		DiasRestantes int     `json:"dias_restantes"`
		Semaforo      string  `json:"semaforo"`
	}

	result := make([]plazoGlobalResponse, 0, len(stored))
	for _, p := range stored {
		remaining := plazo.RemainingBusinessDays(today, p.FechaLimite, holidays)
		sem := string(plazo.EvaluateSemaforo(remaining, plazo.DefaultThresholds))

		if semaforoFilter != "" && sem != semaforoFilter {
			continue
		}

		result = append(result, plazoGlobalResponse{
			ID:            p.ID,
			CasoID:        p.CasoID,
			NumeroOT:      p.NumeroOT,
			ClienteNombre: p.ClienteNombre,
			ClienteRUT:    p.ClienteRUT,
			BancoNombre:   p.BancoNombre,
			Estado:        p.Estado,
			Tipo:          string(p.Tipo),
			FechaInicio:   p.FechaInicio.Format("2006-01-02"),
			FechaLimite:   p.FechaLimite.Format("2006-01-02"),
			DiasHabiles:   p.DiasHabiles,
			DiasRestantes: remaining,
			Semaforo:      sem,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"plazos": result})
}
