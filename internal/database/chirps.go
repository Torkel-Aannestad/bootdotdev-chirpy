package database

import (
	"errors"
	"strconv"
)

type Chirp struct {
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
}

func (db *DB) CreateChirp(body string, userId string) (Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	userIdStr, _ := strconv.Atoi(userId)
	id := len(dbData.Chirps) + 1
	newChirp := Chirp{
		Id:       id,
		Body:     body,
		AuthorId: userIdStr,
	}

	dbData.Chirps[id] = newChirp

	err = db.writeDB(dbData)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, val := range dbStructure.Chirps {
		chirps = append(chirps, val)
	}
	return chirps, nil
}
func (db *DB) GetChirpsByAuthor(authorId int) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := []Chirp{}
	for _, val := range dbStructure.Chirps {
		if val.AuthorId == authorId {
			chirps = append(chirps, val)
		}
		continue
	}
	return chirps, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("could not find chirp by the given id")
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbData, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbData.Chirps, id)

	err = db.writeDB(dbData)
	if err != nil {
		return err
	}

	return nil
}
