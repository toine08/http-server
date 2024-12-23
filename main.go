package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.Handle("/", http.FileServer(http.Dir(".")))
	// Serve the logo.png file at the /assets path
	mux.Handle("/assets/", http.FileServer(http.Dir(".")))
	server.ListenAndServe()
}
