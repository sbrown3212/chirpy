package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJTWandValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "superSecretToken"
	durationLive := time.Minute
	durationExpired := -time.Minute

	tests := []struct {
		name           string
		userID         uuid.UUID
		createSecret   string
		validateSecret string
		expiresIn      time.Duration
		expectedErr    error
	}{
		{
			name:           "Valid JWT",
			userID:         userID,
			createSecret:   tokenSecret,
			validateSecret: tokenSecret,
			expiresIn:      durationLive,
			expectedErr:    nil,
		},
		{
			name:           "Expired JWT",
			userID:         userID,
			createSecret:   tokenSecret,
			validateSecret: tokenSecret,
			expiresIn:      durationExpired,
			expectedErr:    jwt.ErrTokenExpired,
		},
		{
			name:           "Incorrect Secret",
			userID:         userID,
			createSecret:   tokenSecret,
			validateSecret: "incorrectSecret",
			expiresIn:      durationExpired,
			expectedErr:    jwt.ErrTokenSignatureInvalid,
		},
	}

	for _, tt := range tests {
		ss, err := MakeJWT(tt.userID, tt.createSecret, tt.expiresIn)
		if err != nil {
			t.Errorf("MakeJWT() error: %v", err)
		}

		validatedUserID, err := ValidateJWT(ss, tt.validateSecret)
		if tt.expectedErr != nil {
			if err == nil {
				t.Errorf("ValidateJWT() expected error, but got none")
			} else if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateJWT() got error: %v, want: %v", err, tt.expectedErr)
			}
		}

		if (tt.expectedErr == nil) && (err != nil) {
			t.Errorf("ValidateJWT() unexpected error: %v", err)
		}

		if (tt.expectedErr == nil) && (tt.userID != validatedUserID) {
			t.Errorf("Validated uuid does not match original uuid")
		}
	}
}
