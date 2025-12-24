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
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		// fileserverHits: (*atomic.Int32).Load(0),
	}

	fs := http.FileServer(http.Dir(filepathRoot))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)

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
