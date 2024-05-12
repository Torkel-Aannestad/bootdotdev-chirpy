package main

import (
	"log"
	"net/http"
	"time"

	"github.com/torkelaannestad/bootdotdev-chirpy/handlers"
	"github.com/torkelaannestad/bootdotdev-chirpy/internal/database"
)

func newApiConfig(db *database.DB) *handlers.ApiConfig {
	return &handlers.ApiConfig{
		FileserverHits: 0,
		DB:             db,
	}
}

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	db, err := database.NewDb("./database.json")
	if err != nil {
		log.Printf("DB error: %v", err)
	}
	cfg := newApiConfig(&db)

	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("GET /app/*", cfg.MiddlewareMetricsInc(fileServer))

	mux.HandleFunc("GET /api/reset", cfg.HandlerReset)
	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)
	mux.HandleFunc("GET /api/chirps", cfg.HandlerChirpsRetrieve)
	mux.HandleFunc("POST /api/chirps", cfg.HandlerChirpsCreate)
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerMetrics)

	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(srv.ListenAndServe())
}
