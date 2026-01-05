package main

import (
	"encoding/json"
	"net/http"

	"github.com/sbrown3212/chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type paramemters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := paramemters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode parameters", err)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Username/password is incorrect", err)
		return
	}

	ok, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if !ok || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Username/password is incorrect", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})
}
