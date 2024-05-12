package handlers

import (
	"net/http"
	"sort"
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
