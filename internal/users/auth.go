package users

import "time"

func (h *handler) authenticate(login, password string) (*User, error) {
	stmt := `SELECT * FROM users WHERE login = $1 AND password = $2`
	user := newLoginUser(login, password)
	row := h.db.QueryRow(stmt, user.GetLogin(), user.GetPass())
	err := row.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.Deleted)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *handler) updateLastLogin(user *User) (int64, error) {
	user.LastLogin = time.Now()
	return UpdateDB(h.db, user.ID, user)
}

func Authenticate(login, password string) (u *User, err error) {
	u, err = gh.authenticate(login, password)
	if err != nil {
		return
	}

	_, err = gh.updateLastLogin(u)
	return
}
