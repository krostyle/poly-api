package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"poly.app/api/internal/adapters/http/middleware"
	"poly.app/api/internal/domain"
)

type AdminHandler struct {
	bancos   domain.BancoRepository
	usuarios domain.UsuarioRepository
}

func NewAdminHandler(bancos domain.BancoRepository, usuarios domain.UsuarioRepository) *AdminHandler {
	return &AdminHandler{bancos: bancos, usuarios: usuarios}
}

type createBancoRequest struct {
	Nombre string `json:"nombre"`
}

func (h *AdminHandler) CreateBanco(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())

	var req createBancoRequest
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
	json.NewEncoder(w).Encode(bancoJSON{ID: banco.ID, Nombre: banco.Nombre})
}

func (h *AdminHandler) AssignBancoToUsuario(w http.ResponseWriter, r *http.Request) {
	usuarioID := chi.URLParam(r, "id")

	var req struct {
		BancoID string `json:"banco_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.BancoID == "" {
		http.Error(w, `{"error":"banco_id is required"}`, http.StatusBadRequest)
		return
	}

	if err := h.bancos.AsignarUsuario(r.Context(), req.BancoID, usuarioID); err != nil {
		http.Error(w, `{"error":"could not assign banco"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
