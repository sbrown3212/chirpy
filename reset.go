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
		w.Write([]byte("Reset is only allowed in dev environment"))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset database: " + err.Error()))
	}

	msg := fmt.Sprintln("Hits reset to 0 and database reset to initial state.")

	_, err = w.Write([]byte(msg))
	if err != nil {
		log.Println("error writing /reset response")
	}
}
