package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (a *ApiConfig) HandlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	cleanedChirp, err := validateChirp(params.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdChirp, err := a.DB.CreateChirp(cleanedChirp)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not write to db")
		return
	}

	RespondWithJSON(w, http.StatusCreated, createdChirp)
}

func validateChirp(str string) (string, error) {
	const chirpMaxLength = 140
	if len(str) > chirpMaxLength {
		return "", errors.New("chirp is too long")
	}

	type prophaneWords struct {
		list map[string]bool
	}
	badWords := prophaneWords{
		list: map[string]bool{
			"kerfuffle": true,
			"sharbert":  true,
			"fornax":    true,
		},
	}

	words := strings.Fields(str)
	for i, subs := range words {
		if _, ok := badWords.list[strings.ToLower(subs)]; ok {
			words[i] = "****"
		}
	}

	cleanedStr := strings.Join(words, " ")
	return cleanedStr, nil
}
