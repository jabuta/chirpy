package database

import (
	"encoding/json"
	"log"
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

func CreateDB(path string) (*DB, error) {
	var db DB
	filePath := path + "database.json"
	f, err := os.Create(filePath)
	if err != nil {
		return &DB{}, err
	}
	emptyDB, err := json.Marshal(DBstructure{Chirps: map[int]Chirp{}})
	if err != nil {
		return nil, err
	}
	f.Write(emptyDB)
	err = f.Close()
	if err != nil {
		return &DB{}, err
	}
	db.path = filePath
	db.mux = &sync.RWMutex{}
	return &db, nil
}

func (db *DB) load() (DBstructure, error) {
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

func (db *DB) save(memDB DBstructure) error {
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

func (db *DB) CreateChirp(body string) (Chirp, error) {

	db.mux.Lock()
	defer db.mux.Unlock()

	memDB, err := db.load()
	if err != nil {
		log.Print("db load err")
		return Chirp{}, err
	}
	largestK := 0
	for k := range memDB.Chirps {
		if k > largestK {
			largestK = k
		}
	}
	largestK++
	memDB.Chirps[largestK] = Chirp{
		Body: body,
		ID:   largestK,
	}

	err = db.save(memDB)
	if err != nil {
		return Chirp{}, err
	}

	return memDB.Chirps[largestK], nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	memDB, err := db.load()
	if err != nil {
		return []Chirp{}, err
	}
	log.Print(memDB)
	chirpList := make([]Chirp, 0, len(memDB.Chirps))
	for _, chirp := range memDB.Chirps {
		chirpList = append(chirpList, chirp)
	}
	log.Print(chirpList)
	return chirpList, nil
}
