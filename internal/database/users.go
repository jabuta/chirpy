package database

import "log"

type User struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

func (db *DB) CreateUser(email string) (User, error) {

	db.mux.Lock()
	defer db.mux.Unlock()

	memDB, err := db.loadDB()
	if err != nil {
		log.Print("db load err")
		return User{}, err
	}
	largestK := 0
	for k := range memDB.Users {
		if k > largestK {
			largestK = k
		}
	}
	largestK++
	memDB.Users[largestK] = User{
		Email: email,
		ID:    largestK,
	}

	err = db.saveDB(memDB)
	if err != nil {
		return User{}, err
	}

	return memDB.Users[largestK], nil
}

func (db *DB) GetUsersList() ([]User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	memDB, err := db.loadDB()
	if err != nil {
		return []User{}, err
	}
	UserList := make([]User, 0, len(memDB.Users))
	for _, User := range memDB.Users {
		UserList = append(UserList, User)
	}
	return UserList, nil
}

func (db *DB) GetUsersMap() (map[int]User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	memDB, err := db.loadDB()
	if err != nil {
		return map[int]User{}, err
	}
	return memDB.Users, nil
}
