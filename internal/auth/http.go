package auth

import (
	"encoding/json"
	"net/http"
)

type ServiceHTTP struct {
	handler
}

func (svc *ServiceHTTP) authenticate(rw http.ResponseWriter, r *http.Request) {
	// Payload decode
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	token, err, code := svc.auth(credentials)
	if err != nil {
		http.Error(rw, err.Error(), code)
		return
	}

	rw.Write([]byte(token))
}

func HandleHttpAuth(fn authenticateFunc) func(http.ResponseWriter, *http.Request) {
	svc := ServiceHTTP{handler{fn}}

	return svc.authenticate
}
