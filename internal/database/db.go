package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	Path string
	Mu   *sync.RWMutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"user"`
}

func NewDb(path string) (DB, error) {
	db := DB{
		Path: path,
		Mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	if err != nil {
		return DB{}, err
	}
	return db, nil
}
func (db *DB) createDb() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.Path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDb()
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	data, err := os.ReadFile("./database.json")
	if err != nil {
		return DBStructure{}, err
	}

	dbStructure := DBStructure{}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.Path, dat, 0600)
	if err != nil {
		return err
	}

	return nil
}
