package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	type requestVals struct {
		Body string `json:"body"`
	}

	reqBody := requestVals{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	const maxChirpLength = 140

	if len(reqBody.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	if len(reqBody.Body) <= 0 {
		respondWithError(w, http.StatusBadRequest, "chirp where?")
		return
	}

	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	cleanBody := removeBadWords(reqBody.Body, badWords)

	chirp, err := cfg.db.CreateChirp(cleanBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirpsList()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirpStuct, err := cfg.db.GetChirpsMap()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	i, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chirp, ok := chirpStuct[i]
	if !ok {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)
}

func removeBadWords(text string, badWords []string) string {
	words := strings.Split(text, " ")
	for i := range words {
		for _, badWord := range badWords {
			if strings.Contains(strings.ToLower(words[i]), badWord) {
				words[i] = strings.Replace(strings.ToLower(words[i]), badWord, "****", 1)
			}
		}
	}
	return strings.Join(words, " ")
}
