package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = user.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.SetPassword(user.Password)

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

func InsertDB(db *sql.DB, user *User) (int64, error) {
	stmt := `INSERT INTO users (name, login, password, updated_at) VALUES ($1, $2, $3, $4)`

	result, err := db.Exec(stmt, user.Name, user.Login, user.Password, user.UpdatedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
