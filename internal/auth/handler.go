package auth

import "net/http"

type Authenticated interface {
	GetID() int64
	GetName() string
}

type authenticateFunc func(string, string) (Authenticated, error)

type handler struct {
	authenticate authenticateFunc
}

func NewHandlerAuth(fn authenticateFunc) func (http.ResponseWriter, *http.Request) {
	h := handler{fn}

	return h.authHTTP
}
