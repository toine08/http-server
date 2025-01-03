package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/toine08/http-server/internal/auth"
	"github.com/toine08/http-server/internal/database"
)

func (cfg *apiConfig) handleDeleteChirpsById(w http.ResponseWriter, req *http.Request) {
	//verify identity
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No valid bearer token provided", err)
		return
	}

	claims, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token", err)
		return
	}
	//get the ChirpID
	chirpId := req.PathValue("chirpID")
	uuid, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, 404, "Error while getting the value", err)
		return
	}
	if _, err := cfg.dbQueries.DeleteChirpByID(req.Context(), database.DeleteChirpByIDParams{
		ID:     uuid,
		UserID: claims,
	}); err != nil {
		respondWithError(w, http.StatusForbidden, "Error while deleting the chirp", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
