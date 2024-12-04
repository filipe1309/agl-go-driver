package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ReadDB(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err = DecodeAndUpdate(r.Body, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = UpdateDB(h.db, int64(id), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateDB(db *sql.DB, id int64, user *User) (int64, error) {
	user.UpdatedAt = time.Now()
	stmt := `UPDATE users SET name = $1, updated_at = $2, last_login = $3 WHERE id = $4`

	result, err := db.Exec(stmt, user.Name, user.UpdatedAt, user.LastLogin, id)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
