package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/auth"
)

type Session struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *ApiConfig) HandlerAuthLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "something went wrong")
	}
	err = auth.CompareHashAndPassword(user.PasswordHash, params.Password)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "wrong username and password")
		return
	}

	signedToken, err := auth.GenerateJTW(user.Email, user.Id, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken := auth.GenerateRefreshToken()
	err = cfg.DB.SaveRefreshToken(user.Id, refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
	}

	resp := response{
		Session: Session{
			Id:           user.Id,
			Email:        user.Email,
			IsChirpyRed:  user.IsChirpyRed,
			Token:        signedToken,
			RefreshToken: refreshToken,
		},
	}

	RespondWithJSON(w, http.StatusOK, resp)

}
