package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/toine08/http-server/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		fmt.Printf("Error finding user: %v\n", err) // Debug print
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		fmt.Printf("Error comparing passwords: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
