package files

import (
	"database/sql"

	"github.com/filipe1309/agl-go-driver/internal/auth"
	"github.com/filipe1309/agl-go-driver/internal/bucket"
	"github.com/filipe1309/agl-go-driver/internal/queue"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	db     *sql.DB
	bucket *bucket.Bucket
	queue  *queue.Queue
}

func SetRoutes(router chi.Router, db *sql.DB, b *bucket.Bucket, q *queue.Queue) {
	h := handler{db, b, q}

	router.Group(func(r chi.Router) {
		r.Use(auth.Validate)

		r.Post("/folders", h.Create)
		r.Put("/folders/{id}", h.Update)
		r.Delete("/folders/{id}", h.SoftDelete)
	})
}
