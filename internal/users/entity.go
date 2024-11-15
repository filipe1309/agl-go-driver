package users

import "time"

func New(name, login, password string) (*User, error) {
	now := time.Now()
	user := User{
		Name:      name,
		Login:     login,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return &user, nil
}

type User struct {
	ID        int
	Name      string
	Login     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	LastLogin time.Time
	Deleted   bool
}
