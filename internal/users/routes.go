package users

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

var gh handler

type handler struct {
	db *sql.DB
}

func SetRoutes(router chi.Router, db *sql.DB) {
	gh = handler{db: db}

	router.Get("/users", gh.List)
	router.Post("/users", gh.Create)
	router.Get("/users/{id}", gh.GetByID)
	router.Put("/users/{id}", gh.Update)
	router.Delete("/users/{id}", gh.SoftDelete)
}
