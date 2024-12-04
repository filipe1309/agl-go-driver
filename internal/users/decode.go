package users

import (
	"encoding/json"
	"io"
)

func decode(body io.ReadCloser) (*User, error) {
	user := new(User)
	err := json.NewDecoder(body).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func DecodeAndCreate(body io.ReadCloser) (*User, error) {
	user, err := decode(body)
	if err != nil {
		return nil, err
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	user.ChangePassword(user.Password)

	return user, nil
}

func DecodeAndUpdate(body io.ReadCloser, u *User) (*User, error) {
	user, err := decode(body)
	if err != nil {
		return nil, err
	}

	err = user.ChangeName(user.Name)
	if err != nil {
		return nil, err
	}

	return u, nil
}
