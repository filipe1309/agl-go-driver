package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	authpb "github.com/filipe1309/agl-go-driver/proto/v1/auth"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HTTPAuth(path, username, password string) error {
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

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return createTokenCache(string(token))
}

func GRPCAuth(username, password string) error {
	creds := &authpb.CredentialsRequest{
		Username: username,
		Password: password,
	}

	conn := GetGRPCConn()
	defer conn.Close()
	client := authpb.NewAuthServiceClient(conn)

	resp, err := client.Login(context.Background(), creds)
	if err != nil {
		return err
	}

	return createTokenCache(resp.Token)
}

type cacheToken struct {
	Token string `json:"token"`
}

func createTokenCache(token string) error {
	file, err := os.Create(".cacheToken")
	if err != nil {
		return err
	}

	cache := cacheToken{token}

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
