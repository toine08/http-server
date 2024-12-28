package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChipsById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	uuid, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, 404, "Error while getting the value", err)
		return
	}
	row, err := cfg.dbQueries.ChirpByID(r.Context(), uuid)
	if err != nil {
		respondWithError(w, 500, "Error there is no ID who match the chirp", err)
	}
	respondWithJSON(w, 200, Chirp{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Body:      row.Body,
		UserID:    row.UserID,
	})

}
