package main

import (
	"net/http"

	"github.com/toine08/http-server/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {
	if req.ContentLength > 0 {
		respondWithError(w, http.StatusBadRequest, "This endpoint doesn't accept data in the body.", nil)
		return
	}

	refresh_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "An error occured, sorry", err)
		return
	}

	if err := cfg.dbQueries.UpdateRevokedAt(req.Context(), refresh_token); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while revoking the token", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)

}
