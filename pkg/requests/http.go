package requests

import (
	"io"
	"net/http"
)

func AuthenticatedPost(path string, body io.Reader) ([]byte, error) {
	resp, err := doRequest(http.MethodPost, path, body, true)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func AuthenticatedGet(path string) ([]byte, error) {
	resp, err := doRequest(http.MethodGet, path, nil, true)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}
