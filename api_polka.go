package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jabuta/chirpy/internal/auth"
	"github.com/jabuta/chirpy/internal/database"
)

func (cfg *apiConfig) polkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type polkaRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	if err := auth.CheckHeaderApi(r, cfg.polka_key); err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	reqBody := polkaRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to decode json")
		return
	}
	if reqBody.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, "goaway")
		return
	}
	user, err := cfg.db.MakeRed(reqBody.Data.UserID)
	if errors.Is(err, database.UserNotFound) {
		respondWithError(w, http.StatusNotFound, "User Not Found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
