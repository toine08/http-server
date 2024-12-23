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
    // Serve files from the current directory under the /app/ path, stripping the /app/ prefix
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
    // Serve the logo.png file at the /assets path
	mux.Handle("/assets/", http.FileServer(http.Dir(".")))
    // Add a readiness endpoint at /healthz to check if the server is ready to receive requests
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
    // Start the server and listen on the specified port
	server.ListenAndServe()
}
