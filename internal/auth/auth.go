package auth

import (
	"encoding/json"
	"net/http"

	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = "34lddd654j6ha6s54klj7dhja7sadjksldiushf8r9ybc9brbcyr32"

type Claims struct {
	UserID string `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func createToken(user users.User) (string, error) {
	return jwtSecret, nil
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Auth(rw http.ResponseWriter, r *http.Request) {
	// Payload decode
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user (authenticate)

	// Create token

	// Response with token
}
