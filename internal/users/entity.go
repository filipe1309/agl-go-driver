package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"
)

var (
	ErrPasswordEmpty = errors.New("password can't be empty")
	ErrPasswordTooShort = errors.New("password is too short, minimum 6 characters")
)

func New(name, login, password string) (*User, error) {
	now := time.Now()
	user := User{
		Name:      name,
		Login:     login,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := user.SetPassword(password)
	if err != nil {
		return nil, err
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

func (u *User) SetPassword(password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}

	if len(password) < 6 {
		return ErrPasswordTooShort
	}

	u.Password = fmt.Sprintf("%x", md5.Sum([]byte(password)))

	return nil
}
