package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/toine08/http-server/internal/database"
)

func (cfg *apiConfig) handleAllChirps(w http.ResponseWriter, req *http.Request) {
	var chirps []Chirp
	var rows []database.Chirp
	var err error

	s := req.URL.Query().Get("author_id")
	order := req.URL.Query().Get("sort")

	if s == "" {
		if order == "desc" {
			rows, err = cfg.dbQueries.AllChirpsDesc(req.Context())
			if err != nil {
				respondWithError(w, 500, "Error retrieving data", err)
				return
			}
		} else {
			rows, err = cfg.dbQueries.AllChirps(req.Context())
			if err != nil {
				respondWithError(w, 500, "Error retrieving data", err)
				return
			}
		}
	} else {
		userID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, 400, "Invalid author_id", err)
			return
		}
		if order == "desc" {
			rows, err = cfg.dbQueries.AllChirpsByUserIDDesc(req.Context(), userID)
			if err != nil {
				respondWithError(w, 500, "Error retrieving data", err)
				return
			}
		}
	}
	for _, row := range rows {
		var chirp Chirp
		chirp.ID = row.ID
		chirp.CreatedAt = row.CreatedAt
		chirp.UpdatedAt = row.UpdatedAt
		chirp.Body = row.Body
		chirp.UserID = row.UserID
		chirps = append(chirps, chirp)
	}

	respondWithJSON(w, 200, chirps)
}
