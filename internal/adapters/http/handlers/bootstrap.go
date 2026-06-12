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
	uc       *appauth.BootstrapUseCase
	estudios domain.EstudioRepository
	usuarios domain.UsuarioRepository
	bancos   domain.BancoRepository
}

func NewBootstrapHandler(
	uc *appauth.BootstrapUseCase,
	estudios domain.EstudioRepository,
	usuarios domain.UsuarioRepository,
	bancos domain.BancoRepository,
) *BootstrapHandler {
	return &BootstrapHandler{uc: uc, estudios: estudios, usuarios: usuarios, bancos: bancos}
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
	ID                   string `json:"id"`
	Nombre               string `json:"nombre"`
	Email                string `json:"email"`
	Rol                  string `json:"rol"`
	OnboardingCompletado bool   `json:"onboarding_completado"`
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
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok || claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	estudioID := middleware.EstudioIDFromCtx(r.Context())

	usuario, err := h.usuarios.GetByClerkUserID(r.Context(), claims.Subject)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	bancos, err := h.bancos.List(r.Context(), estudioID)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	estudio, err := h.estudios.GetByClerkOrgID(r.Context(), claims.ActiveOrganizationID)
	if err != nil {
		http.Error(w, `{"error":"estudio not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profileResponse{
		Estudio: estudioJSON{ID: estudio.ID, Nombre: estudio.Nombre},
		Usuario: usuarioJSON{ID: usuario.ID, Nombre: usuario.Nombre, Email: usuario.Email, Rol: usuario.Rol, OnboardingCompletado: usuario.OnboardingCompletado},
		Bancos:  toBancoJSONs(bancos),
	})
}

func toBancoJSONs(bs []*domain.Banco) []bancoJSON {
	out := make([]bancoJSON, 0, len(bs))
	for _, b := range bs {
		out = append(out, bancoJSON{ID: b.ID, Nombre: b.Nombre})
	}
	return out
}

func toProfileResponse(e *domain.Estudio, u *domain.Usuario, bs []*domain.Banco) profileResponse {
	bancos := make([]bancoJSON, 0, len(bs))
	for _, b := range bs {
		bancos = append(bancos, bancoJSON{ID: b.ID, Nombre: b.Nombre})
	}
	return profileResponse{
		Estudio: estudioJSON{ID: e.ID, Nombre: e.Nombre},
		Usuario: usuarioJSON{ID: u.ID, Nombre: u.Nombre, Email: u.Email, Rol: u.Rol, OnboardingCompletado: u.OnboardingCompletado},
		Bancos:  bancos,
	}
}
