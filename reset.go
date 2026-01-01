package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	_ = r

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)

	msg := fmt.Sprintln("Reset successful")

	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Println("error writing /reset response")
	}

	// TODO: Make query to delete all users in db
	// TODO: Delete all users in db within this funciton
}
