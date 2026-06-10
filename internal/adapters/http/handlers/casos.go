package handlers

import (
	"net/http"
)

// CasosHandler agrupa los handlers de casos.
// Se expande con los casos de uso en la siguiente fase.
type CasosHandler struct{}

func NewCasosHandler() *CasosHandler {
	return &CasosHandler{}
}

func (h *CasosHandler) Listar(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *CasosHandler) Crear(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *CasosHandler) Obtener(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *CasosHandler) Transicionar(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
