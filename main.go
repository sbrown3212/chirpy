package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("metrics middleware called")
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{}

	fs := http.FileServer(http.Dir(filepathRoot))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerMetricsReset)

	svr := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(svr.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	_ = r

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	msg := "OK"
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Println("error writing /healthz response:", err)
	}
}

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
