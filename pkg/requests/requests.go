package requests

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func doRequest(method, path string, body io.Reader, headers map[string]string, auth bool) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:8090%s", path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if _, ok := headers["Content-Type"]; !ok {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if auth {
		token, err := readCacheToken()
		if err != nil {
			log.Println("No token found, please login")
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return http.DefaultClient.Do(req)
}
