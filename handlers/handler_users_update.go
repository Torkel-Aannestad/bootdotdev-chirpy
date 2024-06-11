package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	bearerToken, ok := r.Header["Authorization"]
	if !ok {
		fmt.Printf("Not found: %v", "Authorization")
		RespondWithError(w, http.StatusUnauthorized, "missing authorization header")
	}
	jwtToken := strings.Split(bearerToken[0], " ")[1]

	userId, err := auth.VerifyTokenAndGetUser(jwtToken, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not authorized")
		return
	}
	idStr, _ := strconv.Atoi(userId)
	user, _ := cfg.DB.GetUserById(idStr)

	fmt.Printf("user: %v", user)

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not parse the request")
	}

	passwordHash, err := auth.GeneratePasswordHash(params.Password)
	if err != nil {
		log.Print("could not encrypt password")
	}

	updatedUser, err := cfg.DB.UpdateUser(user.Id, params.Email, passwordHash)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not update user")
	}

	resp := response{
		User: User{
			Id:    updatedUser.Id,
			Email: updatedUser.Email,
		},
	}

	RespondWithJSON(w, http.StatusOK, resp)

}
