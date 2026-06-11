package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	appdash "poly.app/api/internal/application/dashboard"
	"poly.app/api/internal/adapters/http/middleware"
)

type DashboardHandler struct {
	uc *appdash.DashboardUseCase
}

func NewDashboardHandler(uc *appdash.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{uc: uc}
}

type porVencerResponse struct {
	CasoID        string            `json:"caso_id"`
	BancoID       string            `json:"banco_id"`
	BancoNombre   string            `json:"banco_nombre"`
	NumeroOT      *string           `json:"numero_ot,omitempty"`
	ClienteRUT    string            `json:"cliente_rut"`
	ClienteNombre string            `json:"cliente_nombre"`
	Estado        string            `json:"estado"`
	PlazoCritico  plazoCriticoResp  `json:"plazo_critico"`
}

type plazoCriticoResp struct {
	ID            string `json:"id"`
	Tipo          string `json:"tipo"`
	FechaLimite   string `json:"fecha_limite"`
	DiasRestantes int    `json:"dias_restantes"`
	Semaforo      string `json:"semaforo"`
}

func (h *DashboardHandler) PorVencer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	estudioID := middleware.EstudioIDFromCtx(ctx)
	bancoIDs := middleware.BancoIDsFromCtx(ctx)

	dias := 7
	if q := r.URL.Query().Get("dias"); q != "" {
		if v, err := strconv.Atoi(q); err == nil && v > 0 {
			dias = v
		}
	}
	if bancoID := r.URL.Query().Get("bancoId"); bancoID != "" {
		bancoIDs = filterBancoIDs(bancoIDs, bancoID)
	}

	items, err := h.uc.PorVencer(ctx, estudioID, bancoIDs, dias)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]porVencerResponse, 0, len(items))
	for _, item := range items {
		result = append(result, porVencerResponse{
			CasoID:        item.CasoID,
			BancoID:       item.BancoID,
			BancoNombre:   item.BancoNombre,
			NumeroOT:      item.NumeroOT,
			ClienteRUT:    item.ClienteRUT,
			ClienteNombre: item.ClienteNombre,
			Estado:        item.Estado,
			PlazoCritico: plazoCriticoResp{
				ID:            item.PlazoCritico.ID,
				Tipo:          item.PlazoCritico.Tipo,
				FechaLimite:   item.PlazoCritico.FechaLimite.Format("2006-01-02"),
				DiasRestantes: item.PlazoCritico.DiasRestantes,
				Semaforo:      item.PlazoCritico.Semaforo,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"items": result})
}

type casoNuevoResponse struct {
	CasoID        string  `json:"caso_id"`
	BancoID       string  `json:"banco_id"`
	BancoNombre   string  `json:"banco_nombre"`
	ClienteRUT    string  `json:"cliente_rut"`
	ClienteNombre string  `json:"cliente_nombre"`
	AbogadoID     *string `json:"abogado_id,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

func (h *DashboardHandler) Nuevos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	estudioID := middleware.EstudioIDFromCtx(ctx)
	bancoIDs := middleware.BancoIDsFromCtx(ctx)

	if bancoID := r.URL.Query().Get("bancoId"); bancoID != "" {
		bancoIDs = filterBancoIDs(bancoIDs, bancoID)
	}

	items, err := h.uc.Nuevos(ctx, estudioID, bancoIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]casoNuevoResponse, 0, len(items))
	for _, item := range items {
		result = append(result, casoNuevoResponse{
			CasoID:        item.CasoID,
			BancoID:       item.BancoID,
			BancoNombre:   item.BancoNombre,
			ClienteRUT:    item.ClienteRUT,
			ClienteNombre: item.ClienteNombre,
			AbogadoID:     item.AbogadoID,
			CreatedAt:     item.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"items": result})
}

type casoEstancadoResponse struct {
	CasoID           string  `json:"caso_id"`
	BancoID          string  `json:"banco_id"`
	BancoNombre      string  `json:"banco_nombre"`
	NumeroOT         *string `json:"numero_ot,omitempty"`
	ClienteRUT       string  `json:"cliente_rut"`
	ClienteNombre    string  `json:"cliente_nombre"`
	Estado           string  `json:"estado"`
	UltimoMovimiento string  `json:"ultimo_movimiento"`
	DiasEstancado    int     `json:"dias_estancado"`
}

func (h *DashboardHandler) Estancados(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	estudioID := middleware.EstudioIDFromCtx(ctx)
	bancoIDs := middleware.BancoIDsFromCtx(ctx)

	diasSinMovimiento := 5
	if q := r.URL.Query().Get("dias"); q != "" {
		if v, err := strconv.Atoi(q); err == nil && v > 0 {
			diasSinMovimiento = v
		}
	}
	if bancoID := r.URL.Query().Get("bancoId"); bancoID != "" {
		bancoIDs = filterBancoIDs(bancoIDs, bancoID)
	}

	items, err := h.uc.Estancados(ctx, estudioID, bancoIDs, diasSinMovimiento)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]casoEstancadoResponse, 0, len(items))
	for _, item := range items {
		result = append(result, casoEstancadoResponse{
			CasoID:           item.CasoID,
			BancoID:          item.BancoID,
			BancoNombre:      item.BancoNombre,
			NumeroOT:         item.NumeroOT,
			ClienteRUT:       item.ClienteRUT,
			ClienteNombre:    item.ClienteNombre,
			Estado:           item.Estado,
			UltimoMovimiento: item.UltimoMovimiento.Format("2006-01-02T15:04:05Z"),
			DiasEstancado:    item.DiasEstancado,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"items": result})
}

type cargaAbogadoResponse struct {
	AbogadoID string `json:"abogado_id"`
	Nombre    string `json:"nombre"`
	Total     int    `json:"total"`
	PorVencer int    `json:"por_vencer"`
	Vencidos  int    `json:"vencidos"`
}

func (h *DashboardHandler) PorAbogado(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	estudioID := middleware.EstudioIDFromCtx(ctx)
	bancoIDs := middleware.BancoIDsFromCtx(ctx)

	if bancoID := r.URL.Query().Get("bancoId"); bancoID != "" {
		bancoIDs = filterBancoIDs(bancoIDs, bancoID)
	}

	items, err := h.uc.PorAbogado(ctx, estudioID, bancoIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]cargaAbogadoResponse, 0, len(items))
	for _, item := range items {
		result = append(result, cargaAbogadoResponse{
			AbogadoID: item.AbogadoID,
			Nombre:    item.Nombre,
			Total:     item.Total,
			PorVencer: item.PorVencer,
			Vencidos:  item.Vencidos,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"abogados": result})
}

// filterBancoIDs returns [bancoID] only if it exists in allowed; otherwise returns nil.
func filterBancoIDs(allowed []string, bancoID string) []string {
	for _, id := range allowed {
		if id == bancoID {
			return []string{bancoID}
		}
	}
	return nil
}
