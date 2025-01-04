package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/toine08/http-server/internal/auth"
)

type data struct {
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleWebhooks(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error while retreving the apikey", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "API key not matching", err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if params.Event == "user.upgraded" {
		fmt.Printf("HERE IS THE PARAMS.EVENT: %v, and the ID:%v ", params.Event, params.Data.UserID)
		_, err := cfg.dbQueries.UpdateChirpyById(req.Context(), params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Error while updating user", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, "")
	} else {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}
}
