package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

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
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if splitAuth[0] != "Bearer" || len(splitAuth) < 2 {
		return "", errors.New("invalid authorization header format")
	}

	token := splitAuth[1]
	return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if splitAuth[0] != "ApiKey" || len(splitAuth) < 2 {
		return "", errors.New("invalid authorization header format")
	}

	apiKey := splitAuth[1]
	return apiKey, nil
}
