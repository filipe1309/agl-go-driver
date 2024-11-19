package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	folder := new(Folder)
	err := json.NewDecoder(r.Body).Decode(folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = folder.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = Update(h.db, int64(id), folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Get folder by ID and return it

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

func Update(db *sql.DB, id int64, folder *Folder) (int64, error) {
	folder.UpdatedAt = time.Now()

	stmt := `UPDATE folders SET name = $1, updated_at = $2 WHERE id = $3`
	result, err := db.Exec(stmt, folder.Name, folder.UpdatedAt, id)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
