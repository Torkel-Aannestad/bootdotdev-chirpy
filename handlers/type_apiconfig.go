package handlers

import (
	"sync"

	"github.com/torkelaannestad/bootdotdev-chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits int
	DB             *database.DB
	Mu             sync.Mutex
}
