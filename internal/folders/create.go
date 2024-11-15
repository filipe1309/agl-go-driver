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
	stmt := `INSERT INTO folders (name, parent_id, updated_at) VALUES ($1, $2, $3)`

	result, err := db.Exec(stmt, folder.Name, folder.ParentID, folder.UpdatedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
