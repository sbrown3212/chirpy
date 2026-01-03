package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
)

func handleChirps(w http.ResponseWriter, r *http.Request) {
	type parameteres struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp string `json:"chirp"`
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

	_ = cleanseChirp(params.Body)

	// TODO: create migration for chirps table
	// TODO: create query to add chrip to database
	// TODO: Save chirp to database
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
