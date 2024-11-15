package files

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	file := new(File)
	err := json.NewDecoder(r.Body).Decode(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = file.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := Insert(h.db, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}

func Insert(db *sql.DB, file *File) (int64, error) {
	stmt := `INSERT INTO files (folder_id, owner_id, name, type, path, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	result, err := db.Exec(stmt, file.FolderID, file.OwnerID, file.Name, file.Type, file.Path, file.UpdatedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
