package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type reqBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "no valid apikey in header")
		return
	}

	if apiKey != cfg.PolkaKey {
		RespondWithError(w, http.StatusUnauthorized, "no valid key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	requestBody := reqBody{}
	err = decoder.Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "unable to read request body")
		return
	}

	if requestBody.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.DB.UpdateUserChirpyRed(requestBody.Data.UserID, true)
	if err != nil {
		fmt.Printf("Could not update user to Chirpy Red: %v", err)
		RespondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
