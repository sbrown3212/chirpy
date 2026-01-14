package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	var (
		authorID       uuid.UUID
		filterByAuthor bool
		sortDir        string
	)

	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		parsedUUID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(
				w, http.StatusBadRequest, "unable to parse autor id", err,
			)
			return
		}

		authorID = parsedUUID
		filterByAuthor = true
	}

	switch sortString := r.URL.Query().Get("sort"); sortString {
	case "":
		sortDir = "asc"
	case "asc":
		sortDir = "asc"
	case "desc":
		sortDir = "desc"
	default:
		respondWithError(
			w,
			http.StatusBadRequest,
			"invalid value for sort query parameter",
			nil,
		)
		return
	}

	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get chirps from db", err)
		return
	}

	var chirps []Chirp
	for _, dbChirp := range dbChirps {
		if filterByAuthor && dbChirp.UserID != authorID {
			continue
		}
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	if sortDir == "desc" {
		sort.SliceStable(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "uuid could not be parsed", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to get chirp from database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}
