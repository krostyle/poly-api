package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"poly.app/api/internal/adapters/http/middleware"
	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/plazo"
)

// configurablePlazos lists the deadline types that admins can configure per estudio.
var configurablePlazos = []struct {
	Tipo       plazo.TipoPlazo
	DefaultDias int
	Label      string
}{
	{plazo.TipoAnalisisInterno, 5, "Análisis interno"},
	{plazo.TipoAsignacion, 7, "Asignación de abogado"},
}

type ConfiguracionHandler struct {
	repo domain.ConfiguracionPlazoRepository
}

func NewConfiguracionHandler(repo domain.ConfiguracionPlazoRepository) *ConfiguracionHandler {
	return &ConfiguracionHandler{repo: repo}
}

type configuracionPlazoJSON struct {
	Tipo        string `json:"tipo"`
	Label       string `json:"label"`
	DiasHabiles int    `json:"dias_habiles"`
	EsDefault   bool   `json:"es_default"`
}

func (h *ConfiguracionHandler) Listar(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	configs, err := h.repo.GetByEstudio(r.Context(), estudioID)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	configMap := make(map[plazo.TipoPlazo]int)
	for _, c := range configs {
		configMap[c.TipoPlazo] = c.DiasHabiles
	}

	result := make([]configuracionPlazoJSON, 0, len(configurablePlazos))
	for _, cp := range configurablePlazos {
		dias, ok := configMap[cp.Tipo]
		if !ok {
			dias = cp.DefaultDias
		}
		result = append(result, configuracionPlazoJSON{
			Tipo:        string(cp.Tipo),
			Label:       cp.Label,
			DiasHabiles: dias,
			EsDefault:   !ok,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *ConfiguracionHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	tipo := chi.URLParam(r, "tipo")

	var isConfigurable bool
	for _, cp := range configurablePlazos {
		if string(cp.Tipo) == tipo {
			isConfigurable = true
			break
		}
	}
	if !isConfigurable {
		http.Error(w, `{"error":"tipo de plazo no configurable"}`, http.StatusBadRequest)
		return
	}

	var body struct {
		DiasHabiles int `json:"dias_habiles"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if body.DiasHabiles < 1 || body.DiasHabiles > 90 {
		http.Error(w, `{"error":"dias_habiles debe estar entre 1 y 90"}`, http.StatusBadRequest)
		return
	}

	if err := h.repo.Upsert(r.Context(), estudioID, plazo.TipoPlazo(tipo), body.DiasHabiles); err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
