package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sbrown3212/chirpy/internal/database"
)

type chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"udpated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	type parameteres struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameteres{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters", err)
		return
	}

	const maxLength = 140
	if len(params.Body) > maxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirpBody := cleanseChirp(params.Body)

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirpBody,
		UserID: params.UserID,
	})
	if err != nil {
		log.Printf("error saving chirp to database: %s", err)
	}

	respondWithJSON(w, http.StatusCreated, response{
		chirp: chirp{
			ID:        dbChirp.UserID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		},
	})
}

func cleanseChirp(chirp string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}

	sliceOriginal := strings.Split(chirp, " ")

	lower := strings.ToLower(chirp)
	sliceLower := strings.Split(lower, " ")

	for i, word := range sliceLower {
		if slices.Contains(profanities, word) {
			sliceOriginal[i] = "****"
		}
	}

	clean := strings.Join(sliceOriginal, " ")
	return clean
}
