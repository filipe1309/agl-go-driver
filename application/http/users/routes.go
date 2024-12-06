package users

import (
	"github.com/filipe1309/agl-go-driver/factories"
	"github.com/filipe1309/agl-go-driver/internal/auth"
	"github.com/filipe1309/agl-go-driver/repositories"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	repo    repositories.UserWriteRepository
	factory *factories.UserFactory
}

func SetRoutes(router chi.Router, repo repositories.UserWriteRepository, uf *factories.UserFactory) {
	h := handler{repo, uf}

	router.Route("/users", func(r chi.Router) {

		r.Post("/", h.Create)

		r.Group(func(r chi.Router) {
			r.Use(auth.ValidateTokenMiddleware)

			r.Get("/", h.List)
			r.Get("/{id}", h.GetByID)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.SoftDelete)
		})
	})
}
