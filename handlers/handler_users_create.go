package handlers

import (
	"encoding/json"
	"net/http"
)

func (a *ApiConfig) HandlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not parse the request")
	}

	user, err := a.DB.CreateUser(params.Email, []byte(params.Password))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not create user")
	}

	RespondWithJSON(w, http.StatusCreated, user)

}
