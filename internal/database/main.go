package database

import (
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBstructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

func CreateDB(path string) (DB, error) {
	var db DB
	filePath := path + "database.json"
	f, err := os.Create(filePath)
	if err != nil {
		return DB{}, err
	}
	err = f.Close()
	if err != nil {
		return DB{}, err
	}
	db.path = filePath
	db.mux = &sync.RWMutex{}
	return db, nil
}
