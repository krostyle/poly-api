package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/jackc/pgx/v5/pgxpool"
	appauth "poly.app/api/internal/application/auth"
	appcasos "poly.app/api/internal/application/casos"
	appops "poly.app/api/internal/application/operaciones"
	"poly.app/api/internal/adapters/http/handlers"
	"poly.app/api/internal/adapters/http/middleware"
	"poly.app/api/internal/adapters/persistence"
)

func NewRouter(pool *pgxpool.Pool) http.Handler {
	// ── Repos ────────────────────────────────────────────────────────────────
	estudiosRepo := persistence.NewEstudioRepo(pool)
	usuariosRepo := persistence.NewUsuarioRepo(pool)
	bancosRepo := persistence.NewBancoRepo(pool)
	casosRepo := persistence.NewCasoRepo(pool)
	clientesRepo := persistence.NewClienteRepo(pool)
	operacionesRepo := persistence.NewOperacionRepo(pool)
	auditRepo := persistence.NewAuditRepo(pool)

	// ── Use cases ────────────────────────────────────────────────────────────
	bootstrapUC := appauth.NewBootstrapUseCase(estudiosRepo, usuariosRepo, bancosRepo)
	createCaseUC := appcasos.NewCreateCaseUseCase(casosRepo, clientesRepo, auditRepo)
	updateCaseUC := appcasos.NewUpdateCaseUseCase(casosRepo, auditRepo)
	transicionUC := appcasos.NewTransitionStateUseCase(casosRepo, auditRepo)
	agregarOpUC := appops.NewAgregarOperacionUseCase(casosRepo, operacionesRepo, auditRepo)

	// ── Handlers ─────────────────────────────────────────────────────────────
	bootstrapH := handlers.NewBootstrapHandler(bootstrapUC, estudiosRepo, usuariosRepo, bancosRepo)
	bancosH := handlers.NewBancosHandler(bancosRepo, usuariosRepo)
	casosH := handlers.NewCasosHandler(createCaseUC, updateCaseUC, transicionUC, casosRepo)
	clientesH := handlers.NewClientesHandler(clientesRepo)
	operacionesH := handlers.NewOperacionesHandler(agregarOpUC, operacionesRepo)
	plazosH := handlers.NewPlazosHandler()

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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

		// Bancos + asignaciones
		r.Route("/v1/bancos", func(r chi.Router) {
			r.Get("/", bancosH.Listar)
			r.Post("/", bancosH.Crear)
			r.Patch("/{id}", bancosH.Actualizar)
			r.Delete("/{id}", bancosH.Eliminar)
			r.Get("/{id}/usuarios", bancosH.ListarUsuarios)
			r.Post("/{id}/usuarios", bancosH.AsignarUsuario)
			r.Delete("/{id}/usuarios/{usuarioId}", bancosH.DesasignarUsuario)
		})

		// Usuarios del estudio (para selector de asignación)
		r.Get("/v1/usuarios", bancosH.ListarUsuariosEstudio)

		r.Route("/v1/clientes", func(r chi.Router) {
			r.Post("/", clientesH.Crear)
			r.Get("/{id}", clientesH.Obtener)
		})

		r.Route("/v1/casos", func(r chi.Router) {
			r.Get("/", casosH.Listar)
			r.Post("/", casosH.Crear)
			r.Get("/{id}", casosH.Obtener)
			r.Patch("/{id}", casosH.Actualizar)
			r.Post("/{id}/transicion", casosH.Transicionar)
			r.Post("/{id}/operaciones", operacionesH.Crear)
			r.Get("/{id}/operaciones", operacionesH.Listar)
		})

		r.Route("/v1/casos/{casoID}/plazos", func(r chi.Router) {
			r.Get("/", plazosH.ListarPorCaso)
		})
	})

	return r
}
