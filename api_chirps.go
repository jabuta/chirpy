package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jabuta/chirpy/internal/database"
)

func middlewarePostChirp(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		chirp, err := db.CreateChirp(cleanBody)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, chirp)
	}
}

func middlewareGetChirp(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		chirps, err := db.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, chirps)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
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
