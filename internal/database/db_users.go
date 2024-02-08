package database

import (
	"errors"
	"log"
)

type User struct {
	Email        string `json:"email"`
	PasswordHash []byte `json:"passwordHash"`
	ID           int    `json:"id"`
}
type ReturnUser struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

func (db *DB) CreateUser(email string, hashedPwd []byte) (ReturnUser, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	memDB, err := db.loadDB()
	if err != nil {
		log.Print("db load err")
		return ReturnUser{}, err
	}

	_, ok := memDB.Users[email]
	if ok {
		return ReturnUser{}, errors.New("user already exists")
	}
	uid := len(memDB.Users) + 1
	memDB.Users[email] = User{
		Email:        email,
		PasswordHash: hashedPwd,
		ID:           uid,
	}

	err = db.saveDB(memDB)
	if err != nil {
		return ReturnUser{}, err
	}

	return ReturnUser{
		Email: memDB.Users[email].Email,
		ID:    memDB.Users[email].ID,
	}, nil
}

func (db *DB) ReadUser(email string) (User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	memDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := memDB.Users[email]
	if !ok {
		return User{}, errors.New("no user")
	}

	return user, nil
}
