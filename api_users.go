package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jabuta/chirpy/internal/auth"
)

type userRequestVals struct {
	Password string `json:"password"`
	Email    string `json:"email"`
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

	hashedPwd, err := auth.HashPassword(reqBody.Password)
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
	type loginResponseVals struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		TokenRefresh string `json:"refresh_token"`
	}
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
	if err := auth.CheckPassword(user.PasswordHash, reqBody.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tokenAccess, err := auth.CreateJwtString(user.ID, cfg.jwtSecret, "chirpy-access")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	tokenRefresh, err := auth.CreateJwtString(user.ID, cfg.jwtSecret, "chirpy-refresh")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, loginResponseVals{
		ID:           user.ID,
		Email:        user.Email,
		Token:        tokenAccess,
		TokenRefresh: tokenRefresh,
	})
}

func (cfg *apiConfig) modifyUser(w http.ResponseWriter, r *http.Request) {
	reqBody := userRequestVals{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//get auth from header
	if len(strings.Fields(r.Header.Get("Authorization"))) != 2 {
		respondWithError(w, http.StatusBadRequest, "token bad format")
		return
	}
	authToken := strings.Fields(r.Header.Get("Authorization"))[1]
	uid, err := auth.CheckToken(authToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error hashing password")
		return
	}
	retuser, err := cfg.db.UpdateUser(uid, reqBody.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error saving user")
		return
	}
	respondWithJSON(w, http.StatusOK, retuser)
}
