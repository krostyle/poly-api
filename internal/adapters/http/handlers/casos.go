package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"poly.app/api/internal/adapters/http/middleware"
	appcasos "poly.app/api/internal/application/casos"
	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/caso"
	"poly.app/api/internal/domain/estado"
)

type CasosHandler struct {
	crear      *appcasos.CreateCaseUseCase
	actualizar *appcasos.UpdateCaseUseCase
	transicion *appcasos.TransitionStateUseCase
	repo       domain.CasoRepository
}

func NewCasosHandler(
	crear *appcasos.CreateCaseUseCase,
	actualizar *appcasos.UpdateCaseUseCase,
	transicion *appcasos.TransitionStateUseCase,
	repo domain.CasoRepository,
) *CasosHandler {
	return &CasosHandler{crear: crear, actualizar: actualizar, transicion: transicion, repo: repo}
}

// ── JSON response types ──────────────────────────────────────────────────────

type casoListItemJSON struct {
	ID             string  `json:"id"`
	BancoID        string  `json:"banco_id"`
	BancoNombre    string  `json:"banco_nombre"`
	ClienteID      string  `json:"cliente_id"`
	ClienteRUT     string  `json:"cliente_rut"`
	ClienteNombre  string  `json:"cliente_nombre"`
	AbogadoID      *string `json:"abogado_id"`
	NumeroOT       *string `json:"numero_ot"`
	Estado         string  `json:"estado"`
	FechaDJ        string  `json:"fecha_dj"`
	DenunciaValida bool    `json:"denuncia_valida"`
	CreatedAt      string  `json:"created_at"`
}

type casoJSON struct {
	ID             string  `json:"id"`
	EstudioID      string  `json:"estudio_id"`
	BancoID        string  `json:"banco_id"`
	ClienteID      string  `json:"cliente_id"`
	AbogadoID      *string `json:"abogado_id"`
	NumeroOT       *string `json:"numero_ot"`
	Estado         string  `json:"estado"`
	FechaDJ        string  `json:"fecha_dj"`
	FechaDenuncia  *string `json:"fecha_denuncia"`
	DenunciaValida bool    `json:"denuncia_valida"`
	MotivoTermino  *string `json:"motivo_termino"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type clienteJSON struct {
	ID       string  `json:"id"`
	RUT      string  `json:"rut"`
	Nombre   string  `json:"nombre"`
	Contacto *string `json:"contacto"`
}

type operacionJSON struct {
	ID        string   `json:"id"`
	CasoID    string   `json:"caso_id"`
	MedioPago string   `json:"medio_pago"`
	Relacion  string   `json:"relacion"`
	MontoCLP  int64    `json:"monto_clp"`
	MontoUF   *float64 `json:"monto_uf"`
	FechaOp   string   `json:"fecha_op"`
}

type casoDetalleJSON struct {
	Caso        casoJSON        `json:"caso"`
	Cliente     clienteJSON     `json:"cliente"`
	Operaciones []operacionJSON `json:"operaciones"`
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func toCasoListItemJSON(item *domain.CasoListItem) casoListItemJSON {
	return casoListItemJSON{
		ID:             item.ID,
		BancoID:        item.BancoID,
		BancoNombre:    item.BancoNombre,
		ClienteID:      item.ClienteID,
		ClienteRUT:     item.ClienteRUT,
		ClienteNombre:  item.ClienteNombre,
		AbogadoID:      item.AbogadoID,
		NumeroOT:       item.NumeroOT,
		Estado:         string(item.Estado),
		FechaDJ:        item.FechaDJ.Format("2006-01-02"),
		DenunciaValida: item.DenunciaValida,
		CreatedAt:      item.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func toCasoDetalleJSON(d *domain.CasoDetalle) casoDetalleJSON {
	c := d.Caso
	var fechaDen *string
	if c.FechaDenuncia != nil {
		s := c.FechaDenuncia.Format("2006-01-02")
		fechaDen = &s
	}
	ops := make([]operacionJSON, 0, len(d.Operaciones))
	for _, op := range d.Operaciones {
		ops = append(ops, operacionJSON{
			ID:        op.ID,
			CasoID:    op.CasoID,
			MedioPago: op.MedioPago,
			Relacion:  op.Relacion,
			MontoCLP:  op.MontoCLP,
			MontoUF:   op.MontoUF,
			FechaOp:   op.FechaOp.Format("2006-01-02"),
		})
	}
	return casoDetalleJSON{
		Caso: casoJSON{
			ID:             c.ID,
			EstudioID:      c.EstudioID,
			BancoID:        c.BancoID,
			ClienteID:      c.ClienteID,
			AbogadoID:      c.AbogadoID,
			NumeroOT:       c.NumeroOT,
			Estado:         string(c.Estado),
			FechaDJ:        c.FechaDJ.Format("2006-01-02"),
			FechaDenuncia:  fechaDen,
			DenunciaValida: c.DenunciaValida,
			MotivoTermino:  c.MotivoTermino,
			CreatedAt:      c.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      c.UpdatedAt.UTC().Format(time.RFC3339),
		},
		Cliente: clienteJSON{
			ID:       d.Cliente.ID,
			RUT:      d.Cliente.RUT,
			Nombre:   d.Cliente.Nombre,
			Contacto: d.Cliente.Contacto,
		},
		Operaciones: ops,
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func (h *CasosHandler) Listar(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	bancoIDs := middleware.BancoIDsFromCtx(r.Context())

	items, total, err := h.repo.ListRich(r.Context(), estudioID, domain.CaseFilters{
		BancoIDs: bancoIDs,
		Limit:    50,
	})
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	resp := make([]casoListItemJSON, 0, len(items))
	for _, item := range items {
		resp = append(resp, toCasoListItemJSON(item))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"casos": resp, "total": total})
}

func (h *CasosHandler) Crear(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	usuarioID := middleware.UsuarioIDFromCtx(r.Context())

	var req struct {
		BancoID         string  `json:"banco_id"`
		ClienteRUT      string  `json:"cliente_rut"`
		ClienteNombre   string  `json:"cliente_nombre"`
		ClienteContacto *string `json:"cliente_contacto"`
		FechaDJ         string  `json:"fecha_dj"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.BancoID == "" || req.ClienteRUT == "" || req.ClienteNombre == "" || req.FechaDJ == "" {
		http.Error(w, `{"error":"banco_id, cliente_rut, cliente_nombre and fecha_dj are required"}`, http.StatusBadRequest)
		return
	}

	fechaDJ, err := time.Parse("2006-01-02", req.FechaDJ)
	if err != nil {
		http.Error(w, `{"error":"fecha_dj must be YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}
	if fechaDJ.After(time.Now()) {
		http.Error(w, `{"error":"fecha_dj cannot be in the future"}`, http.StatusBadRequest)
		return
	}

	detalle, err := h.crear.Execute(r.Context(), appcasos.CreateCaseInput{
		EstudioID:       estudioID,
		BancoID:         req.BancoID,
		ClienteRUT:      req.ClienteRUT,
		ClienteNombre:   req.ClienteNombre,
		ClienteContacto: req.ClienteContacto,
		FechaDJ:         fechaDJ,
		UsuarioID:       usuarioID,
	})
	if err != nil {
		http.Error(w, `{"error":"could not create caso"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toCasoDetalleJSON(detalle))
}

func (h *CasosHandler) Obtener(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	id := chi.URLParam(r, "id")

	detalle, err := h.repo.GetDetalle(r.Context(), estudioID, id)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toCasoDetalleJSON(detalle))
}

func (h *CasosHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	usuarioID := middleware.UsuarioIDFromCtx(r.Context())
	id := chi.URLParam(r, "id")

	var req struct {
		AbogadoID      *string `json:"abogado_id"`
		NumeroOT       *string `json:"numero_ot"`
		DenunciaValida *bool   `json:"denuncia_valida"`
		FechaDenuncia  *string `json:"fecha_denuncia"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	input := appcasos.UpdateCaseInput{
		EstudioID:      estudioID,
		CasoID:         id,
		UsuarioID:      usuarioID,
		AbogadoID:      req.AbogadoID,
		NumeroOT:       req.NumeroOT,
		DenunciaValida: req.DenunciaValida,
	}
	if req.FechaDenuncia != nil {
		t, err := time.Parse("2006-01-02", *req.FechaDenuncia)
		if err != nil {
			http.Error(w, `{"error":"fecha_denuncia must be YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
		input.FechaDenuncia = &t
	}

	detalle, err := h.actualizar.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, `{"error":"could not update caso"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toCasoDetalleJSON(detalle))
}

func (h *CasosHandler) Transicionar(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	usuarioID := middleware.UsuarioIDFromCtx(r.Context())
	id := chi.URLParam(r, "id")

	var req struct {
		Estado        string  `json:"estado"`
		MotivoTermino *string `json:"motivo_termino"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	err := h.transicion.Execute(r.Context(), appcasos.TransitionStateInput{
		EstudioID:         estudioID,
		CasoID:            id,
		UsuarioID:         usuarioID,
		NewState:          estado.Estado(req.Estado),
		TerminationReason: req.MotivoTermino,
	})
	if err != nil {
		if isBadRequest(err, req.Estado) {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"error":"could not transition estado"}`, http.StatusInternalServerError)
		return
	}

	detalle, err := h.repo.GetDetalle(r.Context(), estudioID, id)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toCasoDetalleJSON(detalle))
}

// isBadRequest returns true if the error is a domain validation error, not a server error.
func isBadRequest(err error, _ string) bool {
	msg := err.Error()
	return msg == "invalid transition" ||
		msg == "valid denuncia required to enter JUDICIALIZACION" ||
		msg == "termination reason is required when closing a caso as TERMINADO"
}

// Ensure old signature compiles — NewCasosHandler replaces the old zero-arg version.
var _ = (*caso.Caso)(nil)
