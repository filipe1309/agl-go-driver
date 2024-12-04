package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := Decode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := InsertDB(h.db, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = id

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func InsertDB(db *sql.DB, user *User) (id int64, err error) {
	stmt := `INSERT INTO users (name, login, password, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`

	err = db.QueryRow(stmt, user.Name, user.Login, user.Password, user.UpdatedAt).Scan(&id)
	if err != nil {
		return -1, err
	}

	return
}
