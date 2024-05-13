package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/auth"
)

type Session struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (a *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Password           string `json:"password"`
		Email              string `json:"email"`
		Expires_in_seconds int    `json:"expires_in_seconds,omitempty"`
	}
	type response struct {
		Session
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not read request body")
	}

	user, err := a.DB.GetUserByEmail(params.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "something went wrong")
	}
	err = auth.CompareHashAndPassword(user.PasswordHash, params.Password)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "wrong username and password")
		return
	}

	signedToken, err := auth.GenerateJTW(user.Email, user.Id, params.Expires_in_seconds, a.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	resp := response{
		Session: Session{
			Id:    user.Id,
			Email: user.Email,
			Token: signedToken,
		},
	}

	RespondWithJSON(w, http.StatusOK, resp)

}
