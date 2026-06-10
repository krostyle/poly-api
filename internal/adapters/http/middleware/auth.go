package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	keyEstudioID contextKey = "estudio_id"
	keyBancoIDs  contextKey = "banco_ids"
	keyUsuarioID contextKey = "usuario_id"
)

// EstudioIDFromCtx extrae el estudio_id inyectado por el middleware de auth.
func EstudioIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(keyEstudioID).(string)
	return v
}

// BancoIDsFromCtx extrae los banco_ids habilitados para el usuario.
func BancoIDsFromCtx(ctx context.Context) []string {
	v, _ := ctx.Value(keyBancoIDs).([]string)
	return v
}

// UsuarioIDFromCtx extrae el usuario_id del contexto.
func UsuarioIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(keyUsuarioID).(string)
	return v
}

// RequireAuth verifica el JWT de Clerk e inyecta estudio_id y banco_ids en el contexto.
// Implementación completa se realiza al integrar el SDK de Clerk.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: validar JWT con clerk-sdk-go/v2, extraer org_id y user_id,
		// consultar DB para banco_ids habilitados, inyectar en contexto.
		next.ServeHTTP(w, r)
	})
}
