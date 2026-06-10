package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"poly.app/api/internal/adapters/http/handlers"
	"poly.app/api/internal/adapters/http/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)

	// Health check público
	r.Get("/health", handlers.Health)

	// Rutas protegidas
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Use(middleware.RequireTenantScope)

		casosH := handlers.NewCasosHandler()
		r.Route("/v1/casos", func(r chi.Router) {
			r.Get("/", casosH.Listar)
			r.Post("/", casosH.Crear)
			r.Get("/{id}", casosH.Obtener)
			r.Post("/{id}/transicion", casosH.Transicionar)
		})

		plazosH := handlers.NewPlazosHandler()
		r.Route("/v1/casos/{casoID}/plazos", func(r chi.Router) {
			r.Get("/", plazosH.ListarPorCaso)
		})
	})

	return r
}
