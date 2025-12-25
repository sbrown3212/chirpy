package main

import (
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	_ = r

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	msg := http.StatusText(http.StatusOK)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Println("error writing /healthz response:", err)
	}
}
