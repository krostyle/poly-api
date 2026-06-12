package handlers

import (
	"encoding/json"
	"errors"
	"log"
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
	historial  domain.HistorialReader
}

func NewCasosHandler(
	crear *appcasos.CreateCaseUseCase,
	actualizar *appcasos.UpdateCaseUseCase,
	transicion *appcasos.TransitionStateUseCase,
	repo domain.CasoRepository,
	historial domain.HistorialReader,
) *CasosHandler {
	return &CasosHandler{crear: crear, actualizar: actualizar, transicion: transicion, repo: repo, historial: historial}
}

// ── JSON response types ──────────────────────────────────────────────────────

type casoListItemJSON struct {
	ID             string `json:"id"`
	BancoID        string `json:"banco_id"`
	BancoNombre    string `json:"banco_nombre"`
	ClienteID      string `json:"cliente_id"`
	ClienteRUT     string `json:"cliente_rut"`
	ClienteNombre  string `json:"cliente_nombre"`
	AbogadoID      *string `json:"abogado_id"`
	NumeroOT       *string `json:"numero_ot"`
	Estado         string  `json:"estado"`
	FechaDJ        *string `json:"fecha_dj"`
	EstadoDenuncia string  `json:"estado_denuncia"`
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
	FechaDJ        *string `json:"fecha_dj"`
	FechaDenuncia  *string `json:"fecha_denuncia"`
	EstadoDenuncia string  `json:"estado_denuncia"`
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

func formatDatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

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
		FechaDJ:        formatDatePtr(item.FechaDJ),
		EstadoDenuncia: string(item.EstadoDenuncia),
		CreatedAt:      item.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func toCasoDetalleJSON(d *domain.CasoDetalle) casoDetalleJSON {
	c := d.Caso
	fechaDen := formatDatePtr(c.FechaDenuncia)
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
			FechaDJ:        formatDatePtr(c.FechaDJ),
			FechaDenuncia:  fechaDen,
			EstadoDenuncia: string(c.EstadoDenuncia),
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
	usuarioID := middleware.UsuarioIDFromCtx(r.Context())
	bancoIDs := middleware.BancoIDsFromCtx(r.Context())

	q := r.URL.Query()
	filters := domain.CaseFilters{
		BancoIDs:      bancoIDs,
		BancoIDFilter: q.Get("banco_id"),
		Query:         q.Get("q"),
	}
	if est := q.Get("estado"); est != "" {
		e := estado.Estado(est)
		filters.Estado = &e
	}
	if aid := q.Get("abogado_id"); aid != "" {
		if aid == "me" {
			filters.AbogadoID = &usuarioID
		} else {
			filters.AbogadoID = &aid
		}
	}
	if q.Get("excluir_cierre") != "false" {
		filters.ExcluirCierre = true
	}

	items, total, err := h.repo.ListRich(r.Context(), estudioID, filters)
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
		FechaDJ         *string `json:"fecha_dj"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.BancoID == "" || req.ClienteRUT == "" || req.ClienteNombre == "" {
		http.Error(w, `{"error":"banco, RUT y nombre del cliente son obligatorios"}`, http.StatusBadRequest)
		return
	}

	if req.FechaDJ == nil || *req.FechaDJ == "" {
		http.Error(w, `{"error":"fecha_dj es obligatoria"}`, http.StatusBadRequest)
		return
	}
	fechaDJ, err := time.Parse("2006-01-02", *req.FechaDJ)
	if err != nil {
		http.Error(w, `{"error":"fecha_dj debe tener formato YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}
	if fechaDJ.After(time.Now()) {
		http.Error(w, `{"error":"fecha_dj no puede ser una fecha futura"}`, http.StatusBadRequest)
		return
	}

	input := appcasos.CreateCaseInput{
		EstudioID:       estudioID,
		BancoID:         req.BancoID,
		ClienteRUT:      req.ClienteRUT,
		ClienteNombre:   req.ClienteNombre,
		ClienteContacto: req.ClienteContacto,
		FechaDJ:         fechaDJ,
		UsuarioID:       usuarioID,
	}

	detalle, err := h.crear.Execute(r.Context(), input)
	if err != nil {
		log.Printf("[Crear] error: %v", err)
		http.Error(w, `{"error":"No se pudo crear el caso. Intenta nuevamente."}`, http.StatusInternalServerError)
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
		EstadoDenuncia *string `json:"estado_denuncia"`
		FechaDenuncia  *string `json:"fecha_denuncia"`
		FechaDJ        *string `json:"fecha_dj"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	// Tramitadores may only update fecha fields.
	rol := middleware.RolFromCtx(r.Context())
	if rol == "TRAMITADOR" {
		req.AbogadoID = nil
		req.NumeroOT = nil
		req.EstadoDenuncia = nil
	}

	// Validate estado_denuncia if provided.
	var estDenuncia *caso.EstadoDenuncia
	if req.EstadoDenuncia != nil {
		if !caso.IsValidEstadoDenuncia(*req.EstadoDenuncia) {
			http.Error(w, `{"error":"estado_denuncia debe ser PENDIENTE, ACOGIDA o RECHAZADA"}`, http.StatusBadRequest)
			return
		}
		ed := caso.EstadoDenuncia(*req.EstadoDenuncia)
		estDenuncia = &ed
	}

	input := appcasos.UpdateCaseInput{
		EstudioID:      estudioID,
		CasoID:         id,
		UsuarioID:      usuarioID,
		AbogadoID:      req.AbogadoID,
		NumeroOT:       req.NumeroOT,
		EstadoDenuncia: estDenuncia,
	}
	if req.FechaDenuncia != nil {
		if *req.FechaDenuncia == "" {
			input.ClearFechaDenuncia = true
		} else {
			t, err := time.Parse("2006-01-02", *req.FechaDenuncia)
			if err != nil {
				http.Error(w, `{"error":"fecha_denuncia must be YYYY-MM-DD"}`, http.StatusBadRequest)
				return
			}
			input.FechaDenuncia = &t
		}
	}
	if req.FechaDJ != nil && *req.FechaDJ != "" {
		t, err := time.Parse("2006-01-02", *req.FechaDJ)
		if err != nil {
			http.Error(w, `{"error":"fecha_dj debe tener formato YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
		input.FechaDJ = &t
	}

	detalle, err := h.actualizar.Execute(r.Context(), input)
	if err != nil {
		log.Printf("[Actualizar] error: %v", err)
		http.Error(w, `{"error":"No se pudieron guardar los cambios. Intenta nuevamente."}`, http.StatusInternalServerError)
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
		Forzar        bool    `json:"forzar"`
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
		Forzar:            req.Forzar,
	})
	if err != nil {
		if isBadRequest(err) {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, `{"error":"Error al procesar la transición de estado"}`, http.StatusInternalServerError)
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

func (h *CasosHandler) Historial(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	id := chi.URLParam(r, "id")

	entries, err := h.historial.ListByCaso(r.Context(), estudioID, id)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	type entryJSON struct {
		ID            string         `json:"id"`
		Accion        string         `json:"accion"`
		Detalle       map[string]any `json:"detalle"`
		UsuarioNombre string         `json:"usuario_nombre"`
		CreatedAt     string         `json:"created_at"`
	}
	resp := make([]entryJSON, 0, len(entries))
	for _, e := range entries {
		det := e.Detalle
		if det == nil {
			det = map[string]any{}
		}
		resp = append(resp, entryJSON{
			ID:            e.ID,
			Accion:        e.Accion,
			Detalle:       det,
			UsuarioNombre: e.UsuarioNombre,
			CreatedAt:     e.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"historial": resp})
}

func (h *CasosHandler) Eliminar(w http.ResponseWriter, r *http.Request) {
	rol := middleware.RolFromCtx(r.Context())
	if rol != "ADMIN" {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	estudioID := middleware.EstudioIDFromCtx(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.repo.Delete(r.Context(), estudioID, id); err != nil {
		log.Printf("[Eliminar] error: %v", err)
		http.Error(w, `{"error":"No se pudo eliminar el caso. Intenta nuevamente."}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// isBadRequest returns true if the error is a domain validation error, not a server error.
func isBadRequest(err error) bool {
	return errors.Is(err, estado.ErrTransicionNoPermitida) ||
		errors.Is(err, appcasos.ErrMotivoTerminoRequerido) ||
		errors.Is(err, appcasos.ErrMotivoTerminoInvalido) ||
		errors.Is(err, appcasos.ErrDenunciaRechazadaRequerida) ||
		errors.Is(err, appcasos.ErrDenunciaAcogidaRequerida)
}

// keep import used only for interface check
var _ = (*caso.Caso)(nil)
