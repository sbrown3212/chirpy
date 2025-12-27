package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)

		respBody := errorResponse{
			Error: "Unable to decode parameters",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp is too long")

		respBody := errorResponse{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	type validResponse struct {
		Valid bool `json:"valid"`
	}

	respBody := validResponse{Valid: true}
	data, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

type errorResponse struct {
	Error string `json:"error"`
}
