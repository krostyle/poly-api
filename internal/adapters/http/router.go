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
	appdash "poly.app/api/internal/application/dashboard"
	appdocs "poly.app/api/internal/application/documentos"
	appops "poly.app/api/internal/application/operaciones"
	"poly.app/api/internal/adapters/feriados"
	"poly.app/api/internal/adapters/storage"
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
	plazosRepo := persistence.NewPlazoRepo(pool)
	feriadosProvider := feriados.NewDBFeriadoProvider(pool)
	documentosRepo := persistence.NewDocumentoRepo(pool)
	tribunalesRepo := persistence.NewTribunalRepo(pool)
	blobStorage := storage.NewVercelBlobStorage()

	// ── Use cases ────────────────────────────────────────────────────────────
	bootstrapUC := appauth.NewBootstrapUseCase(estudiosRepo, usuariosRepo, bancosRepo)
	createCaseUC := appcasos.NewCreateCaseUseCase(casosRepo, clientesRepo, plazosRepo, feriadosProvider, auditRepo)
	updateCaseUC := appcasos.NewUpdateCaseUseCase(casosRepo, plazosRepo, feriadosProvider, auditRepo)
	transicionUC := appcasos.NewTransitionStateUseCase(casosRepo, plazosRepo, feriadosProvider, auditRepo)
	agregarOpUC := appops.NewAgregarOperacionUseCase(casosRepo, operacionesRepo, plazosRepo, feriadosProvider, auditRepo)
	subirDocUC := appdocs.NewSubirDocumentoUseCase(blobStorage, documentosRepo)
	dashboardUC := appdash.NewDashboardUseCase(pool)

	// ── Handlers ─────────────────────────────────────────────────────────────
	bootstrapH := handlers.NewBootstrapHandler(bootstrapUC, estudiosRepo, usuariosRepo, bancosRepo)
	bancosH := handlers.NewBancosHandler(bancosRepo, usuariosRepo)
	casosH := handlers.NewCasosHandler(createCaseUC, updateCaseUC, transicionUC, casosRepo, auditRepo)
	clientesH := handlers.NewClientesHandler(clientesRepo)
	operacionesH := handlers.NewOperacionesHandler(agregarOpUC, operacionesRepo)
	plazosH := handlers.NewPlazosHandler(plazosRepo, feriadosProvider)
	documentosH := handlers.NewDocumentosHandler(subirDocUC, documentosRepo)
	dashboardH := handlers.NewDashboardHandler(dashboardUC)
	usuariosH := handlers.NewUsuariosHandler(usuariosRepo)
	tribunalesH := handlers.NewTribunalesHandler(tribunalesRepo)

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
		r.Patch("/v1/me/rol", usuariosH.CompletarOnboarding)

		// Bancos — lecturas abiertas, mutaciones solo ADMIN
		r.Route("/v1/bancos", func(r chi.Router) {
			r.Get("/", bancosH.Listar)
			r.Get("/catalogo", bancosH.Catalogo)
			r.With(middleware.RequireRol("ADMIN")).Post("/", bancosH.Crear)
			r.With(middleware.RequireRol("ADMIN")).Patch("/{id}", bancosH.Actualizar)
			r.With(middleware.RequireRol("ADMIN")).Delete("/{id}", bancosH.Eliminar)
			r.Get("/{id}/usuarios", bancosH.ListarUsuarios)
			r.With(middleware.RequireRol("ADMIN")).Post("/{id}/usuarios", bancosH.AsignarUsuario)
			r.With(middleware.RequireRol("ADMIN")).Delete("/{id}/usuarios/{usuarioId}", bancosH.DesasignarUsuario)
		})

		// Usuarios — lectura abierta (se usa para selector de abogado), gestión solo ADMIN
		r.Get("/v1/usuarios", bancosH.ListarUsuariosEstudio)
		r.With(middleware.RequireRol("ADMIN")).Post("/v1/usuarios/invitar", usuariosH.Invitar)
		r.With(middleware.RequireRol("ADMIN")).Patch("/v1/usuarios/{id}/rol", usuariosH.ActualizarRol)

		r.Route("/v1/clientes", func(r chi.Router) {
			r.Post("/", clientesH.Crear)
			r.Get("/{id}", clientesH.Obtener)
			r.Patch("/{id}", clientesH.Actualizar)
		})

		r.Get("/v1/plazos", plazosH.ListarGlobal)

		r.Get("/v1/tribunales", tribunalesH.Listar)
		r.With(middleware.RequireRol("ADMIN")).Post("/v1/tribunales", tribunalesH.Crear)

		r.Route("/v1/casos", func(r chi.Router) {
			r.Get("/", casosH.Listar)
			r.Post("/", casosH.Crear)
			r.Get("/{id}", casosH.Obtener)
			r.Patch("/{id}", casosH.Actualizar)
			r.Delete("/{id}", casosH.Eliminar)
			r.With(middleware.RequireRol("ADMIN", "ABOGADO")).Post("/{id}/transicion", casosH.Transicionar)
			r.Get("/{id}/historial", casosH.Historial)
			r.Post("/{id}/operaciones", operacionesH.Crear)
			r.Get("/{id}/operaciones", operacionesH.Listar)
			r.Get("/{id}/plazos", plazosH.ListarPorCaso)
			r.Post("/{id}/plazos/{plazoID}/cumplir", plazosH.CumplirPlazo)
			r.Get("/{id}/documentos", documentosH.Listar)
			r.Post("/{id}/documentos", documentosH.Subir)
		})

		r.Route("/v1/dashboard", func(r chi.Router) {
			r.Get("/por-vencer", dashboardH.PorVencer)
			r.Get("/nuevos", dashboardH.Nuevos)
			r.Get("/estancados", dashboardH.Estancados)
			r.Get("/por-abogado", dashboardH.PorAbogado)
		})
	})

	return r
}
