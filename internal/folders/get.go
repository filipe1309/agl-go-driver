package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	folder, err := Read(h.db, int64(id))
	if err != nil {
		// TODO: Check if the error is sql.ErrNoRows and return 404
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

func Read(db *sql.DB, id int64) (*Folder, error) {
	stmt := `SELECT * FROM folders WHERE id = $1`
	
	var folder Folder
	row := db.QueryRow(stmt, id)
	err := row.Scan(&folder.ID, &folder.ParentID, &folder.Name, &folder.CreatedAt, &folder.UpdatedAt, &folder.Deleted)
	if err != nil {
		return nil, err
	}

	return &folder, nil
}
