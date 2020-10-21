package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/pxxonline/traefik-plugin-cors/cors"
)

func main() {
	c := cors.AllowAll()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	n := negroni.Classic()
	n.Use(c)
	n.UseHandler(mux)
	n.Run(":3000")
}
