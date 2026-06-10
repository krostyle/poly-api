package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	appauth "poly.app/api/internal/application/auth"
	"poly.app/api/internal/adapters/http/middleware"
	"poly.app/api/internal/domain"
)

type BootstrapHandler struct {
	uc *appauth.BootstrapUseCase
}

func NewBootstrapHandler(uc *appauth.BootstrapUseCase) *BootstrapHandler {
	return &BootstrapHandler{uc: uc}
}

type bootstrapRequest struct {
	OrgName   string `json:"org_name"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

type profileResponse struct {
	Estudio estudioJSON  `json:"estudio"`
	Usuario usuarioJSON  `json:"usuario"`
	Bancos  []bancoJSON  `json:"bancos"`
}

type estudioJSON struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}

type usuarioJSON struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Rol    string `json:"rol"`
}

type bancoJSON struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}

func (h *BootstrapHandler) Bootstrap(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.Subject == "" || claims.ActiveOrganizationID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req bootstrapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	out, err := h.uc.Execute(r.Context(), appauth.BootstrapInput{
		ClerkOrgID:   claims.ActiveOrganizationID,
		ClerkUserID:  claims.Subject,
		ClerkOrgRole: claims.ActiveOrganizationRole,
		OrgName:      req.OrgName,
		UserName:     req.UserName,
		UserEmail:    req.UserEmail,
	})
	if err != nil {
		http.Error(w, `{"error":"bootstrap failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toProfileResponse(out.Estudio, out.Usuario, out.Bancos))
}

func (h *BootstrapHandler) Me(w http.ResponseWriter, r *http.Request) {
	estudioID := middleware.EstudioIDFromCtx(r.Context())
	usuarioID := middleware.UsuarioIDFromCtx(r.Context())
	_ = estudioID
	_ = usuarioID
	// Full Me implementation will query estudio+usuario+bancos from context IDs.
	// For now return the IDs already in context as a minimal response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"estudio_id": estudioID,
		"usuario_id": usuarioID,
	})
}

func toProfileResponse(e *domain.Estudio, u *domain.Usuario, bs []*domain.Banco) profileResponse {
	bancos := make([]bancoJSON, 0, len(bs))
	for _, b := range bs {
		bancos = append(bancos, bancoJSON{ID: b.ID, Nombre: b.Nombre})
	}
	return profileResponse{
		Estudio: estudioJSON{ID: e.ID, Nombre: e.Nombre},
		Usuario: usuarioJSON{ID: u.ID, Nombre: u.Nombre, Email: u.Email, Rol: u.Rol},
		Bancos:  bancos,
	}
}
