package users

func newLoginUser(login, password string) *User {
	u := &User{
		Login:    login,
		Password: password,
	}
	encryptPassword(u)
	return u
}
