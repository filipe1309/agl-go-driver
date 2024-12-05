package auth

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func validate(tokenStr string) (*Claims, error, int) {
	if tokenStr == "" {
		return nil, errors.New("missing token"), http.StatusUnauthorized
	}

	claims := new(Claims)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err, http.StatusUnauthorized
		}

		return nil, err, http.StatusBadRequest
	}

	if !token.Valid {
		return nil, errors.New("invalid token"), http.StatusUnauthorized
	}

	return claims, nil, http.StatusOK
}
