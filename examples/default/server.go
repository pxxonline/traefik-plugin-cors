package main

import (
	"net/http"

	"github.com/pxxonline/traefik-plugin-cors/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	// Use default options
	handler := cors.AllowAll().Handler(mux)
	http.ListenAndServe(":8080", handler)
}
