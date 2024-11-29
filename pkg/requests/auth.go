package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
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

	resp, err := doRequest(http.MethodPost, path, &body, nil, false)
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

	replacer := strings.NewReplacer("\r", "", "\n", "")
	data = []byte(replacer.Replace(string(data)))

	_, err = file.Write(data)

	return err
}

func readCacheToken() (string, error) {
	data, err := os.ReadFile(".cacheToken")
	if err != nil {
		return "", err
	}

	var cache cacheToken
	err = json.Unmarshal(data, &cache)
	return cache.Token, err
}
