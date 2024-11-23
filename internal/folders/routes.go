package folders

import (
	"database/sql"

	"github.com/filipe1309/agl-go-driver/internal/auth"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	db *sql.DB
}

func SetRoutes(router chi.Router, db *sql.DB) {
	h := handler{db: db}

	router.Group(func(r chi.Router) {
		r.Use(auth.Validate)

		r.Get("/folders", h.List)
		r.Post("/folders", h.Create)
		r.Get("/folders/{id}", h.GetByID)
		r.Put("/folders/{id}", h.Update)
		r.Delete("/folders/{id}", h.SoftDelete)
	})
}
