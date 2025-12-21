package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	svr := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Fatal(svr.ListenAndServe())
}
