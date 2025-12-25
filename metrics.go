package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	_ = r

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	msg := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())

	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Println("error writing /metrics response:", err)
	}
}

func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	_ = r

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	cfg.fileserverHits.Store(0)

	msg := fmt.Sprintf("Hits: %v (count reset)", cfg.fileserverHits.Load())

	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Println("error writing /reset response:", err)
	}
}
