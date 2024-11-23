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

func encryptPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func (u *User) SetPassword(password string) {
	u.Password = encryptPassword(password)
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
