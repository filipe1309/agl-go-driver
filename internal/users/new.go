package users

func New(id int64, name, login, password string) (*User, error) {
	u := &User{
		ID:       id,
		Name:     name,
		Login:    login,
		Password: password,
	}

	err := u.Validate()
	if err != nil {
		return nil, err
	}

	// encryptPassword(u)

	return u, nil
}

func newLoginUser(login, password string) *User {
	u := &User{
		Login:    login,
		Password: password,
	}
	encryptPassword(u)
	return u
}
