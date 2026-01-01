package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	_ = r

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)

	html := fmt.Sprintf(`
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
		`, cfg.fileserverHits.Load())

	_, err := w.Write([]byte(html))
	if err != nil {
		log.Println("error writing /metrics response:", err)
	}
}
