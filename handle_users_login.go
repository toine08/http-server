package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/toine08/http-server/internal/auth"
	"github.com/toine08/http-server/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"` // Using a pointer to handle optionality

	}
	type response struct {
		ID            uuid.UUID `json:"id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Token         string    `json:"token"`
		Refresh_Token string    `json:"refresh_token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil || auth.CheckPasswordHash(params.Password, user.HashedPassword) != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	expireIn := time.Hour // Default to 1 hour
	if params.ExpiresInSeconds != nil {
		expireIn = time.Duration(*params.ExpiresInSeconds) * time.Second
		if expireIn > time.Hour {
			expireIn = time.Hour
		}
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, expireIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating the token", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Sorry an error occured while generating the refresh token", err)
	}
	expireAt := time.Now().Add(60 * 24 * time.Hour) // 60 days from now
	cfg.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: sql.NullTime{Time: expireAt, Valid: true},
	})

	respondWithJSON(w, http.StatusOK, response{
		ID:            user.ID,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Email:         user.Email,
		Token:         token,
		Refresh_Token: refreshToken,
	})
}
