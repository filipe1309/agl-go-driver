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

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`
	Deleted   bool      `json:"-"`
}

func encryptPassword(u *User) string {
	u.Password = fmt.Sprintf("%x", md5.Sum([]byte(u.Password)))
	return u.Password
}

func (u *User) GetID() int64 {
	return u.ID
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) ChangeName(name string) error {
	if name == "" {
		return ErrNameEmpty
	}

	u.Name = name

	return nil
}

func (u *User) GetPass() string {
	return u.Password
}

func (u *User) GetLogin() string {
	return u.Login
}

func (u *User) ChangePassword(password string) {
	u.Password = password
	encryptPassword(u)
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameEmpty
	}

	if u.Password == "" {
		return ErrPasswordEmpty
	}

	if len(u.Password) < 6 {
		return ErrPasswordTooShort
	}

	if u.Login == "" {
		return ErrLoginEmpty
	}

	// blankPassword := fmt.Sprintf("%x", md5.Sum([]byte("")))
	// if u.Password == blankPassword {
	// 	return ErrPasswordEmpty
	// }

	return nil
}
