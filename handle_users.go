package main

import (
	"encoding/json"
	"net/http"

	"github.com/sbrown3212/chirpy/internal/auth"
	"github.com/sbrown3212/chirpy/internal/database"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode parameters", err)
		return
	}

	email := params.Email
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "email must not be empty", nil)
		return
	}
	password := params.Password
	if password == "" {
		respondWithError(w, http.StatusBadRequest, "password must not be empty", nil)
		return
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to hash password", err)
		return
	}

	db := cfg.db
	dbUser, err := db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create user in database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.CreatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
	})
}
