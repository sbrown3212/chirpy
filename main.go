package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", fs))

	mux.HandleFunc("/healthz", handlerReadiness)

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
