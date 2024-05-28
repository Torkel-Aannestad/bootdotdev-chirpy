package handlers

import "net/http"

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, _ *http.Request) {
	cfg.Mu.Lock()
	defer cfg.Mu.Unlock()

	cfg.FileserverHits = 0

}
