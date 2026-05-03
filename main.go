package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	mux.Handle("/app", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	mux.HandleFunc("/app/assets", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><pre>\n<a href=\"assets/logo.png\">logo.png</a>\n</pre></html>"))
	})

	mux.HandleFunc("/healthz", handlerReadiness)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
