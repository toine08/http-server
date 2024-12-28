package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/toine08/http-server/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)

	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {

	hits := cfg.fileserverHits.Load()
	html := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
func (cfg *apiConfig) handleReset(w http.ResponseWriter, req *http.Request) {
	log.Printf("Platform is set to: %s", cfg.platform)
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err := cfg.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		log.Println("Error while deleting from the DB", err)
	}
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Users table clear")))
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	plateform := os.Getenv("PLATEFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	cfg := &apiConfig{
		platform:       plateform,
		fileserverHits: atomic.Int32{},
	}
	cfg.dbQueries = database.New(db)
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	// Serve the logo.png file at the /assets path
	mux.Handle("/assets/", http.FileServer(http.Dir(".")))
	// Add a readiness endpoint at /healthz to check if the server is ready to receive requests
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	//mux.HandleFunc("POST /api/validate_chirp", handleValidation)
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	// Start the server and listen on the specified port
	server.ListenAndServe()
}
