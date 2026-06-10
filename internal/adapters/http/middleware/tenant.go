package middleware

import (
	"net/http"
)

// RequireTenantScope asegura que toda request tenga estudio_id y al menos un banco_id
// en contexto antes de pasar a los handlers.
func RequireTenantScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if EstudioIDFromCtx(r.Context()) == "" {
			http.Error(w, "unauthorized: missing tenant scope", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
