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

		router.Post("/", gh.Create)

		router.Group(func(r chi.Router) {
			r.Use(auth.Validate)

			router.Get("/", gh.List)
			router.Get("/{id}", gh.GetByID)
			router.Put("/{id}", gh.Update)
			router.Delete("/{id}", gh.SoftDelete)
		})
	})
}
