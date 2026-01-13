package main

import (
	"encoding/json"
	"net/http"

	"github.com/sbrown3212/chirpy/internal/auth"
	"github.com/sbrown3212/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to find access token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to validate access token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to decode parameters",
			err,
		)
		return
	}

	HashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to hash password",
			err,
		)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: HashedPassword,
	})
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to update user",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email:     updatedUser.Email,
		},
	})
}
