package users

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := ReadAll(h.db)
	if err != nil {
		// TODO: Check if the error is sql.ErrNoRows and return 404
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func ReadAll(db *sql.DB) ([]User, error) {
	stmt := `SELECT id, name, login, created_at, updated_at, last_login FROM users WHERE deleted = false`

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin)
		if err != nil {
			log.Println(err)
			continue
		}

		users = append(users, user)
	}

	return users, nil
}
