package requests

import (
	"io"
	"net/http"
)

func Post(path string, body io.Reader, auth bool) (*http.Response, error) {
	return doRequest(http.MethodPost, path, body, auth)
}
