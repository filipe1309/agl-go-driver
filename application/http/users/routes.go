package users

import (
	"github.com/filipe1309/agl-go-driver/factories"
	"github.com/filipe1309/agl-go-driver/internal/auth"
	"github.com/filipe1309/agl-go-driver/repositories"
	"github.com/go-chi/chi/v5"
)

var gh handler

type handler struct {
	repo    repositories.UserWriteRepository
	factory *factories.UserFactory
}

func SetRoutes(router chi.Router, repo repositories.UserWriteRepository, uf *factories.UserFactory) {
	gh = handler{repo, uf}

	router.Route("/users", func(r chi.Router) {

		r.Post("/", gh.Create)

		r.Group(func(r chi.Router) {
			r.Use(auth.Validate)

			r.Get("/", gh.List)
			r.Get("/{id}", gh.GetByID)
			r.Put("/{id}", gh.Update)
			r.Delete("/{id}", gh.SoftDelete)
		})
	})
}
