package main

import (
	"net/http"

	"github.com/sbrown3212/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"failed to get refresh token",
			err,
		)
		return
	}

	rowsAffected, err := cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to revoke refresh token",
			err,
		)
		return
	}
	if rowsAffected == 0 {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			nil,
		)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
