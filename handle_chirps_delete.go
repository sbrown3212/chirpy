package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/sbrown3212/chirpy/internal/auth"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"failed to find access token",
			err,
		)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtsecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"failed to validate access token",
			err,
		)
		return
	}

	dbChirp, err := cfg.db.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "chirp not found", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "unexpected error", err)
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(
			w,
			http.StatusForbidden,
			"forbidden",
			errors.New("chirp does not belong to user"),
		)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to delete chirp",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
