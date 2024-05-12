package handlers

import "net/http"

func (a *ApiConfig) HandlerReset(w http.ResponseWriter, _ *http.Request) {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	a.FileserverHits = 0

}
