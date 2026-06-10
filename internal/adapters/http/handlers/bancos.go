package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"poly.app/api/internal/adapters/http/middleware"
	"poly.app/api/internal/domain"
)

type BancosHandler struct {
	bancos   domain.BancoRepository
	usuarios domain.UsuarioRepository
}

func NewBancosHandler(bancos domain.BancoRepository, usuarios domain.UsuarioRepository) *BancosHandler {
	return &BancosHandler{bancos: bancos, usuarios: usuarios}
}

type bancoDetalleJSON struct {
	ID        string    `json:"id"`
	Nombre    string    `json:"nombre"`
	CreatedAt time.Time `json:"created_at"`
}

type usuarioBancoJSON struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Rol    string `json:"rol"`
}

func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	if middleware.RolFromCtx(r.Context()) != "ADMIN" {
		http.Error(w, `{"error":"forbidden: se requiere rol ADMIN"}`, http.StatusForbidden)
		return false
	}
	return true
}

func toBancoDetalle(b *domain.Banco) bancoDetalleJSON {
	return bancoDetalleJSON{ID: b.ID, Nombre: b.Nombre, CreatedAt: b.CreatedAt}
}

// GET /v1/bancos
func (h *BancosHandler) Listar(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	bancos, err := h.bancos.List(r.Context(), estudioID)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	out := make([]bancoDetalleJSON, 0, len(bancos))
	for _, b := range bancos {
		out = append(out, toBancoDetalle(b))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"bancos": out})
}

// POST /v1/bancos — ADMIN only
func (h *BancosHandler) Crear(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	var req struct {
		Nombre string `json:"nombre"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Nombre == "" {
		http.Error(w, `{"error":"nombre is required"}`, http.StatusBadRequest)
		return
	}
	banco, err := h.bancos.Create(r.Context(), estudioID, req.Nombre)
	if err != nil {
		http.Error(w, `{"error":"could not create banco"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toBancoDetalle(banco))
}

// PATCH /v1/bancos/{id} — ADMIN only
func (h *BancosHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	id := chi.URLParam(r, "id")
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	var req struct {
		Nombre string `json:"nombre"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Nombre == "" {
		http.Error(w, `{"error":"nombre is required"}`, http.StatusBadRequest)
		return
	}
	banco, err := h.bancos.Update(r.Context(), estudioID, id, req.Nombre)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, `{"error":"banco not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"could not update banco"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toBancoDetalle(banco))
}

// DELETE /v1/bancos/{id} — ADMIN only
func (h *BancosHandler) Eliminar(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	id := chi.URLParam(r, "id")
	estudioID := middleware.EstudioIDFromCtx(r.Context())

	if _, err := h.bancos.GetByID(r.Context(), estudioID, id); err != nil {
		http.Error(w, `{"error":"banco not found"}`, http.StatusNotFound)
		return
	}
	hasCasos, err := h.bancos.HasCasos(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	if hasCasos {
		http.Error(w, `{"error":"el banco tiene casos asociados y no puede eliminarse"}`, http.StatusConflict)
		return
	}
	if err := h.bancos.Delete(r.Context(), estudioID, id); err != nil {
		http.Error(w, `{"error":"could not delete banco"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /v1/bancos/{id}/usuarios — ADMIN only
func (h *BancosHandler) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	id := chi.URLParam(r, "id")
	estudioID := middleware.EstudioIDFromCtx(r.Context())

	if _, err := h.bancos.GetByID(r.Context(), estudioID, id); err != nil {
		http.Error(w, `{"error":"banco not found"}`, http.StatusNotFound)
		return
	}
	usuarios, err := h.bancos.ListUsuarios(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	out := make([]usuarioBancoJSON, 0, len(usuarios))
	for _, u := range usuarios {
		out = append(out, usuarioBancoJSON{ID: u.ID, Nombre: u.Nombre, Email: u.Email, Rol: u.Rol})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"usuarios": out})
}

// POST /v1/bancos/{id}/usuarios — ADMIN only
func (h *BancosHandler) AsignarUsuario(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	bancoID := chi.URLParam(r, "id")
	estudioID := middleware.EstudioIDFromCtx(r.Context())

	if _, err := h.bancos.GetByID(r.Context(), estudioID, bancoID); err != nil {
		http.Error(w, `{"error":"banco not found"}`, http.StatusNotFound)
		return
	}
	var req struct {
		UsuarioID string `json:"usuario_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UsuarioID == "" {
		http.Error(w, `{"error":"usuario_id is required"}`, http.StatusBadRequest)
		return
	}
	// verify usuario belongs to same estudio
	if _, err := h.usuarios.GetByEstudioAndID(r.Context(), estudioID, req.UsuarioID); err != nil {
		http.Error(w, `{"error":"usuario not found in this estudio"}`, http.StatusNotFound)
		return
	}
	if err := h.bancos.AsignarUsuario(r.Context(), bancoID, req.UsuarioID); err != nil {
		http.Error(w, `{"error":"could not assign usuario"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /v1/bancos/{id}/usuarios/{usuarioId} — ADMIN only
func (h *BancosHandler) DesasignarUsuario(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	bancoID := chi.URLParam(r, "id")
	usuarioID := chi.URLParam(r, "usuarioId")
	estudioID := middleware.EstudioIDFromCtx(r.Context())

	if _, err := h.bancos.GetByID(r.Context(), estudioID, bancoID); err != nil {
		http.Error(w, `{"error":"banco not found"}`, http.StatusNotFound)
		return
	}
	if err := h.bancos.DesasignarUsuario(r.Context(), bancoID, usuarioID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, `{"error":"asignación not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"could not remove assignment"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /v1/usuarios — lista todos los usuarios del estudio (ADMIN only, para AsignarUsuarioDialog)
func (h *BancosHandler) ListarUsuariosEstudio(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	usuarios, err := h.usuarios.ListByEstudio(r.Context(), estudioID)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	out := make([]usuarioBancoJSON, 0, len(usuarios))
	for _, u := range usuarios {
		out = append(out, usuarioBancoJSON{ID: u.ID, Nombre: u.Nombre, Email: u.Email, Rol: u.Rol})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"usuarios": out})
}
