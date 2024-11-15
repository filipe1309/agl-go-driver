package users

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = SoftDelete(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func SoftDelete(db *sql.DB, id int64) error {
	stmt := `UPDATE users SET updated_at = $1, deleted = true WHERE id = $2`

	_, err := db.Exec(stmt, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}
