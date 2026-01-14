package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sbrown3212/chirpy/internal/auth"
	"github.com/sbrown3212/chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type paramemters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := paramemters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to decode parameters",
			err,
		)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Username/password is incorrect",
			err,
		)
		return
	}

	ok, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if !ok || err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Username/password is incorrect",
			err,
		)
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtsecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to make JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to create refresh token",
			err,
		)
		return
	}
	_, err = cfg.db.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:  refreshToken,
			UserID: dbUser.ID,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to save refresh token",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
