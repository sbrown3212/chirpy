package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	_ = r

	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		log.Printf("error deleting users from database: %s", err)
	}

	msg := fmt.Sprintln("Reset successful")

	_, err = w.Write([]byte(msg))
	if err != nil {
		log.Println("error writing /reset response")
	}
}
