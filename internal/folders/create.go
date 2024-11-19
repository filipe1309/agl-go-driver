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

	if folder.ParentID.Int64 == 0 {
		folder.ParentID.Valid = false
	}

	id, err := InsertDB(h.db, folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	folder.ID = id

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

func InsertDB(db *sql.DB, folder *Folder) (int64, error) {
	stmt := `INSERT INTO folders (parent_id, name, updated_at) VALUES ($1, $2, $3)`

	var parentID any
	if folder.ParentID.Valid {
		parentID = folder.ParentID.Int64
	} else {
		parentID = nil
	}

	result, err := db.Exec(stmt, parentID, folder.Name, folder.UpdatedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
