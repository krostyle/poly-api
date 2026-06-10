package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/jackc/pgx/v5/pgxpool"
	appauth "poly.app/api/internal/application/auth"
	"poly.app/api/internal/adapters/http/handlers"
	"poly.app/api/internal/adapters/http/middleware"
	"poly.app/api/internal/adapters/persistence"
)

func NewRouter(pool *pgxpool.Pool) http.Handler {
	// Repos
	estudiosRepo := persistence.NewEstudioRepo(pool)
	usuariosRepo := persistence.NewUsuarioRepo(pool)
	bancosRepo := persistence.NewBancoRepo(pool)

	// Use cases
	bootstrapUC := appauth.NewBootstrapUseCase(estudiosRepo, usuariosRepo, bancosRepo)

	// Handlers
	bootstrapH := handlers.NewBootstrapHandler(bootstrapUC)
	adminH := handlers.NewAdminHandler(bancosRepo, usuariosRepo)
	casosH := handlers.NewCasosHandler()
	plazosH := handlers.NewPlazosHandler()

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)

	// ── Public ───────────────────────────────────────────────────────────────
	r.Get("/health", handlers.Health)

	// ── JWT-verified, no tenant scope (bootstrap) ────────────────────────────
	r.Group(func(r chi.Router) {
		r.Use(clerkhttp.WithHeaderAuthorization())
		r.Post("/v1/bootstrap", bootstrapH.Bootstrap)
	})

	// ── Fully protected (JWT + tenant scope) ─────────────────────────────────
	r.Group(func(r chi.Router) {
		r.Use(clerkhttp.WithHeaderAuthorization())
		r.Use(middleware.RequireAuth(pool))
		r.Use(middleware.RequireTenantScope)

		r.Get("/v1/me", bootstrapH.Me)

		r.Route("/v1/admin", func(r chi.Router) {
			r.Post("/bancos", adminH.CreateBanco)
			r.Post("/usuarios/{id}/bancos", adminH.AssignBancoToUsuario)
		})

		r.Route("/v1/casos", func(r chi.Router) {
			r.Get("/", casosH.Listar)
			r.Post("/", casosH.Crear)
			r.Get("/{id}", casosH.Obtener)
			r.Post("/{id}/transicion", casosH.Transicionar)
		})

		r.Route("/v1/casos/{casoID}/plazos", func(r chi.Router) {
			r.Get("/", plazosH.ListarPorCaso)
		})
	})

	return r
}
