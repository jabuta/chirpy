package database

import "log"

type Chirp struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {

	db.mux.Lock()
	defer db.mux.Unlock()

	memDB, err := db.loadDB()
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

	err = db.saveDB(memDB)
	if err != nil {
		return Chirp{}, err
	}

	return memDB.Chirps[largestK], nil
}

func (db *DB) GetChirpsList() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	memDB, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	chirpList := make([]Chirp, 0, len(memDB.Chirps))
	for _, chirp := range memDB.Chirps {
		chirpList = append(chirpList, chirp)
	}
	return chirpList, nil
}

func (db *DB) GetChirpsMap() (map[int]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	memDB, err := db.loadDB()
	if err != nil {
		return map[int]Chirp{}, err
	}
	return memDB.Chirps, nil
}
