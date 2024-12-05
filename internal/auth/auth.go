package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = "34lddd654j6ha6s54klj7dhja7sadjksldiushf8r9ybc9brbcyr32"

type Claims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func createToken(au Authenticated) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &Claims{
		UserID:   au.GetID(),
		UserName: au.GetName(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *handler) auth(creds Credentials) (token string, err error, code int) {
	// Validate authenticatedUser (authenticate)
	authenticatedUser, err := h.authenticate(creds.Username, creds.Password)
	if err != nil {
		return "", err, http.StatusUnauthorized
	}

	// Create token
	token, err = createToken(authenticatedUser)
	if err != nil {
		return "", err, http.StatusInternalServerError
	}

	return token, nil, http.StatusOK
}

func (h *handler) authHTTP(rw http.ResponseWriter, r *http.Request) {
	// Payload decode
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	token, err, code := h.auth(credentials)
	if err != nil {
		http.Error(rw, err.Error(), code)
		return
	}

	rw.Write([]byte(token))
}
