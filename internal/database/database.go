package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBstructure struct {
	Chirps map[int]Chirp        `json:"chirps"`
	Users  map[int]User         `json:"users"`
	Tokens map[string]time.Time `json:"token"`
}

func NewDB(filePath string) (*DB, error) {
	db := &DB{
		path: filePath,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) createDB() error {
	emptyDB := DBstructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
		Tokens: map[string]time.Time{},
	}
	return db.saveDB(emptyDB)
}

func (db *DB) saveDB(memDB DBstructure) error {
	rawDB, err := json.Marshal(memDB)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, rawDB, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBstructure, error) {
	rawDb, err := os.ReadFile(db.path)
	if err != nil {
		return DBstructure{}, err
	}
	parsedDB := DBstructure{}
	err = json.Unmarshal(rawDb, &parsedDB)
	if err != nil {
		return DBstructure{}, err
	}
	return parsedDB, nil
}
