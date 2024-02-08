package main

import (
	"encoding/json"
	"net/http"

	"github.com/jabuta/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type userRequestVals struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) newUser(w http.ResponseWriter, r *http.Request) {
	reqBody := userRequestVals{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len([]byte(reqBody.Password)) > 72 {
		respondWithError(w, http.StatusBadRequest, "Password is too long")
		return
	}
	if len(reqBody.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "Password is too short")
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 0)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(reqBody.Email, hashedPwd)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	reqBody := userRequestVals{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len([]byte(reqBody.Password)) > 72 {
		respondWithError(w, http.StatusBadRequest, "Password is too long")
		return
	}
	if len(reqBody.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "Password is too short")
		return
	}

	user, err := cfg.db.ReadUser(reqBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user not found")
		return
	}
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(reqBody.Password)); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	respondWithJSON(w, http.StatusOK, database.ReturnUser{
		Email: user.Email,
		ID:    user.ID,
	})
}
