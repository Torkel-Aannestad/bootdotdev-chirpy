package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (a *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "fuck off")
	}

	user, err := a.DB.GetUserByEmail(params.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "something went wrong")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "wrong username and password")
	}
	RespondWithJSON(w, http.StatusOK, user.User)

}
