package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"poly.app/api/internal/adapters/http/middleware"
	"poly.app/api/internal/domain"
)

type ClientesHandler struct {
	repo domain.ClienteRepository
}

func NewClientesHandler(repo domain.ClienteRepository) *ClientesHandler {
	return &ClientesHandler{repo: repo}
}

func (h *ClientesHandler) Crear(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())

	var req struct {
		BancoID  string  `json:"banco_id"`
		RUT      string  `json:"rut"`
		Nombre   string  `json:"nombre"`
		Contacto *string `json:"contacto"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.BancoID == "" || req.RUT == "" || req.Nombre == "" {
		http.Error(w, `{"error":"banco_id, rut and nombre are required"}`, http.StatusBadRequest)
		return
	}

	cliente, err := h.repo.Upsert(r.Context(), domain.UpsertClienteInput{
		EstudioID: estudioID,
		BancoID:   req.BancoID,
		RUT:       req.RUT,
		Nombre:    req.Nombre,
		Contacto:  req.Contacto,
	})
	if err != nil {
		http.Error(w, `{"error":"could not create cliente"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toClienteDetailJSON(cliente))
}

func (h *ClientesHandler) Obtener(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	id := chi.URLParam(r, "id")

	cliente, err := h.repo.GetByID(r.Context(), estudioID, id)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toClienteDetailJSON(cliente))
}

type clienteDetailJSON struct {
	ID        string  `json:"id"`
	EstudioID string  `json:"estudio_id"`
	BancoID   string  `json:"banco_id"`
	RUT       string  `json:"rut"`
	Nombre    string  `json:"nombre"`
	Contacto  *string `json:"contacto"`
}

func toClienteDetailJSON(c *domain.Cliente) clienteDetailJSON {
	return clienteDetailJSON{
		ID:        c.ID,
		EstudioID: c.EstudioID,
		BancoID:   c.BancoID,
		RUT:       c.RUT,
		Nombre:    c.Nombre,
		Contacto:  c.Contacto,
	}
}
