package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Auth(path, username, password string) error {
	creds := Credentials{username, password}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(creds)
	if err != nil {
		return err
	}

	resp, err := doRequest("POST", path, &body, false)
	if err != nil {
		return err
	}

	return createTokenCache(resp.Body)
}

type cacheToken struct {
	Token string `json:"token"`
}

func createTokenCache(body io.ReadCloser) error {
	token, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	file, err := os.Create(".cacheToken")
	if err != nil {
		return err
	}

	cache := cacheToken{string(token)}

	data, err := json.Marshal(&cache)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}
