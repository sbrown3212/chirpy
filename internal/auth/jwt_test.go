package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJTW(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "superSecretToken"
	durationLive := time.Minute
	durationExpired := -time.Minute

	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Valid JWT",
			userID:      userID,
			tokenSecret: tokenSecret,
			expiresIn:   durationLive,
			wantErr:     false,
		},
		{
			name:        "Expired JWT",
			userID:      userID,
			tokenSecret: tokenSecret,
			expiresIn:   durationExpired,
			wantErr:     true,
		},
	}

	// TODO: create loop for handling test cases
}

// func TestValidateJWT(t *testing.T) {}
