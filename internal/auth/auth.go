package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	isCorrectPW, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("error comparing password and hash")
	}

	return isCorrectPW, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")

	if bearerToken == "" {
		return bearerToken, fmt.Errorf("authorization header not found")
	}

	split := strings.Split(bearerToken, " ")
	if split[0] != "Bearer" {
		return "", fmt.Errorf("invalid bearer token format")
	}

	token := split[1]
	return token, nil
}
