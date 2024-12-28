package main

import "net/http"

func (cfg *apiConfig) handleAllChirps(w http.ResponseWriter, r *http.Request) {
	var chirps []Chirp
	rows, err := cfg.dbQueries.AllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error retrieving data", err)
		return
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
