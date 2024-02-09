package database

import (
	"errors"
	"log"
)

type User struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	ID           int    `json:"id"`
}
type ReturnUser struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

func (db *DB) CreateUser(email string, hashedPwd string) (ReturnUser, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	memDB, err := db.loadDB()
	if err != nil {
		log.Print("db load err")
		return ReturnUser{}, err
	}

	for _, v := range memDB.Users {
		if v.Email == email {
			return ReturnUser{}, errors.New("user already exists")
		}
	}

	uid := len(memDB.Users) + 1
	memDB.Users[uid] = User{
		Email:        email,
		PasswordHash: hashedPwd,
		ID:           uid,
	}

	err = db.saveDB(memDB)
	if err != nil {
		return ReturnUser{}, err
	}

	return ReturnUser{
		Email: memDB.Users[uid].Email,
		ID:    memDB.Users[uid].ID,
	}, nil
}

func (db *DB) ReadUser(email string) (User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	memDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, v := range memDB.Users {
		if v.Email == email {
			return v, nil
		}
	}
	return User{}, errors.New("no user")
}

func (db *DB) UpdateUser(uid int, email string, passwordHash string) (ReturnUser, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	memDB, err := db.loadDB()
	if err != nil {
		log.Print("db load err")
		return ReturnUser{}, err
	}

	if _, ok := memDB.Users[uid]; !ok {
		return ReturnUser{}, err
	}
	memDB.Users[uid] = User{
		ID:           uid,
		Email:        email,
		PasswordHash: passwordHash,
	}
	if err := db.saveDB(memDB); err != nil {
		return ReturnUser{}, err
	}
	return ReturnUser{
		Email: email,
		ID:    uid,
	}, nil

}
