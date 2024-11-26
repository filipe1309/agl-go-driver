package requests

import (
	"io"
	"net/http"
)

func AuthenticatedPost(path string, body io.Reader) ([]byte, error) {
	return Post(path, body, nil, true)
}

func Post(path string, body io.Reader, headers map[string]string, auth bool) ([]byte, error) {
	resp, err := doRequest(http.MethodPost, path, body, nil, auth)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func AuthenticatedGet(path string) ([]byte, error) {
	resp, err := doRequest(http.MethodGet, path, nil, nil, true)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func AuthenticatedPut(path string, body io.Reader) ([]byte, error) {
	resp, err := doRequest(http.MethodPut, path, body, nil, true)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func AuthenticatedDelete(path string) error {
	_, err := doRequest(http.MethodDelete, path, nil, nil, true)
	return err
}
