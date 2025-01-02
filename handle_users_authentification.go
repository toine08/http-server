package main

import (
	"net/http"
	"os"

	"github.com/toine08/http-server/internal/auth"
)

func (cfg *apiConfig) handleAuthenticatedRequest(w http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No valid bearer token provided", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, os.Getenv("TOKEN_SECRET"))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token", err)
		return
	}

	// You can then use the userID to fetch user-specific data or perform user-related operations
	// Example continuation...
	user, err := cfg.dbQueries.GetUserById(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve user data", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
