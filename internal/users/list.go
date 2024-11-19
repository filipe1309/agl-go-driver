package users

import (
	"database/sql"
	"encoding/json"
	"log"
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

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.Deleted)
		if err != nil {
			log.Println(err)
			continue
		}

		users = append(users, user)
	}

	return users, nil
}
