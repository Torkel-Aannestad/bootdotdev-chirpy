package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error")
	}
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "not a valid token")
		return
	}
	userId, err := auth.VerifyTokenAndGetUser(bearer, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "not a valid token")
		return
	}

	userIdInt, _ := strconv.Atoi(userId)
	chirp, err := cfg.DB.GetChirpById(id)
	if userIdInt != chirp.AuthorId {
		RespondWithError(w, http.StatusForbidden, "Sorry dude")
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "could not find by the given id")
		return
	}

	err = cfg.DB.DeleteChirp(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "could not find by the given id")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
