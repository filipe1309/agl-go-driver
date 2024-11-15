package users

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db *sql.DB
}

func SetRoutes(router chi.Router, db *sql.DB) {
	h := Handler{db: db}

	// router.Get("/users", h.List)
	router.Post("/users", h.Create)
	router.Get("/users/{id}", h.GetByID)
	router.Put("/users/{id}", h.Update)
	router.Delete("/users/{id}", h.SoftDelete)
}
