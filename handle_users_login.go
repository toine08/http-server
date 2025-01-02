package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/toine08/http-server/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"` // Using a pointer to handle optionality
	}
	type response struct {
		ID               uuid.UUID `json:"id"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
		Email            string    `json:"email"`
		ExpiresInSeconds int       `json:"expires_in_seconds"`
		Token            string    `json:"token"`
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

	expiresIn := time.Hour // Default to 1 hour
	if params.ExpiresInSeconds != nil {
		expiresIn = time.Duration(*params.ExpiresInSeconds) * time.Second
		if expiresIn > time.Hour {
			expiresIn = time.Hour
		}
	}

	token, err := auth.MakeJWT(user.ID, os.Getenv("TOKEN_SECRET"), expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating the token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:               user.ID,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		Email:            user.Email,
		ExpiresInSeconds: int(expiresIn.Seconds()), // Convert duration to seconds
		Token:            token,
	})
}
