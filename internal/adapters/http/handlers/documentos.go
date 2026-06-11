package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"poly.app/api/internal/adapters/http/middleware"
	appdocs "poly.app/api/internal/application/documentos"
	"poly.app/api/internal/domain"
)

type DocumentosHandler struct {
	subir *appdocs.SubirDocumentoUseCase
	repo  domain.DocumentoRepository
}

func NewDocumentosHandler(subir *appdocs.SubirDocumentoUseCase, repo domain.DocumentoRepository) *DocumentosHandler {
	return &DocumentosHandler{subir: subir, repo: repo}
}

type documentoResponse struct {
	ID        string  `json:"id"`
	Tipo      string  `json:"tipo"`
	Nombre    string  `json:"nombre"`
	BlobURL   string  `json:"blob_url"`
	SubidoPor *string `json:"subido_por,omitempty"`
	CreatedAt string  `json:"created_at"`
}

func (h *DocumentosHandler) Subir(w http.ResponseWriter, r *http.Request) {
	casoID := chi.URLParam(r, "id")
	usuarioID := middleware.UsuarioIDFromCtx(r.Context())

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, `{"error":"archivo demasiado grande (máx 20 MB)"}`, http.StatusRequestEntityTooLarge)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"campo 'file' requerido"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, 20<<20+1))
	if err != nil {
		http.Error(w, `{"error":"error al leer archivo"}`, http.StatusInternalServerError)
		return
	}
	if len(content) > 20<<20 {
		http.Error(w, `{"error":"archivo supera el límite de 20 MB"}`, http.StatusRequestEntityTooLarge)
		return
	}

	tipo := r.FormValue("tipo")
	if tipo == "" {
		tipo = "OTRO"
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	uid := usuarioID
	doc, err := h.subir.Execute(r.Context(), appdocs.SubirDocumentoInput{
		CasoID:      casoID,
		Tipo:        tipo,
		Nombre:      header.Filename,
		Content:     content,
		ContentType: contentType,
		SubidoPor:   &uid,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "tipo de documento no permitido" || err.Error() == "tipo de archivo no permitido" {
			status = http.StatusBadRequest
		}
		if err.Error() == "archivo supera el límite de 20 MB" {
			status = http.StatusRequestEntityTooLarge
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toDocResponse(doc))
}

func (h *DocumentosHandler) Listar(w http.ResponseWriter, r *http.Request) {
	casoID := chi.URLParam(r, "id")

	docs, err := h.repo.ListByCaso(r.Context(), casoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]documentoResponse, 0, len(docs))
	for _, d := range docs {
		result = append(result, toDocResponse(d))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"documentos": result})
}

func toDocResponse(d *domain.Documento) documentoResponse {
	return documentoResponse{
		ID:        d.ID,
		Tipo:      d.Tipo,
		Nombre:    d.Nombre,
		BlobURL:   d.BlobURL,
		SubidoPor: d.SubidoPor,
		CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
