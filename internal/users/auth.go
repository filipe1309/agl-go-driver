package users

func (h *handler) authenticate(login, password string) (*User, error) {
	stmt := `SELECT * FROM users WHERE login = $1 AND password = $2`

	var user User
	row := h.db.QueryRow(stmt, login, encryptPassword(password))
	err := row.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.Deleted)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func Authenticate(login, password string) (*User, error) {
	return gh.authenticate(login, password)
}
