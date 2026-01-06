package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashedPassword(password string) (string, error) {
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
