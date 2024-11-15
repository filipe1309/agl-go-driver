package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"
)

var (
	ErrPasswordEmpty    = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters long")
	ErrNameEmpty        = errors.New("name is required")
	ErrLoginEmpty       = errors.New("login is required")
)

func New(name, login, password string) (*User, error) {
	user := User{
		Name:      name,
		Login:     login,
		UpdatedAt: time.Now(),
	}
	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}

	err = user.Validate()
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

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameEmpty
	}

	if u.Login == "" {
		return ErrLoginEmpty
	}

	blankPassword := fmt.Sprintf("%x", md5.Sum([]byte("")))
	if u.Password == blankPassword {
		return ErrPasswordEmpty
	}

	return nil
}
