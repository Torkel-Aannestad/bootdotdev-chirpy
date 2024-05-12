package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

func (a *ApiConfig) HandlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.DB.GetChirps()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not read db")
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	RespondWithJSON(w, http.StatusOK, chirps)

}

func (a *ApiConfig) HandlerChirpsRetrieveByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error")
	}

	chirp, err := a.DB.GetChirpById(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "could not find by the given id")
		return
	}

	RespondWithJSON(w, http.StatusOK, chirp)

}
