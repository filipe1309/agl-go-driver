package users

import (
	"encoding/json"
	"io"
)

func Decode(body io.ReadCloser) (*User, error) {
	user := new(User)
	err := json.NewDecoder(body).Decode(user)
	if err != nil {
		return nil, err
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	user.SetPassword(user.Password)

	return user, nil
}
