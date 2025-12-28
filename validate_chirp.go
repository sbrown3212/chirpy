package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVal struct {
		Cleaned string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode parameters", err)
		return
	}

	const maxLength = 140
	if len(params.Body) > maxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleansed := cleanse(params.Body)

	respondWithJSON(w, http.StatusOK, returnVal{Cleaned: cleansed})
}

func cleanse(chirp string) string {
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
