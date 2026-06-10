package handlers

import "net/http"

type PlazosHandler struct{}

func NewPlazosHandler() *PlazosHandler {
	return &PlazosHandler{}
}

func (h *PlazosHandler) ListarPorCaso(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
