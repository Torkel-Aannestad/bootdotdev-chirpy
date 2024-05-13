package database

import "errors"

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"passordhash"`
}

func (db *DB) CreateUser(email string, hashedPassord string) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbData.Users) + 1
	newUser := User{
		Id:           id,
		Email:        email,
		PasswordHash: hashedPassord,
	}

	dbData.Users[id] = newUser

	err = db.writeDB(dbData)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}
func (db *DB) UpdateUser(id int, email string, hashedPassord string) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	updatedUser := User{
		Id:           id,
		Email:        email,
		PasswordHash: hashedPassord,
	}
	if _, ok := dbData.Users[id]; !ok {
		return User{}, err
	}
	dbData.Users[id] = updatedUser

	err = db.writeDB(dbData)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (db *DB) GetUserById(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{}
	for k, v := range dbStructure.Users {
		if v.Email == email {
			user = dbStructure.Users[k]
		}
	}

	return user, nil
}
