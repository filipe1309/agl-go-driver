package users

import (
	"database/sql"

	"github.com/filipe1309/agl-go-driver/internal/auth"
	"github.com/go-chi/chi/v5"
)

var gh handler

type handler struct {
	db *sql.DB
}

func SetRoutes(router chi.Router, db *sql.DB) {
	gh = handler{db: db}

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
