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
	keyRol       contextKey = "rol"
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

func RolFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(keyRol).(string)
	return v
}

// RequireRol returns a middleware that allows only the listed roles.
func RequireRol(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rol := RolFromCtx(r.Context())
			for _, allowed := range roles {
				if rol == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		})
	}
}

// RequireAuth reads the Clerk session claims injected by clerkhttp.WithHeaderAuthorization,
// looks up the estudio and usuario in the DB, and injects the tenant scope into context.
//
// ADMINs receive the full list of estudio banco IDs so they can see all cases regardless
// of which bancos they are explicitly assigned to in usuarios_bancos.
func RequireAuth(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	estudiosRepo := persistence.NewEstudioRepo(pool)
	usuariosRepo := persistence.NewUsuarioRepo(pool)
	bancosRepo := persistence.NewBancoRepo(pool)

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

			var bancoIDs []string
			if usuario.Rol == "ADMIN" {
				// Admins see all cases across all bancos of the estudio.
				bancos, _ := bancosRepo.List(r.Context(), estudio.ID)
				for _, b := range bancos {
					bancoIDs = append(bancoIDs, b.ID)
				}
			} else {
				bancoIDs, _ = usuariosRepo.GetBancoIDs(r.Context(), usuario.ID)
			}

			ctx := context.WithValue(r.Context(), keyEstudioID, estudio.ID)
			ctx = context.WithValue(ctx, keyUsuarioID, usuario.ID)
			ctx = context.WithValue(ctx, keyBancoIDs, bancoIDs)
			ctx = context.WithValue(ctx, keyRol, usuario.Rol)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
