package database

import (
	"errors"
	"time"
)

func (db *DB) AddRevocation(token string) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	memDB, err := db.loadDB()
	if err != nil {
		return err
	}
	memDB.Tokens[token] = time.Now().UTC()
	err = db.saveDB(memDB)
	return err
}

func (db *DB) TokenIsValid(token string) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	memDB, err := db.loadDB()
	if err != nil {
		return err
	}
	if _, ok := memDB.Tokens[token]; ok {
		return errors.New("token revoked")
	}
	return nil
}
