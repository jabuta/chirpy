package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	type requestVals struct {
		Email string `json:"email"`
	}

	reqBody := requestVals{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(reqBody.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.db.GetUsersList()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func (cfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request) {
	userStuct, err := cfg.db.GetUsersMap()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	i, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, ok := userStuct[i]
	if !ok {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
