package handlers

import (
	"fmt"
	"net/http"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerAuthRefresh(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)

	if r.Body != http.NoBody {
		fmt.Println("Refresh token req includes body")
		RespondWithError(w, http.StatusUnauthorized, "Request should only have a header")
		return
	}

	user, err := cfg.DB.UserForRefreshToken(refreshToken)
	if err != nil {
		fmt.Println("Token not found or expired")
		RespondWithError(w, http.StatusUnauthorized, "Not a valid token")
		return
	}

	signedToken, err := auth.GenerateJTW(user.Email, user.Id, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	resp := response{
		Token: signedToken,
	}

	RespondWithJSON(w, http.StatusOK, resp)

}
