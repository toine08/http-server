package main

import (
	"net/http"
	"time"

	"github.com/toine08/http-server/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, req *http.Request) {
	if req.ContentLength > 0 {
		respondWithError(w, http.StatusBadRequest, "This endpoint doesn't accept data in the body.", nil)
		return
	}

	refresh_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "An error occured, sorry", err)
		return
	}

	tokenInfo, err := cfg.dbQueries.SelectRefreshTokenByToken(req.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve token info", err)
		return
	}

	if tokenInfo.ExpiresAt.Valid && tokenInfo.ExpiresAt.Time.After(time.Now()) && (!tokenInfo.RevokedAt.Valid || tokenInfo.RevokedAt.Time.IsZero()) {
		expiresIn := 1 * time.Hour

		token, err := auth.MakeJWT(tokenInfo.UserID, cfg.tokenSecret, time.Duration(expiresIn))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error while creating the token", err)
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
		return
	} else {
		respondWithError(w, http.StatusUnauthorized, "Token is expired or has been revoked", err)
		return
	}
}
