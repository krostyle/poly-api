package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkorginv "github.com/clerk/clerk-sdk-go/v2/organizationinvitation"
	clerkorgmem "github.com/clerk/clerk-sdk-go/v2/organizationmembership"
	"github.com/go-chi/chi/v5"
	"poly.app/api/internal/adapters/http/middleware"
	"poly.app/api/internal/domain"
)

type UsuariosHandler struct {
	usuarios domain.UsuarioRepository
}

func NewUsuariosHandler(usuarios domain.UsuarioRepository) *UsuariosHandler {
	return &UsuariosHandler{usuarios: usuarios}
}

var validRoles = map[string]bool{"ADMIN": true, "ABOGADO": true, "TRAMITADOR": true}

func polyRoleToClerk(rol string) string {
	if rol == "ADMIN" {
		return "org:admin"
	}
	return "org:member"
}

// POST /v1/usuarios/invitar — ADMIN only
// Sends a Clerk org invitation; when the user accepts and logs in, bootstrap auto-creates them.
func (h *UsuariosHandler) Invitar(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}

	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok || claims.ActiveOrganizationID == "" {
		http.Error(w, `{"error":"no active organization in token"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		Email string `json:"email"`
		Rol   string `json:"rol"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, `{"error":"email is required"}`, http.StatusBadRequest)
		return
	}
	if req.Rol == "" {
		req.Rol = "ABOGADO"
	}
	if !validRoles[req.Rol] {
		http.Error(w, `{"error":"rol must be ADMIN, ABOGADO or TRAMITADOR"}`, http.StatusBadRequest)
		return
	}

	clerkRole := polyRoleToClerk(req.Rol)
	_, err := clerkorginv.Create(r.Context(), &clerkorginv.CreateParams{
		OrganizationID: claims.ActiveOrganizationID,
		EmailAddress:   &req.Email,
		Role:           &clerkRole,
	})
	if err != nil {
		http.Error(w, `{"error":"could not send invitation: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "invitation sent"})
}

// PATCH /v1/me/rol — any authenticated user
// Lets the user choose their own role (ABOGADO or TRAMITADOR) on first login.
// Also marks onboarding_completado = true.
func (h *UsuariosHandler) CompletarOnboarding(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok || claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		Rol string `json:"rol"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.Rol != "ABOGADO" && req.Rol != "TRAMITADOR" {
		http.Error(w, `{"error":"el rol debe ser ABOGADO o TRAMITADOR"}`, http.StatusBadRequest)
		return
	}

	updated, err := h.usuarios.CompleteOnboarding(r.Context(), claims.Subject, req.Rol)
	if err != nil {
		http.Error(w, `{"error":"no se pudo actualizar el rol"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":                    updated.ID,
		"nombre":                updated.Nombre,
		"email":                 updated.Email,
		"rol":                   updated.Rol,
		"onboarding_completado": updated.OnboardingCompletado,
	})
}

// PATCH /v1/usuarios/:id/rol — ADMIN only
// Updates the user's Poly role and syncs the Clerk org membership role when ADMIN changes.
func (h *UsuariosHandler) ActualizarRol(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}

	id := chi.URLParam(r, "id")
	estudioID := middleware.EstudioIDFromCtx(r.Context())

	var req struct {
		Rol string `json:"rol"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !validRoles[req.Rol] {
		http.Error(w, `{"error":"rol must be ADMIN, ABOGADO or TRAMITADOR"}`, http.StatusBadRequest)
		return
	}

	// Fetch current user to get clerk_user_id and current role.
	current, err := h.usuarios.GetByEstudioAndID(r.Context(), estudioID, id)
	if err != nil {
		http.Error(w, `{"error":"usuario not found"}`, http.StatusNotFound)
		return
	}

	// Sync Clerk org membership when the ADMIN boundary is crossed.
	if (current.Rol == "ADMIN") != (req.Rol == "ADMIN") {
		claims, ok := clerk.SessionClaimsFromContext(r.Context())
		if ok && claims.ActiveOrganizationID != "" {
			newClerkRole := polyRoleToClerk(req.Rol)
			_, _ = clerkorgmem.Update(r.Context(), &clerkorgmem.UpdateParams{
				OrganizationID: claims.ActiveOrganizationID,
				UserID:         current.ClerkUserID,
				Role:           &newClerkRole,
			})
		}
	}

	updated, err := h.usuarios.UpdateRol(r.Context(), estudioID, id, req.Rol)
	if err != nil {
		http.Error(w, `{"error":"could not update rol"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":     updated.ID,
		"nombre": updated.Nombre,
		"email":  updated.Email,
		"rol":    updated.Rol,
	})
}
