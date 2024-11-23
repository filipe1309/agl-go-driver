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
	return token.SignedString(jwtSecret)
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *handler) auth(rw http.ResponseWriter, r *http.Request) {
	// Payload decode
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate authenticatedUser (authenticate)
	authenticatedUser, err := h.authenticate(credentials.Username, credentials.Password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	// Create token
	token, err := createToken(authenticatedUser)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response with token
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"token": token})
}
