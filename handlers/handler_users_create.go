package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/auth"
)

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (cfg *ApiConfig) HandlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not parse the request")
	}

	passwordHash, err := auth.GeneratePasswordHash(params.Password)
	if err != nil {
		log.Print("could not encrypt password")
	}

	user, err := cfg.DB.CreateUser(params.Email, passwordHash)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not create user")
	}

	resp := response{
		User: User{
			Id:          user.Id,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	}

	RespondWithJSON(w, http.StatusCreated, resp)

}
