package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/database"
)

func (cfg *ApiConfig) HandlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	chirps := []database.Chirp{}
	if authorId != "" {
		authorIdInt, _ := strconv.Atoi(authorId)
		dat, err := cfg.DB.GetChirpsByAuthor(authorIdInt)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "could not read db")
			return
		}
		chirps = dat

	} else {
		dat, err := cfg.DB.GetChirps()
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "could not read db")
			return
		}
		chirps = dat
	}

	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id > chirps[j].Id
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})
	}

	RespondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *ApiConfig) HandlerChirpsRetrieveByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error")
	}

	chirp, err := cfg.DB.GetChirpById(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "could not find by the given id")
		return
	}

	RespondWithJSON(w, http.StatusOK, chirp)

}
