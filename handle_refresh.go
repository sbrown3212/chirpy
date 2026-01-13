package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/sbrown3212/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"failed to get refresh token",
			err,
		)
	}

	dbRT, err := cfg.db.GetRefreshTokenByToken(r.Context(), refreshToken)

	// token not found
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			err,
		)
		return
	}

	// expired
	if time.Now().After(dbRT.ExpiresAt) {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			nil,
		)
		return
	}

	// revoked
	if dbRT.RevokedAt.Valid {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			nil,
		)
		return
	}

	// unexpected error
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"unexpected error",
			err,
		)
		return
	}

	newAccessToken, err := auth.MakeJWT(dbRT.UserID, cfg.jwtsecret, time.Hour)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to create new access token",
			err,
		)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: newAccessToken,
	})
}
