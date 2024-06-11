package handlers

import (
	"sync"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits int
	DB             *database.DB
	JWTSecret      string
	PolkaKey       string
	Mu             sync.Mutex
}
