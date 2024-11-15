package folders

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	db *sql.DB
}

func SetRoutes(router chi.Router, db *sql.DB) {
	h := handler{db: db}

	// router.Get("/folders", h.List)
	router.Post("/folders", h.Create)
	// router.Get("/folders/{id}", h.GetByID)
	router.Put("/folders/{id}", h.Update)
	// router.Delete("/folders/{id}", h.SoftDelete)
}
