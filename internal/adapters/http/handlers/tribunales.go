package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"poly.app/api/internal/domain"
)

type TribunalesHandler struct {
	repo domain.TribunalRepository
}

func NewTribunalesHandler(repo domain.TribunalRepository) *TribunalesHandler {
	return &TribunalesHandler{repo: repo}
}

type tribunalJSON struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Region string `json:"region"`
}

func (h *TribunalesHandler) Listar(w http.ResponseWriter, r *http.Request) {
	tribunales, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	resp := make([]tribunalJSON, 0, len(tribunales))
	for _, t := range tribunales {
		resp = append(resp, tribunalJSON{ID: t.ID, Nombre: t.Nombre, Region: t.Region})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"tribunales": resp})
}

func (h *TribunalesHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nombre string `json:"nombre"`
		Region string `json:"region"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Nombre == "" || req.Region == "" {
		http.Error(w, `{"error":"nombre y region son obligatorios"}`, http.StatusBadRequest)
		return
	}

	t, err := h.repo.Create(r.Context(), req.Nombre, req.Region)
	if err != nil {
		log.Printf("[Tribunales] crear error: %v", err)
		http.Error(w, `{"error":"No se pudo crear el tribunal. Intenta nuevamente."}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tribunalJSON{ID: t.ID, Nombre: t.Nombre, Region: t.Region})
}
