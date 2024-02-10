package database

import (
	"errors"
	"log"
)

type User struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	ID           int    `json:"id"`
	ChirpyRed    bool   `json:"is_chirpy_red"`
}
type ReturnUser struct {
	Email     string `json:"email"`
	ID        int    `json:"id"`
	ChirpyRed bool   `json:"is_chirpy_red"`
}

var UserNotFound = errors.New("user Not Found")

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
		ChirpyRed:    false,
	}

	err = db.saveDB(memDB)
	if err != nil {
		return ReturnUser{}, err
	}

	return ReturnUser{
		Email:     memDB.Users[uid].Email,
		ID:        memDB.Users[uid].ID,
		ChirpyRed: memDB.Users[uid].ChirpyRed,
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

	if user, ok := memDB.Users[uid]; !ok {
		return ReturnUser{}, err
	} else {
		if email != "" {
			user.Email = email
		}
		if passwordHash != "" {
			user.PasswordHash = passwordHash
		}

		memDB.Users[uid] = user
	}

	if err := db.saveDB(memDB); err != nil {
		return ReturnUser{}, err
	}
	return ReturnUser{
		Email:     memDB.Users[uid].Email,
		ID:        memDB.Users[uid].ID,
		ChirpyRed: memDB.Users[uid].ChirpyRed,
	}, nil
}

func (db *DB) MakeRed(uid int) (ReturnUser, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	memDB, err := db.loadDB()
	if err != nil {
		log.Print("db load err")
		return ReturnUser{}, err
	}

	if user, ok := memDB.Users[uid]; !ok {
		return ReturnUser{}, UserNotFound
	} else {
		user.ChirpyRed = true
		memDB.Users[uid] = user
	}
	if err := db.saveDB(memDB); err != nil {
		return ReturnUser{}, err
	}
	return ReturnUser{
		Email:     memDB.Users[uid].Email,
		ID:        memDB.Users[uid].ID,
		ChirpyRed: memDB.Users[uid].ChirpyRed,
	}, nil
}
