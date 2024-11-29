package requests

import (
	"errors"
	"io"
	"net/http"
)

func validateResponse(resp *http.Response) ([]byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 399 {
		return nil, errors.New(string(data))
	}

	return data, nil
}

func AuthenticatedPost(path string, body io.Reader) ([]byte, error) {
	return Post(path, body, nil, true)
}

func Post(path string, body io.Reader, headers map[string]string, auth bool) ([]byte, error) {
	resp, err := doRequest(http.MethodPost, path, body, nil, auth)
	if err != nil {
		return nil, err
	}

	return validateResponse(resp)
}

func AuthenticatedGet(path string) ([]byte, error) {
	resp, err := doRequest(http.MethodGet, path, nil, nil, true)
	if err != nil {
		return nil, err
	}

	return validateResponse(resp)
}

func AuthenticatedPut(path string, body io.Reader) ([]byte, error) {
	resp, err := doRequest(http.MethodPut, path, body, nil, true)
	if err != nil {
		return nil, err
	}

	return validateResponse(resp)
}

func AuthenticatedDelete(path string) error {
	_, err := doRequest(http.MethodDelete, path, nil, nil, true)
	return err
}
