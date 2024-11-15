package users

import (
	"database/sql"
	"net/http"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
}

func Insert(db *sql.DB, user *User) (int64, error) {
	stmt := `INSERT INTO users (name, login, password, updated_at) VALUES ($1, $2, $3, $4)`

	result, err := db.Exec(stmt, user.Name, user.Login, user.Password, user.UpdatedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
