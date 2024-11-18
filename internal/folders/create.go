package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
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

	id, err := Insert(h.db, folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	folder.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

func Insert(db *sql.DB, folder *Folder) (int64, error) {
	stmt := `INSERT INTO folders (parent_id, name, updated_at) VALUES ($1, $2, $3)`

	result, err := db.Exec(stmt, folder.ParentID, folder.Name, folder.UpdatedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
