package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := ReadAllDB(h.db)
	if err != nil {
		// TODO: Check if the error is sql.ErrNoRows and return 404
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func ReadAllDB(db *sql.DB) ([]User, error) {
	stmt := `SELECT * FROM users WHERE deleted = FALSE`

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return RestoreAll(rows)
}
