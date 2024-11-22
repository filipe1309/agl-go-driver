package files

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// Update updates a file name
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := ReadDB(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if file.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	_, err = UpdateDB(h.db, int64(id), file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}

func UpdateDB(db *sql.DB, id int64, file *File) (int64, error) {
	file.UpdatedAt = time.Now()

	stmt := `UPDATE files SET name = $1, updated_at = $2, deleted = $3 WHERE id = $4`
	result, err := db.Exec(stmt, file.Name, file.UpdatedAt, file.Deleted, id)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
