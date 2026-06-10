package middleware

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"poly.app/api/internal/adapters/persistence"
)

type contextKey string

const (
	keyEstudioID contextKey = "estudio_id"
	keyBancoIDs  contextKey = "banco_ids"
	keyUsuarioID contextKey = "usuario_id"
)

func EstudioIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(keyEstudioID).(string)
	return v
}

func BancoIDsFromCtx(ctx context.Context) []string {
	v, _ := ctx.Value(keyBancoIDs).([]string)
	return v
}

func UsuarioIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(keyUsuarioID).(string)
	return v
}

// RequireAuth reads the Clerk session claims injected by clerkhttp.WithHeaderAuthorization,
// looks up the estudio and usuario in the DB, and injects the tenant scope into context.
func RequireAuth(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	estudiosRepo := persistence.NewEstudioRepo(pool)
	usuariosRepo := persistence.NewUsuarioRepo(pool)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok || claims == nil || claims.Subject == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			if claims.ActiveOrganizationID == "" {
				http.Error(w, `{"error":"no active organization"}`, http.StatusUnauthorized)
				return
			}

			estudio, err := estudiosRepo.GetByClerkOrgID(r.Context(), claims.ActiveOrganizationID)
			if err != nil {
				http.Error(w, `{"error":"estudio not found — run bootstrap"}`, http.StatusUnauthorized)
				return
			}

			usuario, err := usuariosRepo.GetByClerkUserID(r.Context(), claims.Subject)
			if err != nil {
				http.Error(w, `{"error":"usuario not found — run bootstrap"}`, http.StatusUnauthorized)
				return
			}

			bancoIDs, _ := usuariosRepo.GetBancoIDs(r.Context(), usuario.ID)

			ctx := context.WithValue(r.Context(), keyEstudioID, estudio.ID)
			ctx = context.WithValue(ctx, keyUsuarioID, usuario.ID)
			ctx = context.WithValue(ctx, keyBancoIDs, bancoIDs)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
