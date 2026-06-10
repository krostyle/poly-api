package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"poly.app/api/internal/adapters/http/middleware"
	appops "poly.app/api/internal/application/operaciones"
	"poly.app/api/internal/domain"
)

type OperacionesHandler struct {
	agregar *appops.AgregarOperacionUseCase
	repo    domain.OperacionRepository
}

func NewOperacionesHandler(agregar *appops.AgregarOperacionUseCase, repo domain.OperacionRepository) *OperacionesHandler {
	return &OperacionesHandler{agregar: agregar, repo: repo}
}

func (h *OperacionesHandler) Crear(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	usuarioID := middleware.UsuarioIDFromCtx(r.Context())
	casoID := chi.URLParam(r, "id")

	var req struct {
		MedioPago string   `json:"medio_pago"`
		Relacion  string   `json:"relacion"`
		MontoCLP  int64    `json:"monto_clp"`
		MontoUF   *float64 `json:"monto_uf"`
		FechaOp   string   `json:"fecha_op"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.MedioPago == "" || req.Relacion == "" || req.FechaOp == "" {
		http.Error(w, `{"error":"medio_pago, relacion and fecha_op are required"}`, http.StatusBadRequest)
		return
	}

	fechaOp, err := time.Parse("2006-01-02", req.FechaOp)
	if err != nil {
		http.Error(w, `{"error":"fecha_op must be YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	op, err := h.agregar.Execute(r.Context(), appops.AgregarOperacionInput{
		EstudioID: estudioID,
		CasoID:    casoID,
		UsuarioID: usuarioID,
		MedioPago: req.MedioPago,
		Relacion:  req.Relacion,
		MontoCLP:  req.MontoCLP,
		MontoUF:   req.MontoUF,
		FechaOp:   fechaOp,
	})
	if err != nil {
		if err.Error() == "monto_clp must be greater than zero" {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"error":"could not create operacion"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(operacionJSON{
		ID:        op.ID,
		CasoID:    op.CasoID,
		MedioPago: op.MedioPago,
		Relacion:  op.Relacion,
		MontoCLP:  op.MontoCLP,
		MontoUF:   op.MontoUF,
		FechaOp:   op.FechaOp.Format("2006-01-02"),
	})
}

func (h *OperacionesHandler) Listar(w http.ResponseWriter, r *http.Request) {
	casoID := chi.URLParam(r, "id")

	ops, err := h.repo.ListByCaso(r.Context(), casoID)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	resp := make([]operacionJSON, 0, len(ops))
	for _, op := range ops {
		resp = append(resp, operacionJSON{
			ID:        op.ID,
			CasoID:    op.CasoID,
			MedioPago: op.MedioPago,
			Relacion:  op.Relacion,
			MontoCLP:  op.MontoCLP,
			MontoUF:   op.MontoUF,
			FechaOp:   op.FechaOp.Format("2006-01-02"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"operaciones": resp})
}
